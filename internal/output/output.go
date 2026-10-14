// Package output renders command results as either a styled JSON envelope
// or raw JSON for piping into other tools, and reports errors as text or JSON
// with an exit code that says what kind of failure it was.
package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// Exit codes, so a script or agent can tell whether retrying might help.
const (
	ExitOK        = 0
	ExitError     = 1 // anything not covered below
	ExitUsage     = 2 // bad flags, arguments or input; fix the command
	ExitAuth      = 3 // missing, invalid or insufficient credentials (401, 403)
	ExitNotFound  = 4 // the resource doesn't exist (404)
	ExitRejected  = 5 // the API rejected the request (400, 409, 422); fix the request
	ExitRetryable = 6 // rate limited (429), server error (5xx) or no response; retrying may help
)

// Envelope wraps successful results.
type Envelope struct {
	OK   bool        `json:"ok"`
	Data interface{} `json:"data,omitempty"`
}

// ErrorEnvelope wraps failure results. Error is always the human-readable
// message; the other fields are filled in when the API returned problem details.
type ErrorEnvelope struct {
	OK        bool           `json:"ok"`
	Error     string         `json:"error"`
	ExitCode  int            `json:"exitCode"`
	Status    int            `json:"status,omitempty"`
	Code      string         `json:"code,omitempty"`
	Detail    string         `json:"detail,omitempty"`
	RequestID string         `json:"requestId,omitempty"`
	Errors    map[string]any `json:"errors,omitempty"`
}

// PrintOptions controls how PrintJSON renders results.
type PrintOptions struct {
	Quiet bool // when true, emit only the raw data without an envelope
}

// PrintJSON writes data to w as pretty-printed JSON.
func PrintJSON(w io.Writer, data any, opts PrintOptions) error {
	enc := newEncoder(w)
	if opts.Quiet {
		if data == nil {
			return nil
		}
		return enc.Encode(data)
	}
	return enc.Encode(Envelope{OK: true, Data: data})
}

// PrintError writes an error to stderr, as an ErrorEnvelope when asJSON is set.
func PrintError(err error, asJSON bool) {
	if asJSON {
		_ = newEncoder(os.Stderr).Encode(NewErrorEnvelope(err))
		return
	}
	fmt.Fprintln(os.Stderr, "Error:", err)
}

// newEncoder writes indented JSON, leaving <, > and & as they are: output is
// read by people and programs, not embedded in HTML, and email bodies are full of them.
func newEncoder(w io.Writer) *json.Encoder {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc
}

// NewErrorEnvelope describes err, including any problem details the API sent.
func NewErrorEnvelope(err error) ErrorEnvelope {
	env := ErrorEnvelope{OK: false, Error: err.Error(), ExitCode: ExitCode(err)}
	var apiErr apiError
	if errors.As(err, &apiErr) {
		env.Status = apiErr.Status()
		problem := apiErr.Problem()
		env.Code, _ = problem["code"].(string)
		env.Detail, _ = problem["detail"].(string)
		env.RequestID, _ = problem["requestId"].(string)
		env.Errors, _ = problem["errors"].(map[string]any)
	}
	return env
}

// apiError is the part of the SDKs' APIError this package relies on, as an
// interface so it works with whichever generated SDK a command calls.
type apiError interface {
	error
	Status() int             // HTTP status, or 0 when no response was received
	Problem() map[string]any // decoded problem-details body, or nil
}

// UsageError reports a command that can't run as given: unknown flags, a
// missing argument, malformed input. Retrying the same command won't help.
type UsageError struct {
	Err error
}

func (e *UsageError) Error() string { return e.Err.Error() }
func (e *UsageError) Unwrap() error { return e.Err }

// Usagef returns a UsageError with a formatted message.
func Usagef(format string, args ...any) error {
	return &UsageError{Err: fmt.Errorf(format, args...)}
}

// AsUsage marks err as a UsageError, unless it already is one or is nil.
func AsUsage(err error) error {
	var usage *UsageError
	if err == nil || errors.As(err, &usage) {
		return err
	}
	return &UsageError{Err: err}
}

// CredentialsError reports that no usable credentials were configured.
type CredentialsError struct {
	Message string
}

func (e *CredentialsError) Error() string { return e.Message }

// ExitCode maps err to one of the Exit* codes.
func ExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	var usage *UsageError
	if errors.As(err, &usage) {
		return ExitUsage
	}
	var creds *CredentialsError
	if errors.As(err, &creds) {
		return ExitAuth
	}
	var apiErr apiError
	if errors.As(err, &apiErr) {
		switch status := apiErr.Status(); {
		case status == 0, status == 429, status >= 500:
			return ExitRetryable
		case status == 401, status == 403:
			return ExitAuth
		case status == 404:
			return ExitNotFound
		case status >= 400:
			return ExitRejected
		}
	}
	return ExitError
}
