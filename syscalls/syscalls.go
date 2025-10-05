package syscalls

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

func OpenChrome(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "-a", "Google Chrome", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("google-chrome", url)
	}

	return cmd.Start()
}

// Not rly a syscall but close enough
func NormalizeRawURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty url")
	}
	// If no scheme present, assume https
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if u.Host == "" {
		return "", fmt.Errorf("missing host in URL")
	}
	// Normalize trailing slash: remove single/multiple trailing slashes from path
	if u.Path == "/" {
		u.Path = ""
	} else {
		u.Path = strings.TrimRight(u.Path, "/")
	}
	return u.String(), nil
}
