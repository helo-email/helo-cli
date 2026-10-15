# Changelog

## 0.1.2

- `helo` can now be installed with Homebrew (`brew install helo-email/tap/helo`) or with an install script
  (`curl -fsSL https://raw.githubusercontent.com/helo-email/helo-cli/main/install.sh | sh`), which verifies
  the download against the release's `checksums.txt`. Both avoid macOS's Gatekeeper warning for
  browser-downloaded binaries.

## 0.1.1

- Releases now include prebuilt binaries for macOS, Linux and Windows on amd64 and arm64, with a
  `checksums.txt`, so `helo` can be installed without Go: see the README. `helo --version` reports the
  release in these builds too.

## 0.1.0

- First release. `helo` has a command for every operation in the Helo API: channels, domains, sending,
  broadcasts, activity, statistics, suppressions and webhooks.
- Commands with a request body take it whole as JSON with `--data`, as well as field by field; flags
  override fields from `--data`. JSON values can be read from a file (`@path`) or stdin (`@-`), and
  unknown fields are rejected rather than ignored.
- Results are JSON on stdout. Errors go to stderr, as JSON with the API's error code, request ID and field
  errors when `--json` is set or stderr isn't a terminal, and the exit code says what kind of failure it
  was: 2 usage, 3 credentials, 4 not found, 5 rejected request, 6 worth retrying.
- `--dry-run` prints the request a command would send, with the API key redacted, without sending it.
- The API key is read from `--token`, `$HELO_API_KEY` or `$HELO_TOKEN`, and `--channel-id` defaults to
  `$HELO_CHANNEL_ID`.
- `--clear-channel-id` and `--clear-additional-headers` on `helo webhooks update` (or a `null` in
  `--data`) remove a webhook's channel scope or custom headers.
