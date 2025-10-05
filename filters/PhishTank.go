package filters

import (
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/briananakpintar/phishnet/syscalls"
)

const PhishTankFilterName string = "PhishTankFilter"

// PhishTankFilter checks an incoming URL against a hard-coded map of known bad URLs.
type PhishTankFilter struct {
	bad map[string]bool
}

func (f *PhishTankFilter) Configure(cfg map[string]string) error {
	// Allow adding additional URLs via configuration values if provided
	for _, v := range cfg {
		v, err := syscalls.NormalizeRawURL(strings.TrimSpace(v))
		if err != nil {
			continue
		}
		if v == "" {
			continue
		}
		f.bad[v] = true
	}
	return nil
}

func (f *PhishTankFilter) Run(u string) (FilterResult, error) {
	key := strings.TrimSpace(u)
	if key == "" {
		return FilterResult{Proceed: true, Reason: fmt.Sprintf("[%s] empty URL", PhishTankFilterName)}, nil
	}
	if _, found := f.bad[key]; found {
		return FilterResult{
			Proceed: false,
			Reason:  fmt.Sprintf("[%s] URL blocked: %s", PhishTankFilterName, u),
		}, nil
	}
	return FilterResult{Proceed: true, Reason: fmt.Sprintf("[%s] not listed", PhishTankFilterName)}, nil
}

// locateDataFiles searches common locations for the PhishTank data files and
// returns the first pair of paths where either the gob or csv exists.
func locateDataFiles() (gobPath, csvPath string, err error) {
	candidates := []string{}
	// Prefer data directory next to the executable (useful when installed)
	if ex, e := os.Executable(); e == nil {
		exDir := filepath.Dir(ex)
		candidates = append(candidates, filepath.Join(exDir, "data"))
	}
	// Then usual project/workdir locations
	candidates = append(candidates, "./data", "data", filepath.Join("..", "data"))

	for _, d := range candidates {
		g := filepath.Join(d, "verified_online.gob")
		c := filepath.Join(d, "verified_online.csv")
		if _, err := os.Stat(g); err == nil {
			return g, c, nil
		}
		if _, err := os.Stat(c); err == nil {
			return g, c, nil
		}
	}
	return "", "", fmt.Errorf("could not find verified_online.gob or verified_online.csv in candidate locations")
}

func loadPhishTankGob() (map[string]bool, error) {
	gobPath, csvPath, err := locateDataFiles()
	if err != nil {
		return nil, err
	}

	// If gob exists, load it
	if _, err := os.Stat(gobPath); err == nil {
		f, err := os.Open(gobPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		dec := gob.NewDecoder(f)
		var urls []string
		if err := dec.Decode(&urls); err != nil {
			return nil, err
		}
		m := make(map[string]bool, len(urls))
		for _, u := range urls {
			u = strings.TrimSpace(u)
			if u == "" {
				continue
			}
			m[u] = true
		}
		return m, nil
	}

	// Gob not found: try to build from CSV then load
	if _, err := os.Stat(csvPath); err != nil {
		return nil, fmt.Errorf("neither %s nor %s found", gobPath, csvPath)
	}
	if err := buildGobFromCSV(csvPath, gobPath); err != nil {
		return nil, err
	}
	// Try loading again
	return loadPhishTankGob()
}

// buildGobFromCSV reads the CSV and writes a gob containing the list of URLs.
func buildGobFromCSV(csvPath, gobPath string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	// CSV has header; we'll skip it if present
	var urls []string
	first := true
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if first {
			first = false
			// If header row contains the expected first column name, skip it
			if len(rec) > 0 && strings.HasPrefix(strings.ToLower(rec[0]), "phish_id") {
				continue
			}
		}
		if len(rec) < 2 {
			continue
		}
		u := strings.TrimSpace(rec[1])
		u, err = syscalls.NormalizeRawURL(u)
		if err != nil {
			continue
		}
		if u == "" {
			continue
		}
		urls = append(urls, u)
	}

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(gobPath), 0755); err != nil {
		return err
	}
	gf, err := os.Create(gobPath)
	if err != nil {
		return err
	}
	defer gf.Close()
	enc := gob.NewEncoder(gf)
	if err := enc.Encode(urls); err != nil {
		return err
	}
	return nil
}

func NewPhishTank() Filter {
	m, err := loadPhishTankGob()
	if err != nil {
		// If loading fails, fall back to empty map but do not panic the program
		fmt.Println("Warning: failed to load PhishTank data:", err)
		m = make(map[string]bool)
	}
	return &PhishTankFilter{bad: m}
}

func init() {
	RegisterFilter(PhishTankFilterName, NewPhishTank)
}
