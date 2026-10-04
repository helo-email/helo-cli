// Package commands implements all CLI subcommands. The generated code in this
// package is intentionally light — most behaviour comes from the wrapped SDK.
package commands

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/helo-email/helo-sdk-go"

	"github.com/helo-email/helo-cli/internal/output"
)

// userAgent identifies the CLI rather than the SDK it wraps, so requests made
// by a person at a terminal are distinguishable from an application's. Like the
// SDK, the CLI is versioned by git tag, so the version comes from the build
// info the toolchain embeds.
var userAgent = buildUserAgent()

func buildUserAgent() string {
	version := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		version = info.Main.Version
	}
	return fmt.Sprintf("helo-cli/%s (%s; %s/%s)",
		sanitizeVersion(version), runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// A module version can be a placeholder like "(devel)" for an untagged local
// build. The parentheses there would open a comment the User-Agent never
// closes, so strip anything that is not a token character.
func sanitizeVersion(version string) string {
	cleaned := strings.Map(func(r rune) rune {
		if r == '(' || r == ')' || r == ' ' {
			return -1
		}
		return r
	}, version)
	if cleaned == "" {
		return "unknown"
	}
	return cleaned
}

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

	opts := []helo.HeloOption{
		helo.WithUserAgent(userAgent),
	}
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
