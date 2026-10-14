# helo CLI

CLI for Helo API. Wraps the [helo](https://github.com/helo-email/helo-sdk-go) Go SDK.

## Install

Download a prebuilt binary from the [latest release](https://github.com/helo-email/helo-cli/releases/latest), for macOS, Linux or
Windows on amd64 or arm64. On macOS or Linux:

```bash
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')
curl -fsSL "https://github.com/helo-email/helo-cli/releases/latest/download/helo_${os}_${arch}.tar.gz" \
  | tar -xz -C ~/.local/bin helo
```

That installs into `~/.local/bin`, which must exist and be on your `PATH`; use another directory if you
prefer. Each release also has a `checksums.txt` of SHA-256 sums. Or, with Go installed:

```bash
go install github.com/helo-email/helo-cli/cmd/helo@latest
```

`@latest` is the newest release; pin one with `@v0.1.1`, say. `helo --version` reports the
release and the SDK version it was built with. Each release's notes are in [CHANGELOG.md](CHANGELOG.md) and on
the repository's GitHub releases.


Or build from a local checkout:

```bash
go build -o helo ./cmd/helo
```

`go.mod` requires a released version of the SDK. A local build also uses a checkout of the SDK next to this
repository, through a `go.work` file the generator writes (and `.gitignore` keeps out of the repository).

## Authentication

Export your API key:

```bash
export HELO_API_KEY="your-api-key"
```

Or pass it per invocation with `--token`. `$HELO_TOKEN` is also read, after `$HELO_API_KEY`.

Commands that send through a channel take `--channel-id`, which defaults to `$HELO_CHANNEL_ID`.

## Usage

```bash
helo --help
helo channels --help
```

### Input

Commands that take a request body accept each field as a flag, and the whole body as JSON with `--data`.
Flags override fields from `--data`. JSON values can be given inline, as `@path` to read a file, or as `@-`
to read standard input:

```bash
helo <command> --data @request.json
helo <command> --data @request.json --some-field override
generate-request | helo <command> --data @-
```

Unknown JSON fields are rejected rather than silently dropped. Each flag's help shows the JSON it takes.

`--dry-run` prints the request a command would send, with the API key redacted, without sending it.

### Output

Results are written to stdout as a JSON envelope:

```json
{
  "ok": true,
  "data": { ... }
}
```

Pass `--quiet` to print just the raw data (e.g. for piping into `jq`).

Errors are written to stderr: as text at a terminal, and as JSON when `--json` is set or stderr isn't a
terminal. The JSON carries the API's problem details when there are any:

```json
{
  "ok": false,
  "error": "One or more fields were invalid. (status 422, code validation_failed, request 0HNL29NEOHAS9:00000004)",
  "exitCode": 5,
  "status": 422,
  "code": "validation_failed",
  "detail": "One or more fields were invalid.",
  "requestId": "0HNL29NEOHAS9:00000004",
  "errors": { "/to/0/email": [{ "message": "...", "code": "..." }] }
}
```

### Exit codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | Other error |
| 2 | Usage error: bad flags, arguments or input |
| 3 | Missing, invalid or insufficient credentials (401, 403) |
| 4 | Not found (404) |
| 5 | Request rejected (400, 409, 422): fix it before retrying |
| 6 | Rate limited (429), server error (5xx) or no response: retrying may help |

### Scripts and agents

- Rely on the exit code, and on the JSON error on stderr, rather than parsing messages.
- Retry only exit code 6, and for requests that send something, pass the same `--idempotency-key` each
  time so a retry can't send twice.
- Build large or nested bodies as a JSON file and pass `--data @file`, rather than quoting JSON on the
  command line.
- Use `--dry-run` to show exactly what a command will send before running it.

### Resources

- `helo channels` — Create and manage communication channels for organizing messages.
- `helo activity` — Track and retrieve message activity, including messages, delivery and engagement events.
- `helo domains` — Register, verify, and manage domains for email sending.
- `helo sending` — Send transactional and broadcast emails.
- `helo broadcasts` — Manage and track broadcast email campaigns.
- `helo statistics` — Access activity statistics.
- `helo suppressions` — Manage Suppressions
- `helo webhooks` — Create and manage webhooks for event notifications.
