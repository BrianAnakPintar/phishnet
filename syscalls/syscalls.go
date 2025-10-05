package syscalls

import (
	"os/exec"
	"runtime"
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
