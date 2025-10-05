package filters

import (
	"fmt"
	"strings"
)

const PhishTankFilterName string = "PhishTankFilter"

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

func NewPhishTank() Filter {
	return &PhishTankFilter{}
}

func init() {
	RegisterFilter(PhishTankFilterName, NewPhishTank)
}
