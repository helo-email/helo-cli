# helo CLI

Generated CLI for Helo API. Wraps the [helo](https://github.com/balvig/helo-sdk-go) Go SDK.

## Install

```bash
go install github.com/balvig/helo-cli/cmd/helo@latest
```

Or build from a local checkout:

```bash
go build -o helo ./cmd/helo
```

## Authentication

Export your access token:

```bash
export HELO_TOKEN="your-token"
```

Or pass it per invocation with `--token`.

## Usage

```bash
helo --help
helo channels --help
```

### Output

By default the CLI emits a JSON envelope:

```json
{
  "ok": true,
  "data": { ... }
}
```

Pass `--quiet` to print just the raw data (e.g. for piping into `jq`).

### Resources

- `helo channels` — Create and manage communication channels for organizing messages.
- `helo activity` — Track and retrieve message activity, including messages, delivery and engagement events.
- `helo domains` — Register, verify, and manage domains for email sending.
- `helo sending` — Send transactional and broadcast emails.
- `helo broadcasts` — Manage and track broadcast email campaigns.
- `helo statistics` — Access activity statistics.
- `helo suppressions` — Manage Suppressions
- `helo webhook-endpoints` — Create and manage webhook endpoints for event notifications.
