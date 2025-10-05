package fishnet

import (
	"fmt"
	"strings"

	f "github.com/briananakpintar/phishnet/filters"
)

type filterConfig struct {
	name   string
	params map[string]string
	inst   f.Filter
}

type FilterChain struct {
	filters []filterConfig
}

func NewFilterChain() *FilterChain {
	return &FilterChain{filters: []filterConfig{}}
}

// Adds a new filter to the filter chain.
func (c *FilterChain) Add(name string, params map[string]string) error {
	inst, err := f.CreateFilter(name)
	if err != nil {
		return err
	}
	if err := inst.Configure(params); err != nil {
		return err
	}
	c.filters = append(c.filters, filterConfig{name: name, params: params, inst: inst})
	return nil
}

func (c *FilterChain) Run(url string) (allowed bool, reason string, err error) {
	var logs []string

	if len(c.filters) == 0 {
		return true, "no filters configured", nil
	}

	for _, fc := range c.filters {
		res, err := fc.inst.Run(url)
		if err != nil {
			return false, "", fmt.Errorf("filter %s error: %w", fc.name, err)
		}

		if res.Proceed {
			logs = append(logs, fmt.Sprintf("[%s] PASS", fc.name))
		} else {
			// record the blocking filter and return the whole stack
			logs = append(logs, fmt.Sprintf("[%s] FAIL\n%s", fc.name, res.Reason))
			return false, fmt.Sprintf("Filter chain log:\n%s\n", strings.Join(logs, "\n")), nil
		}
	}

	// all filters passed — return the stack showing passes
	return true, fmt.Sprintf("Filter chain log:\n%s\n", strings.Join(logs, "\n")), nil
}
