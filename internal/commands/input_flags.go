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

// bodyFields records the top-level fields --data set, and which it set to
// null, so a required field can come from either --data or its own flag, and a
// null can clear a field that accepts one.
type bodyFields struct {
	set  map[string]bool
	null map[string]bool
}

// decodeBody parses --data, a whole request body, into v.
func decodeBody(value string, v any) (bodyFields, error) {
	fields := bodyFields{set: map[string]bool{}, null: map[string]bool{}}
	if value == "" {
		return fields, nil
	}
	data, err := readInput("data", value)
	if err != nil {
		return fields, err
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fields, output.Usagef("--data must be a JSON object: %v", err)
	}
	if err := decodeJSON("data", data, v); err != nil {
		return fields, err
	}
	for name, value := range raw {
		fields.set[name] = true
		fields.null[name] = string(bytes.TrimSpace(value)) == "null"
	}
	return fields, nil
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
