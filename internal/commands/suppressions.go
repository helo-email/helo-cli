package commands

import (
	"github.com/spf13/cobra"

	"fmt"
	"github.com/helo-email/helo-sdk-go"
)

// NewSuppressionsCmd builds the cobra command group for Suppressions.
func NewSuppressionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suppressions",
		Short: "Manage Suppressions",
	}
	cmd.AddCommand(newSuppressionsListCmd())
	cmd.AddCommand(newSuppressionsCreateCmd())
	cmd.AddCommand(newSuppressionsRemoveCmd())
	return cmd
}

func newSuppressionsListCmd() *cobra.Command {
	var channelID string
	var mailType string
	var reason string
	var email string
	var limit int
	var offset int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List suppressions",
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
			if reason != "" {
				switch reason {
				case "bounce", "complaint", "unsubscribe", "manual":
				default:
					return fmt.Errorf("invalid value %q for --reason: must be one of bounce|complaint|unsubscribe|manual", reason)
				}
			}

			params := &helo.SuppressionsListParams{}
			params.ChannelID = channelID
			params.MailType = helo.MailType(mailType)
			params.Reason = helo.SuppressionReason(reason)
			params.Email = email
			params.Limit = limit
			params.Offset = offset

			result, err := client.Suppressions.List(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId (required)")
	_ = cmd.MarkFlagRequired("channel-id")
	cmd.Flags().StringVar(&mailType, "mail-type", "", "mailType (required, one of: transactional|broadcast)")
	_ = cmd.MarkFlagRequired("mail-type")
	_ = cmd.RegisterFlagCompletionFunc("mail-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"transactional", "broadcast"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&reason, "reason", "", "reason (one of: bounce|complaint|unsubscribe|manual)")
	_ = cmd.RegisterFlagCompletionFunc("reason", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"bounce", "complaint", "unsubscribe", "manual"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&email, "email", "", "email")
	cmd.Flags().IntVar(&limit, "limit", 0, "limit")
	cmd.Flags().IntVar(&offset, "offset", 0, "offset")

	return cmd
}

func newSuppressionsCreateCmd() *cobra.Command {
	var channelID string
	var mailType string
	var emails []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create suppressions",
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

			params := &helo.CreateSuppressionsRequest{}
			params.ChannelID = channelID
			if channelID == "" {
				return fmt.Errorf("--channel-id is required")
			}
			params.MailType = helo.MailType(mailType)
			if mailType == "" {
				return fmt.Errorf("--mail-type is required")
			}
			params.Emails = emails

			result, err := client.Suppressions.Create(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId (required)")
	_ = cmd.MarkFlagRequired("channel-id")
	cmd.Flags().StringVar(&mailType, "mail-type", "", "mailType (required, one of: transactional|broadcast)")
	_ = cmd.MarkFlagRequired("mail-type")
	_ = cmd.RegisterFlagCompletionFunc("mail-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"transactional", "broadcast"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringSliceVar(&emails, "emails", nil, "emails (required)")
	_ = cmd.MarkFlagRequired("emails")

	return cmd
}

func newSuppressionsRemoveCmd() *cobra.Command {
	var channelID string
	var mailType string
	var emails []string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove suppressions",
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

			params := &helo.RemoveSuppressionsRequest{}
			params.ChannelID = channelID
			if channelID == "" {
				return fmt.Errorf("--channel-id is required")
			}
			params.MailType = helo.MailType(mailType)
			if mailType == "" {
				return fmt.Errorf("--mail-type is required")
			}
			params.Emails = emails

			result, err := client.Suppressions.Remove(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId (required)")
	_ = cmd.MarkFlagRequired("channel-id")
	cmd.Flags().StringVar(&mailType, "mail-type", "", "mailType (required, one of: transactional|broadcast)")
	_ = cmd.MarkFlagRequired("mail-type")
	_ = cmd.RegisterFlagCompletionFunc("mail-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"transactional", "broadcast"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringSliceVar(&emails, "emails", nil, "emails (required)")
	_ = cmd.MarkFlagRequired("emails")

	return cmd
}
