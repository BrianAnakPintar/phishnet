package guardwhale

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	_ "github.com/briananakpintar/phishnet/filters"
	fishnet "github.com/briananakpintar/phishnet/fishnet"
	guardwhalepb "github.com/briananakpintar/phishnet/gen/guardwhalepb"
	"github.com/briananakpintar/phishnet/syscalls"
)

// ScanURL normalizes the provided raw URL, runs the filter chain and returns (allowed, reason, error).
// error is only non-nil for unexpected failures; allowed/reason convey policy result.
func ScanURL(raw string) (bool, string, error) {
	normalized, err := syscalls.NormalizeRawURL(raw)
	if err != nil {
		return false, fmt.Sprintf("invalid url: %v", err), nil
	}

	filterChain := fishnet.NewFilterChain()
	if err := fishnet.ParseIntoChain(filterChain, ""); err != nil {
		return false, fmt.Sprintf("failed to parse filters: %v", err), err
	}

	pass, reason, err := filterChain.Run(normalized)
	if err != nil {
		return false, fmt.Sprintf("error running filters: %v", err), err
	}
	return pass, reason, nil
}

// StartHTTPServer starts an HTTP server exposing /scan (GET ?url=... or POST JSON {"url":...}).
func StartHTTPServer(port int) error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var raw string
		if r.Method == http.MethodGet {
			raw = r.URL.Query().Get("url")
		} else if r.Method == http.MethodPost {
			var body struct {
				URL string `json:"url"`
			}
			dec := json.NewDecoder(r.Body)
			if err := dec.Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body"})
				return
			}
			raw = body.URL
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
			return
		}

		if raw == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "missing url"})
			return
		}

		pass, reason, err := ScanURL(raw)
		fmt.Println("Scanned URL:", raw, "Allowed:", pass, "Reason:", reason)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": reason})
			return
		}

		resp := map[string]interface{}{"allowed": pass, "reason": reason}
		if !pass {
			w.WriteHeader(http.StatusForbidden)
		}
		json.NewEncoder(w).Encode(resp)
	}

	http.HandleFunc("/scan", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "phishnet server running. Use /scan?url=<url> or POST /scan with {\"url\":\"...\"}\n")
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("phishnet HTTP server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}

// gRPC server implementation
type grpcServer struct {
	guardwhalepb.UnimplementedGuardWhaleServer
}

func (s *grpcServer) Scan(ctx context.Context, req *guardwhalepb.ScanRequest) (*guardwhalepb.ScanResponse, error) {
	url := req.GetUrl()
	pass, reason, err := ScanURL(url)
	if err != nil {
		// return response with failure reason (and no error) so client gets policy result
		return &guardwhalepb.ScanResponse{Allowed: false, Reason: reason}, nil
	}
	return &guardwhalepb.ScanResponse{Allowed: pass, Reason: reason}, nil
}

// StartGRPCServer starts a gRPC server on the given port and registers the GuardWhale service.
func StartGRPCServer(port int) error {
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	grpcSrv := grpc.NewServer()
	guardwhalepb.RegisterGuardWhaleServer(grpcSrv, &grpcServer{})
	log.Printf("gRPC server listening on %s", addr)
	return grpcSrv.Serve(lis)
}
