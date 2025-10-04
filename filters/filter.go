package filters

import (
	"fmt"
	"strings"
)

type FilterResult struct {
	Proceed bool
	Reason  string
}

type Filter interface {
	// Configure the filter with parameters from DSL (string->string)
	Configure(map[string]string) error

	// Run the filter on the provided URL
	Run(u string) (FilterResult, error)
}

// Factory function for a filter instance
type FilterFactory func() Filter

var filterRegistry = map[string]FilterFactory{}

func RegisterFilter(name string, factory FilterFactory) {
	filterRegistry[strings.ToLower(name)] = factory
}

func CreateFilter(name string) (Filter, error) {
	factory, ok := filterRegistry[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("unknown filter: %s", name)
	}
	return factory(), nil
}
