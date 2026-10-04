package commands

import (
	"github.com/spf13/cobra"

	"fmt"
	"github.com/helo-email/helo-sdk-go"
	"time"
)

// NewStatisticsCmd builds the cobra command group for Statistics.
func NewStatisticsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "statistics",
		Short: "Access activity statistics.",
	}
	cmd.AddCommand(newStatisticsRetrieveHourlyCmd())
	cmd.AddCommand(newStatisticsRetrieveDailyCmd())
	cmd.AddCommand(newStatisticsRetrieveTotalsCmd())
	return cmd
}

func newStatisticsRetrieveHourlyCmd() *cobra.Command {
	var channelID string
	var from string
	var to string
	var tags []string

	cmd := &cobra.Command{
		Use:   "retrieve-hourly",
		Short: "Retrieve hourly statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}

			params := &helo.StatisticsRetrieveHourlyParams{}
			params.ChannelID = channelID
			if from != "" {
				t, err := time.Parse(time.RFC3339, from)
				if err != nil {
					return fmt.Errorf("invalid timestamp for --from: %w", err)
				}
				params.From = t
			}
			if to != "" {
				t, err := time.Parse(time.RFC3339, to)
				if err != nil {
					return fmt.Errorf("invalid timestamp for --to: %w", err)
				}
				params.To = t
			}
			params.Tags = tags

			result, err := client.Statistics.RetrieveHourly(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&from, "from", "", "from (required)")
	_ = cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&to, "to", "", "to (required)")
	_ = cmd.MarkFlagRequired("to")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags")

	return cmd
}

func newStatisticsRetrieveDailyCmd() *cobra.Command {
	var channelID string
	var from string
	var to string
	var tags []string
	var timezone string

	cmd := &cobra.Command{
		Use:   "retrieve-daily",
		Short: "Retrieve daily statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			if from != "" {
				if _, err := time.Parse(time.DateOnly, from); err != nil {
					return fmt.Errorf("invalid value %q for --from: must be YYYY-MM-DD (e.g. 2025-02-27)", from)
				}
			}
			if to != "" {
				if _, err := time.Parse(time.DateOnly, to); err != nil {
					return fmt.Errorf("invalid value %q for --to: must be YYYY-MM-DD (e.g. 2025-02-27)", to)
				}
			}

			params := &helo.StatisticsRetrieveDailyParams{}
			params.ChannelID = channelID
			params.From = from
			params.To = to
			params.Tags = tags
			params.Timezone = timezone

			result, err := client.Statistics.RetrieveDaily(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&from, "from", "", "from (required)")
	_ = cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&to, "to", "", "to (required)")
	_ = cmd.MarkFlagRequired("to")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags")
	cmd.Flags().StringVar(&timezone, "timezone", "", "timezone (required)")
	_ = cmd.MarkFlagRequired("timezone")

	return cmd
}

func newStatisticsRetrieveTotalsCmd() *cobra.Command {
	var channelID string
	var from string
	var to string
	var tags []string

	cmd := &cobra.Command{
		Use:   "retrieve-totals",
		Short: "Retrieve all time statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}

			params := &helo.StatisticsRetrieveTotalsParams{}
			params.ChannelID = channelID
			if from != "" {
				t, err := time.Parse(time.RFC3339, from)
				if err != nil {
					return fmt.Errorf("invalid timestamp for --from: %w", err)
				}
				params.From = t
			}
			if to != "" {
				t, err := time.Parse(time.RFC3339, to)
				if err != nil {
					return fmt.Errorf("invalid timestamp for --to: %w", err)
				}
				params.To = t
			}
			params.Tags = tags

			result, err := client.Statistics.RetrieveTotals(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&from, "from", "", "from (required)")
	_ = cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&to, "to", "", "to (required)")
	_ = cmd.MarkFlagRequired("to")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags")

	return cmd
}
