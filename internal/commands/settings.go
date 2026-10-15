package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Settings come from the environment, else from a credentials file of
// KEY=value lines. Agents and other non-interactive shells often don't inherit
// the variables a user's shell profile exports, and the file lets them find
// the API key without it being passed as --token, where process listings and
// transcripts would show it.
//
//	HELO_API_KEY=...
//	HELO_CHANNEL_ID=...
//
// The file is read, never executed: lines may start with "export ", values may
// be quoted, and anything else (comments, other variables) is ignored.

// credentialsFileEnv names the variable that overrides the file's location.
const credentialsFileEnv = "HELO_CREDENTIALS_FILE"

// setting returns the first of names that is set, looking in the environment
// before the credentials file. A command reads only a few settings, so the
// file is simply read again for each.
func setting(names ...string) string {
	for _, name := range names {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}
	credentials := readCredentials()
	for _, name := range names {
		if value := credentials[name]; value != "" {
			return value
		}
	}
	return ""
}

// credentialsPath is $HELO_CREDENTIALS_FILE, else
// ${XDG_CONFIG_HOME:-~/.config}/helo/credentials.
func credentialsPath() string {
	if path := os.Getenv(credentialsFileEnv); path != "" {
		return path
	}
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "helo", "credentials")
}

func readCredentials() map[string]string {
	credentials := map[string]string{}
	path := credentialsPath()
	if path == "" {
		return credentials
	}
	file, err := os.Open(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "warning: can't read %s: %v\n", path, err)
		}
		return credentials
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimPrefix(line, "export ")
		name, value, ok := strings.Cut(line, "=")
		name = strings.TrimSpace(name)
		if !ok || !strings.HasPrefix(name, "HELO_") {
			continue
		}
		value = strings.TrimSpace(value)
		if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
			value = value[1 : len(value)-1]
		}
		credentials[name] = value
	}
	return credentials
}
