package commands

import (
	"github.com/spf13/cobra"

	"github.com/helo-email/helo-sdk-go"
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
		Long:  "Fetches hourly aggregated statistics.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.StatisticsRetrieveHourlyParams{}
			params.ChannelID = channelID
			if from != "" {
				t, err := parseTimeFlag("from", from)
				if err != nil {
					return usageError(err)
				}
				params.From = t
			}
			if to != "" {
				t, err := parseTimeFlag("to", to)
				if err != nil {
					return usageError(err)
				}
				params.To = t
			}
			params.Tags = tags

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Statistics.RetrieveHourly(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&from, "from", "", "from (required; RFC 3339 or a duration ago, e.g. 24h)")
	_ = cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&to, "to", "", "to (required; RFC 3339 or a duration ago, e.g. 24h)")
	_ = cmd.MarkFlagRequired("to")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags (comma-separated)")

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
		Long:  "Fetches daily aggregated statistics.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if from != "" {
				d, err := parseDateFlag("from", from)
				if err != nil {
					return usageError(err)
				}
				from = d
			}
			if to != "" {
				d, err := parseDateFlag("to", to)
				if err != nil {
					return usageError(err)
				}
				to = d
			}

			params := &helo.StatisticsRetrieveDailyParams{}
			params.ChannelID = channelID
			params.From = from
			params.To = to
			params.Tags = tags
			params.Timezone = timezone

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Statistics.RetrieveDaily(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&from, "from", "", "from (required; YYYY-MM-DD or a duration ago, e.g. 7d)")
	_ = cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&to, "to", "", "to (required; YYYY-MM-DD or a duration ago, e.g. 7d)")
	_ = cmd.MarkFlagRequired("to")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags (comma-separated)")
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
		Long:  "Fetches cumulative statistics.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.StatisticsRetrieveTotalsParams{}
			params.ChannelID = channelID
			if from != "" {
				t, err := parseTimeFlag("from", from)
				if err != nil {
					return usageError(err)
				}
				params.From = t
			}
			if to != "" {
				t, err := parseTimeFlag("to", to)
				if err != nil {
					return usageError(err)
				}
				params.To = t
			}
			params.Tags = tags

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Statistics.RetrieveTotals(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&from, "from", "", "from (required; RFC 3339 or a duration ago, e.g. 24h)")
	_ = cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&to, "to", "", "to (required; RFC 3339 or a duration ago, e.g. 24h)")
	_ = cmd.MarkFlagRequired("to")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags (comma-separated)")

	return cmd
}
