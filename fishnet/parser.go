package fishnet

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
Fishnet is a DSL where users can declare custom filters for PhishNet
*/

// FilterSpec represents a parsed filter declaration from the DSL.
type FilterSpec struct {
	Name   string
	Params map[string]string
}

// ParseFile parses a .fn DSL file and returns the declared filters in order.
// The DSL supports blocks of the form:
//
// FilterName:[
//
//	key=value
//	key2: value2
//
// ]
//
// Lines starting with '//' are treated as comments and ignored. Whitespace is trimmed.
func ParseFile(path string) ([]FilterSpec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var specs []FilterSpec
	var current *FilterSpec

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// strip inline comments (//)
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}
		if line == "" {
			continue
		}

		// If we're not currently inside a block, look for the start pattern
		if current == nil {
			// Accept patterns like: Name:[  or Name : [
			if strings.HasSuffix(line, ":[") {
				name := strings.TrimSpace(strings.TrimSuffix(line, ":["))
				if name == "" {
					return nil, fmt.Errorf("invalid filter declaration: %q", line)
				}
				current = &FilterSpec{Name: name, Params: make(map[string]string)}
				continue
			}
			// allow a one-line empty-params form: Name:[]
			if strings.HasSuffix(line, ":[]") {
				name := strings.TrimSpace(strings.TrimSuffix(line, ":[]"))
				if name == "" {
					return nil, fmt.Errorf("invalid filter declaration: %q", line)
				}
				specs = append(specs, FilterSpec{Name: name, Params: make(map[string]string)})
				continue
			}
			// ignore other lines outside blocks
			continue
		}

		// We are inside a block. Looking for end marker or key/value lines
		if line == "]" || strings.HasPrefix(line, "]") {
			// end of block
			specs = append(specs, *current)
			current = nil
			continue
		}

		// parse key/value line. Support key=value and key: value
		l := strings.TrimSuffix(line, ",") // allow trailing comma
		var key, val string
		if i := strings.Index(l, "="); i >= 0 {
			key = strings.TrimSpace(l[:i])
			val = strings.TrimSpace(l[i+1:])
		} else if i := strings.Index(l, ":"); i >= 0 {
			key = strings.TrimSpace(l[:i])
			val = strings.TrimSpace(l[i+1:])
		} else {
			// treat single tokens as values with an empty key
			key = strings.TrimSpace(l)
			val = ""
		}
		// remove surrounding quotes from val
		val = strings.Trim(val, "\"')")
		if key != "" {
			current.Params[key] = val
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan file: %w", err)
	}

	if current != nil {
		return nil, fmt.Errorf("unterminated filter block for %q", current.Name)
	}

	return specs, nil
}

// ParseIntoChain reads the DSL file at path and adds each declared filter into the provided FilterChain.
// If path is empty, it will attempt to use "fishnet/bootstrap.fn" in the current working directory.
func ParseIntoChain(c *FilterChain, path string) error {
	if path == "" {
		path = "/home/brian/Documents/PersonalProjects/LocalGuardWhale/fishnet/bootstrap.fn"
	}
	specs, err := ParseFile(path)
	if err != nil {
		return err
	}

	for _, s := range specs {
		if err := c.Add(s.Name, s.Params); err != nil {
			return fmt.Errorf("add filter %s: %w", s.Name, err)
		}
	}
	return nil
}
