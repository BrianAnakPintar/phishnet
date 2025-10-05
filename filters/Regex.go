package filters

import (
	"fmt"
	"regexp"
)

const RegexFilterName string = "RegexFilter"

type RegexFilter struct {
	expressions []*regexp.Regexp
}

func (f *RegexFilter) Configure(config map[string]string) error {
	for _, val := range config {
		f.expressions = append(f.expressions, regexp.MustCompile(val))
	}
	return nil
}

// Run the filter on the provided URL
// If one of the expressions match. We return UNSAFE.
func (f *RegexFilter) Run(url string) (FilterResult, error) {
	res := FilterResult{Proceed: true, Reason: fmt.Sprintf("[%s] No regex matches", FilterName)}
	for _, expr := range f.expressions {
		match := expr.MatchString(url)
		if match {
			return FilterResult{
				Proceed: false,
				Reason:  fmt.Sprintf("[%s] URL: %s matched with pattern: %s\n", FilterName, url, expr),
			}, nil
		}
	}
	return res, nil
}

func NewRegexFilter() Filter {
	return &RegexFilter{expressions: make([]*regexp.Regexp, 0)}
}

func init() {
	RegisterFilter(FilterName, NewRegexFilter)
}
