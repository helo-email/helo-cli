package commands

import (
	"github.com/spf13/cobra"

	"fmt"
	"github.com/helo-email/helo-sdk-go"
)

// NewBroadcastsCmd builds the cobra command group for Broadcasts.
func NewBroadcastsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broadcasts",
		Short: "Manage and track broadcast email campaigns.",
	}
	cmd.AddCommand(newBroadcastsListCmd())
	cmd.AddCommand(newBroadcastsRetrieveCmd())
	cmd.AddCommand(newBroadcastsListFailuresCmd())
	cmd.AddCommand(newBroadcastsListSuppressionsCmd())
	return cmd
}

func newBroadcastsListCmd() *cobra.Command {
	var channelID string
	var status string
	var subject string
	var from string
	var to string
	var limit int
	var offset int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List broadcasts",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			if status != "" {
				switch status {
				case "accepted", "processing", "completed", "canceled":
				default:
					return fmt.Errorf("invalid value %q for --status: must be one of accepted|processing|completed|canceled", status)
				}
			}

			params := &helo.BroadcastsListParams{}
			params.ChannelID = channelID
			params.Status = helo.BroadcastStatus(status)
			params.Subject = subject
			if from != "" {
				t, err := parseTimeFlag("from", from)
				if err != nil {
					return err
				}
				params.From = t
			}
			if to != "" {
				t, err := parseTimeFlag("to", to)
				if err != nil {
					return err
				}
				params.To = t
			}
			params.Limit = limit
			params.Offset = offset

			result, err := client.Broadcasts.List(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId (required)")
	_ = cmd.MarkFlagRequired("channel-id")
	cmd.Flags().StringVar(&status, "status", "", "status (one of: accepted|processing|completed|canceled)")
	_ = cmd.RegisterFlagCompletionFunc("status", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"accepted", "processing", "completed", "canceled"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&subject, "subject", "", "subject")
	cmd.Flags().StringVar(&from, "from", "", "from (RFC 3339 or a duration ago, e.g. 24h)")
	cmd.Flags().StringVar(&to, "to", "", "to (RFC 3339 or a duration ago, e.g. 24h)")
	cmd.Flags().IntVar(&limit, "limit", 0, "limit")
	cmd.Flags().IntVar(&offset, "offset", 0, "offset")

	return cmd
}

func newBroadcastsRetrieveCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "retrieve <id>",
		Short: "Retrieve a broadcast",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			result, err := client.Broadcasts.Retrieve(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}

func newBroadcastsListFailuresCmd() *cobra.Command {
	var limit int
	var offset int

	cmd := &cobra.Command{
		Use:   "list-failures <id>",
		Short: "List failed broadcast messages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			params := &helo.BroadcastsListFailuresParams{}
			params.Limit = limit
			params.Offset = offset

			result, err := client.Broadcasts.ListFailures(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 0, "limit")
	cmd.Flags().IntVar(&offset, "offset", 0, "offset")

	return cmd
}

func newBroadcastsListSuppressionsCmd() *cobra.Command {
	var limit int
	var offset int

	cmd := &cobra.Command{
		Use:   "list-suppressions <id>",
		Short: "List broadcast suppressed recipients",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			params := &helo.BroadcastsListSuppressionsParams{}
			params.Limit = limit
			params.Offset = offset

			result, err := client.Broadcasts.ListSuppressions(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 0, "limit")
	cmd.Flags().IntVar(&offset, "offset", 0, "offset")

	return cmd
}
