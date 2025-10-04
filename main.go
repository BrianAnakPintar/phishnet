package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	fishnet "github.com/briananakpintar/phishnet/fishnet"
	"github.com/briananakpintar/phishnet/ui"
)

func openChrome(url string) error {
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

func usage() {
	fmt.Println("Usage: phishnet <url>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	raw := os.Args[1]
	filterChain := fishnet.NewFilterChain()
	customParams := make(map[string]string)
	customParams["RickRoll"] = regexp.QuoteMeta(raw)

	filterChain.Add("PhishTank", nil)
	filterChain.Add("Regex", customParams)

	pass, reason, err := filterChain.Run(raw)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	if !pass {
		fmt.Printf("Denied entry to site: %s\n%s\n", raw, reason)
		ui.Popup(reason)
		return
	}

	openChrome(raw)
}
