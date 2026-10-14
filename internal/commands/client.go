// Package commands implements all CLI subcommands. The generated code in this
// package is intentionally light — most behaviour comes from the wrapped SDK.
package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
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
	return fmt.Sprintf("helo-cli/%s (%s; %s/%s)",
		sanitizeVersion(cliVersion()), runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func cliVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "unknown"
}

// Version describes the CLI build and the SDK version compiled into it.
func Version() string {
	sdkVersion := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/helo-email/helo-sdk-go" {
				sdkVersion = dep.Version
				if dep.Replace != nil {
					sdkVersion = dep.Replace.Version
				}
				break
			}
		}
	}
	if sdkVersion == "" {
		sdkVersion = "local"
	}
	return fmt.Sprintf("%s (helo-sdk-go %s)", cliVersion(), sdkVersion)
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
	DryRun  bool
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

// Environment variables the API key is read from, in order, after --token.
// HELO_API_KEY is the name the API documentation uses.
var tokenEnvVars = []string{"HELO_API_KEY", "HELO_TOKEN"}

// newClient builds an SDK client from the current flags + environment.
func newClient() (*helo.Helo, error) {
	f := Flags()

	token := f.Token
	for _, name := range tokenEnvVars {
		if token != "" {
			break
		}
		token = os.Getenv(name)
	}
	if token == "" && f.DryRun {
		token = "dry-run" // never sent, and redacted from the printed request
	}
	if token == "" {
		return nil, &output.CredentialsError{Message: "no API key: pass --token or set $" + strings.Join(tokenEnvVars, " or $")}
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
	if f.DryRun {
		opts = append(opts, helo.WithHTTPClient(&http.Client{Transport: dryRunTransport{quiet: f.Quiet}}))
	}

	return helo.NewHelo(token, opts...), nil
}

// ErrDryRun stops a --dry-run request once it has been printed.
var ErrDryRun = errors.New("dry run: request not sent")

// IsDryRun reports whether err is a --dry-run request stopping as intended.
func IsDryRun(err error) bool {
	return errors.Is(err, ErrDryRun)
}

// dryRunTransport prints each request instead of sending it.
type dryRunTransport struct {
	quiet bool
}

type dryRunRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    any               `json:"body,omitempty"`
}

func (t dryRunTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	printed := dryRunRequest{Method: req.Method, URL: req.URL.String(), Headers: map[string]string{}}
	keys := make([]string, 0, len(req.Header))
	for key := range req.Header {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := req.Header.Get(key)
		if strings.EqualFold(key, "Authorization") {
			value = "Bearer <redacted>"
		}
		printed.Headers[key] = value
	}
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		// Decoded so it prints as readable JSON: the SDK's encoding escapes the
		// <, > and & that email bodies are full of.
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.UseNumber()
		if err := dec.Decode(&printed.Body); err != nil {
			printed.Body = string(body)
		}
	}

	var data any = map[string]any{"dryRun": true, "request": printed}
	if t.quiet {
		data = printed
	}
	if err := output.PrintJSON(os.Stdout, data, output.PrintOptions{Quiet: t.quiet}); err != nil {
		return nil, err
	}
	return nil, ErrDryRun
}

// emit writes a successful result using the configured output mode.
func emit(data any) error {
	f := Flags()
	return output.PrintJSON(os.Stdout, data, output.PrintOptions{Quiet: f.Quiet})
}
