// Package cli wires up the root cobra command and global flags.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/balvig/helo-cli/internal/commands"
)

// GlobalFlags collects flags shared by every subcommand.
type GlobalFlags struct {
	Token   string
	BaseURL string
	JSON    bool
	Quiet   bool
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

// NewRootCmd builds the root command for helo.
func NewRootCmd() *cobra.Command {
	flags := &GlobalFlags{}

	cmd := &cobra.Command{
		Use:           "helo",
		Short:         "Command-line interface for Helo API",
		Long:          "Helo API",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetContext(WithGlobalFlags(cmd.Context(), flags))
			commands.SetGlobalFlags(commands.GlobalFlags{
				Token:   flags.Token,
				BaseURL: flags.BaseURL,
				JSON:    flags.JSON,
				Quiet:   flags.Quiet,
			})
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&flags.Token, "token", "", "API access token (overrides $HELO_TOKEN)")
	cmd.PersistentFlags().StringVar(&flags.BaseURL, "base-url", "", "API base URL override")
	cmd.PersistentFlags().BoolVar(&flags.JSON, "json", false, "Force JSON output")
	cmd.PersistentFlags().BoolVar(&flags.Quiet, "quiet", false, "Output data only, no envelope")

	cmd.AddCommand(
		commands.NewChannelsCmd(),
		commands.NewActivityCmd(),
		commands.NewDomainsCmd(),
		commands.NewSendingCmd(),
		commands.NewBroadcastsCmd(),
		commands.NewStatisticsCmd(),
		commands.NewSuppressionsCmd(),
		commands.NewWebhookEndpointsCmd(),
	)

	return cmd
}
