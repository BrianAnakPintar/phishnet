package filters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Name used to register the filter
const GeminiFilterName string = "GeminiFilter"

// GeminiFilter queries Google Gemini to determine
// whether a given URL is likely a phishing URL.
//
// Configurable properties (via Configure):
// - API_KEY (string): Google API key. If empty the filter will
//   be skipped (allows the URL).
// - MODEL (string, optional): model name to use. Default is "gemini-1.5-flash".

type GeminiFilter struct {
	API_KEY string
	Model   string
	client  *http.Client
}

func (f *GeminiFilter) Configure(config map[string]string) error {
	if key, ok := config["API_KEY"]; ok {
		f.API_KEY = key
	}
	if m, ok := config["MODEL"]; ok {
		f.Model = m
	} else {
		f.Model = "gemini-2.5-flash"
	}
	f.client = &http.Client{Timeout: 10 * time.Second}
	return nil
}

// Run queries Google Gemini. If no API_KEY is configured
// the filter will be skipped and allow the URL.
func (f *GeminiFilter) Run(u string) (FilterResult, error) {
	if f.API_KEY == "" || f.API_KEY == "NIL" {
		return FilterResult{Proceed: true, Reason: fmt.Sprintf("[%s] No API key configured, skipping", GeminiFilterName)}, nil
	}

	// Build a short, specific prompt that requests a concise answer.
	prompt := fmt.Sprintf("You are a security assistant. Answer ONLY 'YES' or 'NO' followed by a short reason (one sentence).\nIs the following URL likely a phishing site?\nURL: %s", u)

	// Construct the request payload for Gemini API
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return FilterResult{}, fmt.Errorf("failed to marshal Gemini request: %w", err)
	}

	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", f.Model)
	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return FilterResult{}, fmt.Errorf("invalid endpoint URL: %w", err)
	}
	q := reqURL.Query()
	q.Set("key", f.API_KEY)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", reqURL.String(), bytes.NewReader(b))
	if err != nil {
		return FilterResult{}, fmt.Errorf("failed to create Gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return FilterResult{}, fmt.Errorf("gemini request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FilterResult{}, fmt.Errorf("failed to read Gemini response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return FilterResult{}, fmt.Errorf("gemini API returned status %d: %s", resp.StatusCode, string(body))
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return FilterResult{}, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	// Extract text from Gemini response
	text := extractGeminiText(parsed)
	if text == "" {
		return FilterResult{Proceed: true, Reason: fmt.Sprintf("[%s] No response text, allowing", GeminiFilterName)}, nil
	}

	return interpretAnswer(text)
}

// extractGeminiText pulls the generated text from Gemini's response JSON.
func extractGeminiText(resp map[string]interface{}) string {
	if candidates, ok := resp["candidates"].([]interface{}); ok && len(candidates) > 0 {
		if cand, ok := candidates[0].(map[string]interface{}); ok {
			if content, ok := cand["content"].(map[string]interface{}); ok {
				if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
					if part, ok := parts[0].(map[string]interface{}); ok {
						if text, ok := part["text"].(string); ok {
							return text
						}
					}
				}
			}
		}
	}
	return ""
}

// interpretAnswer inspects the model reply for a clear YES/NO. Defaults to
// allowing the URL when uncertain.
func interpretAnswer(reply string) (FilterResult, error) {
	r := strings.TrimSpace(strings.ToLower(reply))

	fmt.Println("Gemini reply:", r)
	// Look for an explicit yes/no at the start
	if strings.HasPrefix(r, "yes") || strings.HasPrefix(r, "y") {
		return FilterResult{Proceed: false, Reason: fmt.Sprintf("[%s] Model answered YES: %s", GeminiFilterName, strings.TrimSpace(reply))}, nil
	}
	if strings.HasPrefix(r, "no") || strings.HasPrefix(r, "n") {
		return FilterResult{Proceed: true, Reason: fmt.Sprintf("[%s] Model answered NO: %s", GeminiFilterName, strings.TrimSpace(reply))}, nil
	}
	// Try to find explicit words anywhere in the reply
	if strings.Contains(r, "yes") {
		return FilterResult{Proceed: false, Reason: fmt.Sprintf("[%s] Model indicated phishing: %s", GeminiFilterName, strings.TrimSpace(reply))}, nil
	}
	if strings.Contains(r, "no") {
		return FilterResult{Proceed: true, Reason: fmt.Sprintf("[%s] Model indicated not phishing: %s", GeminiFilterName, strings.TrimSpace(reply))}, nil
	}
	// Unclear -> allow but indicate uncertainty
	return FilterResult{Proceed: true, Reason: fmt.Sprintf("[%s] Unclear model response, allowing: %s", GeminiFilterName, strings.TrimSpace(reply))}, nil
}

func NewGeminiFilter() Filter {
	return &GeminiFilter{API_KEY: "NIL"}
}

func init() {
	RegisterFilter(GeminiFilterName, NewGeminiFilter)
}
