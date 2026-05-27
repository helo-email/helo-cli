// Package commands implements all CLI subcommands. The generated code in this
// package is intentionally light — most behaviour comes from the wrapped SDK.
package commands

import (
	"errors"
	"os"
	"sync"

	"github.com/balvig/helo-sdk-go"

	"github.com/balvig/helo-cli/internal/output"
)

// GlobalFlags carries the values of the root command's persistent flags.
// SetGlobalFlags is called from internal/cli at PreRun.
type GlobalFlags struct {
	Token   string
	BaseURL string
	JSON    bool
	Quiet   bool
}

var (
	globalFlagsMu sync.RWMutex
	currentFlags  GlobalFlags
)

// SetGlobalFlags records the latest persistent flag values.
func SetGlobalFlags(f GlobalFlags) {
	globalFlagsMu.Lock()
	defer globalFlagsMu.Unlock()
	currentFlags = f
}

// Flags returns the current persistent flag values.
func Flags() GlobalFlags {
	globalFlagsMu.RLock()
	defer globalFlagsMu.RUnlock()
	return currentFlags
}

// newClient builds an SDK client from the current flags + environment.
func newClient() (*helo.Helo, error) {
	f := Flags()

	token := f.Token
	if token == "" {
		token = os.Getenv("HELO_TOKEN")
	}
	if token == "" {
		return nil, errors.New("authentication token not set; provide --token or export HELO_TOKEN")
	}

	opts := []helo.HeloOption{}
	baseURL := f.BaseURL
	if baseURL == "" {
		baseURL = os.Getenv("HELO_BASE_URL")
	}
	if baseURL != "" {
		opts = append(opts, helo.WithBaseURL(baseURL))
	}

	return helo.NewHelo(token, opts...), nil
}

// emit writes a successful result using the configured output mode.
func emit(data any) error {
	f := Flags()
	return output.PrintJSON(os.Stdout, data, output.PrintOptions{Quiet: f.Quiet})
}
