// Package main is the entry point for the helo CLI.
package main

import (
	"os"

	"github.com/helo-email/helo-cli/internal/cli"
)

func main() {
	os.Exit(cli.Execute(os.Args[1:]))
}
