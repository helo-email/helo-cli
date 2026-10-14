package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/helo-email/helo-cli/internal/output"
)

// JSON flags and --data take their value inline, or indirectly so large or
// awkward-to-quote input needn't pass through the shell: "@path" reads a file
// and "@-" reads standard input.

// channelIDEnv names the variable --channel-id falls back to.
const channelIDEnv = "HELO_CHANNEL_ID"

// stdinRead records that a flag consumed standard input, which can only be
// read once.
var stdinRead bool

func readInput(flag, value string) ([]byte, error) {
	switch {
	case value == "@-":
		if stdinRead {
			return nil, output.Usagef("--%s: standard input was already read by another flag", flag)
		}
		stdinRead = true
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, output.Usagef("--%s: reading standard input: %v", flag, err)
		}
		return data, nil
	case strings.HasPrefix(value, "@"):
		data, err := os.ReadFile(value[1:])
		if err != nil {
			return nil, output.Usagef("--%s: %v", flag, err)
		}
		return data, nil
	}
	return []byte(value), nil
}

// decodeJSONFlag parses a JSON flag value into v. Fields v doesn't have are
// rejected, so a misspelled key fails rather than being silently dropped.
func decodeJSONFlag(flag, value string, v any) error {
	data, err := readInput(flag, value)
	if err != nil {
		return err
	}
	return decodeJSON(flag, data, v)
}

func decodeJSON(flag string, data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return output.Usagef("invalid JSON for --%s: %v", flag, err)
	}
	if dec.More() {
		return output.Usagef("invalid JSON for --%s: unexpected data after the JSON value", flag)
	}
	return nil
}

// decodeBody parses --data, a whole request body, into v and reports which
// top-level fields it set, so a required field can come from either --data
// or its own flag.
func decodeBody(value string, v any) (map[string]bool, error) {
	if value == "" {
		return map[string]bool{}, nil
	}
	data, err := readInput("data", value)
	if err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, output.Usagef("--data must be a JSON object: %v", err)
	}
	if err := decodeJSON("data", data, v); err != nil {
		return nil, err
	}
	set := make(map[string]bool, len(fields))
	for name := range fields {
		set[name] = true
	}
	return set, nil
}

// envDefault returns value, or the environment variable name when value is empty.
func envDefault(value, name string) string {
	if value != "" {
		return value
	}
	return os.Getenv(name)
}

func usageErrorf(format string, args ...any) error {
	return output.Usagef(format, args...)
}

func usageError(err error) error {
	return output.AsUsage(err)
}
