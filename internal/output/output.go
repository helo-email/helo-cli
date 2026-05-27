// Package output renders command results as either a styled JSON envelope
// or raw JSON for piping into other tools.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Envelope wraps successful results.
type Envelope struct {
	OK   bool        `json:"ok"`
	Data interface{} `json:"data,omitempty"`
}

// ErrorEnvelope wraps failure results.
type ErrorEnvelope struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// PrintOptions controls how PrintJSON renders results.
type PrintOptions struct {
	Quiet bool // when true, emit only the raw data without an envelope
}

// PrintJSON writes data to w as pretty-printed JSON.
func PrintJSON(w io.Writer, data any, opts PrintOptions) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if opts.Quiet {
		if data == nil {
			return nil
		}
		return enc.Encode(data)
	}
	return enc.Encode(Envelope{OK: true, Data: data})
}

// PrintError writes an error to stderr as JSON when machine output is requested.
func PrintError(err error, asJSON bool) {
	if asJSON {
		enc := json.NewEncoder(os.Stderr)
		enc.SetIndent("", "  ")
		_ = enc.Encode(ErrorEnvelope{OK: false, Error: err.Error()})
		return
	}
	fmt.Fprintln(os.Stderr, "Error:", err)
}
