package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/briananakpintar/phishnet/guardwhale"
	"github.com/briananakpintar/phishnet/syscalls"
	"github.com/briananakpintar/phishnet/ui"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  phishnet <url>               # scan a single URL and open if allowed")
	fmt.Println("  phishnet -server [-port N]   # run HTTP server on port N (default 8080)")
	fmt.Println("  phishnet -grpc [-port N]     # run gRPC server on port N (default 8080)")
}

func main() {
	serverMode := flag.Bool("server", false, "Run HTTP server")
	grpcMode := flag.Bool("grpc", false, "Run gRPC server")
	port := flag.Int("port", 8080, "Port for server (HTTP or gRPC)")
	flag.Parse()

	if *grpcMode {
		log.Printf("Starting phishnet gRPC server on :%d", *port)
		if err := guardwhale.StartGRPCServer(*port); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
		return
	}

	if *serverMode {
		log.Printf("Starting phishnet HTTP server on :%d", *port)
		if err := guardwhale.StartHTTPServer(*port); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
		return
	}

	if flag.NArg() < 1 {
		usage()
		return
	}

	raw := flag.Arg(0)
	pass, reason, err := guardwhale.ScanURL(raw)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if !pass {
		fmt.Printf("Denied entry to site: %s\n%s\n", raw, reason)
		ui.Popup(reason, raw)
		return
	}

	// Open the site in Chrome. Normalize again to ensure consistent form.
	normalized, err := syscalls.NormalizeRawURL(raw)
	if err != nil {
		// fallback to raw if normalization unexpectedly fails
		syscalls.OpenChrome(raw)
		return
	}
	syscalls.OpenChrome(normalized)
}
