package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	_ "github.com/briananakpintar/phishnet/filters"
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

	// Populate the chain from the DSL file (defaults to fishnet/bootstrap.fn)
	if err := fishnet.ParseIntoChain(filterChain, ""); err != nil {
		fmt.Printf("Failed to parse filters from DSL: %v\n", err)
	}

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
