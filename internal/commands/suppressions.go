package commands

import (
	"github.com/spf13/cobra"

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
		Long:  "Retrieves a list of suppressed email addresses for a channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mailType != "" {
				switch mailType {
				case "transactional", "broadcast":
				default:
					return usageErrorf("invalid value %q for --mail-type: must be one of transactional|broadcast", mailType)
				}
			}
			if reason != "" {
				switch reason {
				case "bounce", "complaint", "unsubscribe", "manual":
				default:
					return usageErrorf("invalid value %q for --reason: must be one of bounce|complaint|unsubscribe|manual", reason)
				}
			}

			params := &helo.SuppressionsListParams{}
			params.ChannelID = channelID
			params.MailType = helo.MailType(mailType)
			params.Reason = helo.SuppressionReason(reason)
			params.Email = email
			params.Limit = limit
			params.Offset = offset

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Suppressions.List(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId (required)")
	_ = cmd.MarkFlagRequired("channel-id")
	cmd.Flags().StringVar(&mailType, "mail-type", "", "mailType (required; one of: transactional|broadcast)")
	_ = cmd.MarkFlagRequired("mail-type")
	_ = cmd.RegisterFlagCompletionFunc("mail-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"transactional", "broadcast"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&reason, "reason", "", "reason (one of: bounce|complaint|unsubscribe|manual)")
	_ = cmd.RegisterFlagCompletionFunc("reason", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"bounce", "complaint", "unsubscribe", "manual"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&email, "email", "", "email")
	cmd.Flags().IntVar(&limit, "limit", 0, "limit (min 1; max 500)")
	cmd.Flags().IntVar(&offset, "offset", 0, "offset (min 0)")

	return cmd
}

func newSuppressionsCreateCmd() *cobra.Command {
	var channelID string
	var mailType string
	var emails []string
	var requestData string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create suppressions",
		Long:  "Adds email addresses to the suppression list to prevent future sends.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mailType != "" {
				switch mailType {
				case "transactional", "broadcast":
				default:
					return usageErrorf("invalid value %q for --mail-type: must be one of transactional|broadcast", mailType)
				}
			}

			params := &helo.CreateSuppressionsRequest{}
			dataFields, err := decodeBody(requestData, params)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("channel-id") {
				params.ChannelID = channelID
			}
			if !cmd.Flags().Changed("channel-id") && !dataFields.set["channelId"] {
				return usageErrorf("--channel-id is required (or set \"channelId\" in --data)")
			}
			if cmd.Flags().Changed("mail-type") {
				params.MailType = helo.MailType(mailType)
			}
			if !cmd.Flags().Changed("mail-type") && !dataFields.set["mailType"] {
				return usageErrorf("--mail-type is required (or set \"mailType\" in --data)")
			}
			if cmd.Flags().Changed("emails") {
				params.Emails = emails
			}
			if !cmd.Flags().Changed("emails") && !dataFields.set["emails"] {
				return usageErrorf("--emails is required (or set \"emails\" in --data)")
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Suppressions.Create(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId (required)")
	cmd.Flags().StringVar(&mailType, "mail-type", "", "mailType (required; one of: transactional|broadcast)")
	_ = cmd.RegisterFlagCompletionFunc("mail-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"transactional", "broadcast"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringSliceVar(&emails, "emails", nil, "emails (required; comma-separated)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}

func newSuppressionsRemoveCmd() *cobra.Command {
	var channelID string
	var mailType string
	var emails []string
	var requestData string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove suppressions",
		Long:  "Removes email addresses from the suppression list to allow future sends.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mailType != "" {
				switch mailType {
				case "transactional", "broadcast":
				default:
					return usageErrorf("invalid value %q for --mail-type: must be one of transactional|broadcast", mailType)
				}
			}

			params := &helo.RemoveSuppressionsRequest{}
			dataFields, err := decodeBody(requestData, params)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("channel-id") {
				params.ChannelID = channelID
			}
			if !cmd.Flags().Changed("channel-id") && !dataFields.set["channelId"] {
				return usageErrorf("--channel-id is required (or set \"channelId\" in --data)")
			}
			if cmd.Flags().Changed("mail-type") {
				params.MailType = helo.MailType(mailType)
			}
			if !cmd.Flags().Changed("mail-type") && !dataFields.set["mailType"] {
				return usageErrorf("--mail-type is required (or set \"mailType\" in --data)")
			}
			if cmd.Flags().Changed("emails") {
				params.Emails = emails
			}
			if !cmd.Flags().Changed("emails") && !dataFields.set["emails"] {
				return usageErrorf("--emails is required (or set \"emails\" in --data)")
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Suppressions.Remove(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId (required)")
	cmd.Flags().StringVar(&mailType, "mail-type", "", "mailType (required; one of: transactional|broadcast)")
	_ = cmd.RegisterFlagCompletionFunc("mail-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"transactional", "broadcast"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringSliceVar(&emails, "emails", nil, "emails (required; comma-separated)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}
