package main

import (
	"fmt"
	"os"

	_ "github.com/briananakpintar/phishnet/filters"
	fishnet "github.com/briananakpintar/phishnet/fishnet"
	"github.com/briananakpintar/phishnet/syscalls"
	"github.com/briananakpintar/phishnet/ui"
)

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
		ui.Popup(reason, raw)
		return
	}

	syscalls.OpenChrome(raw)
}
