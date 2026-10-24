// Package cli wires up the root cobra command and global flags.
package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/helo-email/helo-cli/internal/commands"
	"github.com/helo-email/helo-cli/internal/output"
)

// GlobalFlags collects flags shared by every subcommand.
type GlobalFlags struct {
	Token   string
	BaseURL string
	JSON    bool
	Quiet   bool
	DryRun  bool
}

type globalFlagsKey struct{}

// WithGlobalFlags stores GlobalFlags on a context.
func WithGlobalFlags(ctx context.Context, flags *GlobalFlags) context.Context {
	return context.WithValue(ctx, globalFlagsKey{}, flags)
}

// FlagsFromContext retrieves the GlobalFlags previously stored on a context.
func FlagsFromContext(ctx context.Context) *GlobalFlags {
	if f, ok := ctx.Value(globalFlagsKey{}).(*GlobalFlags); ok {
		return f
	}
	return &GlobalFlags{}
}

const longHelp = "Helo Email API (https://helohq.com)" + `

Authentication:
  Pass an API key with --token, or set $HELO_API_KEY ($HELO_TOKEN also works).
  Commands that send through a channel use --channel-id, else $HELO_CHANNEL_ID.
  Variables not in the environment are read from ~/.config/helo/credentials, a
  file of HELO_API_KEY=... lines ($HELO_CREDENTIALS_FILE moves it).

Input:
  Commands that take a request body accept the whole body as JSON with --data,
  and each field as its own flag; flags override fields from --data. JSON flag
  values can be given inline, as @path to read a file, or as @- to read stdin.
  Unknown JSON fields are rejected rather than ignored.

Output:
  Results are written to stdout as {"ok": true, "data": ...}; --quiet prints just
  the data. Errors go to stderr, as JSON ({"ok": false, "error": ..., "code": ...,
  "requestId": ..., "exitCode": ...}) when --json is set or stderr isn't a terminal.
  --dry-run prints the request a command would send, with the API key redacted,
  without sending it.

Exit codes:
  0  success
  1  other error
  2  usage error: bad flags, arguments or input
  3  missing, invalid or insufficient credentials (401, 403)
  4  not found (404)
  5  request rejected (400, 409, 422): fix it before retrying
  6  rate limited (429), server error (5xx) or no response: retrying may help`

// NewRootCmd builds the root command for helo.
func NewRootCmd() *cobra.Command {
	flags := &GlobalFlags{}

	cmd := &cobra.Command{
		Use:           "helo",
		Short:         "Command-line interface for Helo API",
		Long:          longHelp,
		Version:       commands.Version(),
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetContext(WithGlobalFlags(cmd.Context(), flags))
			commands.SetGlobalFlags(commands.GlobalFlags{
				Token:   flags.Token,
				BaseURL: flags.BaseURL,
				JSON:    flags.JSON,
				Quiet:   flags.Quiet,
				DryRun:  flags.DryRun,
			})
			return nil
		},
	}
	cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return output.AsUsage(err)
	})

	cmd.PersistentFlags().StringVar(&flags.Token, "token", "", "API key (overrides $HELO_API_KEY and $HELO_TOKEN)")
	cmd.PersistentFlags().StringVar(&flags.BaseURL, "base-url", "", "API base URL override (default $HELO_BASE_URL or $HELO_API_URL)")
	cmd.PersistentFlags().BoolVar(&flags.JSON, "json", false, "Report errors as JSON (the default when stderr isn't a terminal)")
	cmd.PersistentFlags().BoolVar(&flags.Quiet, "quiet", false, "Output data only, no envelope")
	cmd.PersistentFlags().BoolVar(&flags.DryRun, "dry-run", false, "Print the request instead of sending it")

	cmd.AddCommand(
		commands.NewChannelsCmd(),
		commands.NewActivityCmd(),
		commands.NewDomainsCmd(),
		commands.NewSendingCmd(),
		commands.NewBroadcastsCmd(),
		commands.NewStatisticsCmd(),
		commands.NewSuppressionsCmd(),
		commands.NewTemplatesCmd(),
		commands.NewWebhooksCmd(),
	)

	return cmd
}

// Execute runs the CLI and returns its exit code (see longHelp).
func Execute(args []string) int {
	cmd := NewRootCmd()
	cmd.SetArgs(args)
	started := markStarted(cmd)

	err := cmd.Execute()
	switch {
	case err == nil, commands.IsDryRun(err):
		return output.ExitOK
	case !*started:
		// cobra rejected the command line (an unknown command, a missing
		// argument or required flag) before any command ran.
		err = output.AsUsage(err)
	}
	output.PrintError(err, wantJSONErrors(args))
	return output.ExitCode(err)
}

// markStarted wraps every command's RunE to record that one began running,
// which separates cobra's own validation errors from a command's failures.
func markStarted(root *cobra.Command) *bool {
	started := new(bool)
	var wrap func(*cobra.Command)
	wrap = func(c *cobra.Command) {
		if run := c.RunE; run != nil {
			c.RunE = func(cmd *cobra.Command, args []string) error {
				*started = true
				return run(cmd, args)
			}
		}
		for _, child := range c.Commands() {
			wrap(child)
		}
	}
	wrap(root)
	return started
}

// wantJSONErrors reports whether errors should be JSON: when asked for, or when
// a program rather than a person is reading stderr. It reads args itself since
// flags may not have been parsed when the error happened.
func wantJSONErrors(args []string) bool {
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "--json" || arg == "--json=true" {
			return true
		}
	}
	info, err := os.Stderr.Stat()
	return err != nil || info.Mode()&os.ModeCharDevice == 0
}
