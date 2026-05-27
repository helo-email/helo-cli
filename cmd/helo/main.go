// Package main is the entry point for the helo CLI.
package main

import (
	"fmt"
	"os"

	"github.com/balvig/helo-cli/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
