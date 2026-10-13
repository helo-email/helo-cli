package commands

import (
	"github.com/spf13/cobra"

	"fmt"
	"github.com/helo-email/helo-sdk-go"
)

// NewActivityCmd builds the cobra command group for Activity.
func NewActivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activity",
		Short: "Track and retrieve message activity, including messages, delivery and engagement events.",
	}
	cmd.AddCommand(newActivityListEventsCmd())
	cmd.AddCommand(newActivityListMessagesCmd())
	cmd.AddCommand(newActivityRetrieveMessageCmd())
	return cmd
}

func newActivityListEventsCmd() *cobra.Command {
	var channelID string
	var messageID string
	var after int64
	var startDate string
	var endDate string
	var limit int
	var recipient string
	var subject string
	var tags []string
	var mailType string
	var eventTypes []string

	cmd := &cobra.Command{
		Use:   "list-events",
		Short: "List activity events",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			if mailType != "" {
				switch mailType {
				case "transactional", "broadcast":
				default:
					return fmt.Errorf("invalid value %q for --mail-type: must be one of transactional|broadcast", mailType)
				}
			}
			for _, v := range eventTypes {
				switch v {
				case "accepted", "processed", "delivered", "bounced", "opened", "clicked", "complained", "unsubscribed", "resubscribed":
				default:
					return fmt.Errorf("invalid value %q for --event-types: must be one of accepted|processed|delivered|bounced|opened|clicked|complained|unsubscribed|resubscribed", v)
				}
			}

			params := &helo.ActivityListEventsParams{}
			params.ChannelID = channelID
			params.MessageID = messageID
			params.After = after
			if startDate != "" {
				t, err := parseTimeFlag("start-date", startDate)
				if err != nil {
					return err
				}
				params.StartDate = t
			}
			if endDate != "" {
				t, err := parseTimeFlag("end-date", endDate)
				if err != nil {
					return err
				}
				params.EndDate = t
			}
			params.Limit = limit
			params.Recipient = recipient
			params.Subject = subject
			params.Tags = tags
			params.MailType = mailType
			for _, v := range eventTypes {
				params.EventTypes = append(params.EventTypes, helo.EventType(v))
			}

			result, err := client.Activity.ListEvents(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&messageID, "message-id", "", "messageId")
	cmd.Flags().Int64Var(&after, "after", 0, "after")
	cmd.Flags().StringVar(&startDate, "start-date", "", "startDate (RFC 3339 or a duration ago, e.g. 24h)")
	cmd.Flags().StringVar(&endDate, "end-date", "", "endDate (RFC 3339 or a duration ago, e.g. 24h)")
	cmd.Flags().IntVar(&limit, "limit", 0, "limit")
	cmd.Flags().StringVar(&recipient, "recipient", "", "recipient")
	cmd.Flags().StringVar(&subject, "subject", "", "subject")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags")
	cmd.Flags().StringVar(&mailType, "mail-type", "", "mailType (one of: transactional|broadcast)")
	_ = cmd.RegisterFlagCompletionFunc("mail-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"transactional", "broadcast"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringSliceVar(&eventTypes, "event-types", nil, "eventTypes (one of: accepted|processed|delivered|bounced|opened|clicked|complained|unsubscribed|resubscribed)")
	_ = cmd.RegisterFlagCompletionFunc("event-types", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"accepted", "processed", "delivered", "bounced", "opened", "clicked", "complained", "unsubscribed", "resubscribed"}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func newActivityListMessagesCmd() *cobra.Command {
	var channelID string
	var after int64
	var startDate string
	var endDate string
	var limit int
	var recipient string
	var subject string
	var tags []string
	var mailType string
	var status string

	cmd := &cobra.Command{
		Use:   "list-messages",
		Short: "List messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			if mailType != "" {
				switch mailType {
				case "transactional", "broadcast":
				default:
					return fmt.Errorf("invalid value %q for --mail-type: must be one of transactional|broadcast", mailType)
				}
			}
			if status != "" {
				switch status {
				case "queued", "sent":
				default:
					return fmt.Errorf("invalid value %q for --status: must be one of queued|sent", status)
				}
			}

			params := &helo.ActivityListMessagesParams{}
			params.ChannelID = channelID
			params.After = after
			if startDate != "" {
				t, err := parseTimeFlag("start-date", startDate)
				if err != nil {
					return err
				}
				params.StartDate = t
			}
			if endDate != "" {
				t, err := parseTimeFlag("end-date", endDate)
				if err != nil {
					return err
				}
				params.EndDate = t
			}
			params.Limit = limit
			params.Recipient = recipient
			params.Subject = subject
			params.Tags = tags
			params.MailType = mailType
			params.Status = helo.MessageStatus(status)

			result, err := client.Activity.ListMessages(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().Int64Var(&after, "after", 0, "after")
	cmd.Flags().StringVar(&startDate, "start-date", "", "startDate (RFC 3339 or a duration ago, e.g. 24h)")
	cmd.Flags().StringVar(&endDate, "end-date", "", "endDate (RFC 3339 or a duration ago, e.g. 24h)")
	cmd.Flags().IntVar(&limit, "limit", 0, "limit")
	cmd.Flags().StringVar(&recipient, "recipient", "", "recipient")
	cmd.Flags().StringVar(&subject, "subject", "", "subject")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags")
	cmd.Flags().StringVar(&mailType, "mail-type", "", "mailType (one of: transactional|broadcast)")
	_ = cmd.RegisterFlagCompletionFunc("mail-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"transactional", "broadcast"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&status, "status", "", "status (one of: queued|sent)")
	_ = cmd.RegisterFlagCompletionFunc("status", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"queued", "sent"}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func newActivityRetrieveMessageCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "retrieve-message <id>",
		Short: "Retrieve message details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			result, err := client.Activity.RetrieveMessage(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}
