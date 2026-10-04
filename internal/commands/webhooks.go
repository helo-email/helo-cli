package commands

import (
	"github.com/spf13/cobra"

	"encoding/json"
	"fmt"
	"github.com/helo-email/helo-sdk-go"
)

// NewWebhooksCmd builds the cobra command group for Webhooks.
func NewWebhooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhooks",
		Short: "Create and manage webhooks for event notifications.",
	}
	cmd.AddCommand(newWebhooksListCmd())
	cmd.AddCommand(newWebhooksCreateCmd())
	cmd.AddCommand(newWebhooksRetrieveCmd())
	cmd.AddCommand(newWebhooksUpdateCmd())
	cmd.AddCommand(newWebhooksDeleteCmd())
	cmd.AddCommand(newWebhooksRegenerateSigningKeyCmd())
	return cmd
}

func newWebhooksListCmd() *cobra.Command {
	var limit int
	var offset int
	var channelIds []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all webhooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}

			params := &helo.WebhooksListParams{}
			params.Limit = limit
			params.Offset = offset
			params.ChannelIds = channelIds

			result, err := client.Webhooks.List(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 0, "limit")
	cmd.Flags().IntVar(&offset, "offset", 0, "offset")
	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "channelIds")

	return cmd
}

func newWebhooksCreateCmd() *cobra.Command {
	var url string
	var events []string
	var channelID string
	var additionalHeaders string
	var enabled bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			for _, v := range events {
				switch v {
				case "message-accepted", "message-processed", "email-delivered", "email-bounced", "email-opened", "link-clicked", "recipient-complained", "recipient-unsubscribed", "recipient-resubscribed", "domain-key-verified", "domain-key-verification-failed", "return-path-domain-verified", "return-path-domain-verification-failed":
				default:
					return fmt.Errorf("invalid value %q for --events: must be one of message-accepted|message-processed|email-delivered|email-bounced|email-opened|link-clicked|recipient-complained|recipient-unsubscribed|recipient-resubscribed|domain-key-verified|domain-key-verification-failed|return-path-domain-verified|return-path-domain-verification-failed", v)
				}
			}

			params := &helo.CreateWebhookRequest{}
			params.URL = url
			if url == "" {
				return fmt.Errorf("--url is required")
			}
			for _, v := range events {
				params.Events = append(params.Events, helo.WebhookEvent(v))
			}
			params.ChannelID = channelID
			if additionalHeaders != "" {
				var v []helo.WebhookHeader
				if err := json.Unmarshal([]byte(additionalHeaders), &v); err != nil {
					return fmt.Errorf("invalid JSON for --additional-headers: %w", err)
				}
				params.AdditionalHeaders = v
			}
			params.Enabled = enabled

			result, err := client.Webhooks.Create(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "url (required)")
	_ = cmd.MarkFlagRequired("url")
	cmd.Flags().StringSliceVar(&events, "events", nil, "events (required, one of: message-accepted|message-processed|email-delivered|email-bounced|email-opened|link-clicked|recipient-complained|recipient-unsubscribed|recipient-resubscribed|domain-key-verified|domain-key-verification-failed|return-path-domain-verified|return-path-domain-verification-failed)")
	_ = cmd.MarkFlagRequired("events")
	_ = cmd.RegisterFlagCompletionFunc("events", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"message-accepted", "message-processed", "email-delivered", "email-bounced", "email-opened", "link-clicked", "recipient-complained", "recipient-unsubscribed", "recipient-resubscribed", "domain-key-verified", "domain-key-verification-failed", "return-path-domain-verified", "return-path-domain-verification-failed"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&additionalHeaders, "additional-headers", "", "additionalHeaders (JSON)")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "enabled")

	return cmd
}

func newWebhooksRetrieveCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "retrieve <id>",
		Short: "Retrieve a webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			result, err := client.Webhooks.Retrieve(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}

func newWebhooksUpdateCmd() *cobra.Command {
	var url string
	var events []string
	var channelID string
	var additionalHeaders string
	var enabled bool

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			for _, v := range events {
				switch v {
				case "message-accepted", "message-processed", "email-delivered", "email-bounced", "email-opened", "link-clicked", "recipient-complained", "recipient-unsubscribed", "recipient-resubscribed", "domain-key-verified", "domain-key-verification-failed", "return-path-domain-verified", "return-path-domain-verification-failed":
				default:
					return fmt.Errorf("invalid value %q for --events: must be one of message-accepted|message-processed|email-delivered|email-bounced|email-opened|link-clicked|recipient-complained|recipient-unsubscribed|recipient-resubscribed|domain-key-verified|domain-key-verification-failed|return-path-domain-verified|return-path-domain-verification-failed", v)
				}
			}
			id := args[0]

			params := &helo.UpdateWebhookRequest{}
			params.URL = url
			for _, v := range events {
				params.Events = append(params.Events, helo.WebhookEvent(v))
			}
			params.ChannelID = channelID
			if additionalHeaders != "" {
				var v []helo.WebhookHeader
				if err := json.Unmarshal([]byte(additionalHeaders), &v); err != nil {
					return fmt.Errorf("invalid JSON for --additional-headers: %w", err)
				}
				params.AdditionalHeaders = v
			}
			params.Enabled = enabled

			result, err := client.Webhooks.Update(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "url")
	cmd.Flags().StringSliceVar(&events, "events", nil, "events (one of: message-accepted|message-processed|email-delivered|email-bounced|email-opened|link-clicked|recipient-complained|recipient-unsubscribed|recipient-resubscribed|domain-key-verified|domain-key-verification-failed|return-path-domain-verified|return-path-domain-verification-failed)")
	_ = cmd.RegisterFlagCompletionFunc("events", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"message-accepted", "message-processed", "email-delivered", "email-bounced", "email-opened", "link-clicked", "recipient-complained", "recipient-unsubscribed", "recipient-resubscribed", "domain-key-verified", "domain-key-verification-failed", "return-path-domain-verified", "return-path-domain-verification-failed"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&additionalHeaders, "additional-headers", "", "additionalHeaders (JSON)")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "enabled")

	return cmd
}

func newWebhooksDeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			if err := client.Webhooks.Delete(cmd.Context(), id); err != nil {
				return err
			}
			return emit(map[string]string{"status": "ok"})
		},
	}

	return cmd
}

func newWebhooksRegenerateSigningKeyCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "regenerate-signing-key <id>",
		Short: "Regenerate webhook signing key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			result, err := client.Webhooks.RegenerateSigningKey(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}
