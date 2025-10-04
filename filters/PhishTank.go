package filters

import (
	"fmt"
	"strings"
)

// PhishTankFilter checks an incoming URL against a hard-coded map of known bad URLs.
type PhishTankFilter struct {
	bad map[string]struct{}
}

func (f *PhishTankFilter) Configure(cfg map[string]string) error {
	// Allow adding additional URLs via configuration values if provided
	for _, v := range cfg {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		f.bad[v] = struct{}{}
	}
	return nil
}

func (f *PhishTankFilter) Run(u string) (FilterResult, error) {
	key := strings.TrimSpace(u)
	if key == "" {
		return FilterResult{Proceed: true, Reason: "[PhishTank] empty URL"}, nil
	}
	if _, found := f.bad[key]; found {
		return FilterResult{
			Proceed: false,
			Reason:  fmt.Sprintf("[PhishTank] URL blocked: %s", u),
		}, nil
	}
	return FilterResult{Proceed: true, Reason: "[PhishTank] not listed"}, nil
}

func NewPhishTank() Filter {
	// Hard-coded hashmap of malicious URLs for now.
	bad := map[string]struct{}{
		"http://malicious.example/":   {},
		"https://phish.example/login": {},
		"http://example.com/rickroll": {},
		"https://badactor.test/steal": {},
	}
	return &PhishTankFilter{bad: bad}
}

func init() {
	RegisterFilter("PhishTank", NewPhishTank)
}
