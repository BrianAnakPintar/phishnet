package filters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GoogleSafeBrowsingFilter struct {
	API_KEY string
}

func (f *GoogleSafeBrowsingFilter) Configure(config map[string]string) error {
	key, ok := config["API_KEY"]
	if !ok {
		return fmt.Errorf("GoogleSafeBrowsing filter requires an API_KEY configuration")
	}
	f.API_KEY = key
	return nil
}

// Run the filter on the provided URL
func (f *GoogleSafeBrowsingFilter) Run(url string) (FilterResult, error) {
	// If no API key configured, skip this filter (allow)
	if f.API_KEY == "" {
		return FilterResult{Proceed: true, Reason: "[GoogleSafeBrowsing] No API key configured, skipping"}, nil
	}

	// Build request body according to GSB v4 API
	reqBody := map[string]interface{}{
		"client": map[string]string{
			"clientId":      "localguardwhale",
			"clientVersion": "1.0",
		},
		"threatInfo": map[string]interface{}{
			"threatTypes":      []string{"MALWARE", "SOCIAL_ENGINEERING", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION"},
			"platformTypes":    []string{"ANY_PLATFORM"},
			"threatEntryTypes": []string{"URL"},
			"threatEntries":    []map[string]string{{"url": url}},
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return FilterResult{}, fmt.Errorf("failed to marshal GSB request: %w", err)
	}

	endpoint := fmt.Sprintf("https://safebrowsing.googleapis.com/v4/threatMatches:find?key=%s", f.API_KEY)
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return FilterResult{}, fmt.Errorf("google safe browsing request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FilterResult{}, fmt.Errorf("failed to read GSB response: %w", err)
	}

	// The API returns 200 with an empty object when no matches are found
	if resp.StatusCode != http.StatusOK {
		return FilterResult{}, fmt.Errorf("gsb API returned status %d: %s", resp.StatusCode, string(body))
	}

	var parsed struct {
		Matches []map[string]interface{} `json:"matches"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return FilterResult{}, fmt.Errorf("failed to parse GSB response: %w", err)
	}

	if len(parsed.Matches) == 0 {
		return FilterResult{Proceed: true, Reason: "[GoogleSafeBrowsing] No threats found"}, nil
	}

	// Build a concise reason from the first match
	first := parsed.Matches[0]
	reason := fmt.Sprintf("[GoogleSafeBrowsing] Threat detected: %v", first)
	return FilterResult{Proceed: false, Reason: reason}, nil
}

func NewGSBFilter() Filter {
	return &GoogleSafeBrowsingFilter{API_KEY: "NIL"}
}

func init() {
	RegisterFilter("GoogleSafeBrowsing", NewGSBFilter)
}
