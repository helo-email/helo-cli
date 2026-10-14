package commands

import (
	"github.com/spf13/cobra"

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
		Long:  "Retrieves all webhooks configured for the account.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.WebhooksListParams{}
			params.Limit = limit
			params.Offset = offset
			params.ChannelIds = channelIds

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Webhooks.List(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 0, "limit")
	cmd.Flags().IntVar(&offset, "offset", 0, "offset")
	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "channelIds (comma-separated)")

	return cmd
}

func newWebhooksCreateCmd() *cobra.Command {
	var url string
	var events []string
	var channelID string
	var additionalHeaders string
	var enabled bool
	var requestData string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a webhook",
		Long:  "Registers a new webhook to receive event notifications.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, v := range events {
				switch v {
				case "message-accepted", "message-processed", "email-delivered", "email-bounced", "email-opened", "link-clicked", "recipient-complained", "recipient-unsubscribed", "recipient-resubscribed", "domain-key-verified", "domain-key-verification-failed", "return-path-domain-verified", "return-path-domain-verification-failed":
				default:
					return usageErrorf("invalid value %q for --events: must be one of message-accepted|message-processed|email-delivered|email-bounced|email-opened|link-clicked|recipient-complained|recipient-unsubscribed|recipient-resubscribed|domain-key-verified|domain-key-verification-failed|return-path-domain-verified|return-path-domain-verification-failed", v)
				}
			}

			params := &helo.CreateWebhookRequest{}
			dataFields, err := decodeBody(requestData, params)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("url") {
				params.URL = url
			}
			if !cmd.Flags().Changed("url") && !dataFields["url"] {
				return usageErrorf("--url is required (or set \"url\" in --data)")
			}
			if cmd.Flags().Changed("events") {
				params.Events = nil
				for _, v := range events {
					params.Events = append(params.Events, helo.WebhookEvent(v))
				}
			}
			if !cmd.Flags().Changed("events") && !dataFields["events"] {
				return usageErrorf("--events is required (or set \"events\" in --data)")
			}
			if cmd.Flags().Changed("channel-id") {
				params.ChannelID = channelID
			}
			if cmd.Flags().Changed("additional-headers") {
				var v []helo.WebhookHeader
				if err := decodeJSONFlag("additional-headers", additionalHeaders, &v); err != nil {
					return err
				}
				params.AdditionalHeaders = v
			}
			if cmd.Flags().Changed("enabled") {
				params.Enabled = &enabled
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Webhooks.Create(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "The HTTPS URL to deliver webhook events to (required)")
	cmd.Flags().StringSliceVar(&events, "events", nil, "List of event types that will trigger this webhook (required; one of: message-accepted|message-processed|email-delivered|email-bounced|email-opened|link-clicked|recipient-complained|recipient-unsubscribed|recipient-resubscribed|domain-key-verified|domain-key-verification-failed|return-path-domain-verified|return-path-domain-verification-failed; comma-separated)")
	_ = cmd.RegisterFlagCompletionFunc("events", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"message-accepted", "message-processed", "email-delivered", "email-bounced", "email-opened", "link-clicked", "recipient-complained", "recipient-unsubscribed", "recipient-resubscribed", "domain-key-verified", "domain-key-verification-failed", "return-path-domain-verified", "return-path-domain-verification-failed"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&channelID, "channel-id", "", "Optional channel scope. If set, only events for this channel trigger the webhook")
	cmd.Flags().StringVar(&additionalHeaders, "additional-headers", "", "Custom headers to include in each webhook request. Each header name is prefixed with X-Customer- (JSON: [{\"name\": string, \"value\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "Whether the webhook is active or not. Defaults to true")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}

func newWebhooksRetrieveCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "retrieve <id>",
		Short: "Retrieve a webhook",
		Long:  "Fetches the details and configuration of a specific webhook.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
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
	var requestData string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a webhook",
		Long:  "Modifies an existing webhook by ID.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, v := range events {
				switch v {
				case "message-accepted", "message-processed", "email-delivered", "email-bounced", "email-opened", "link-clicked", "recipient-complained", "recipient-unsubscribed", "recipient-resubscribed", "domain-key-verified", "domain-key-verification-failed", "return-path-domain-verified", "return-path-domain-verification-failed":
				default:
					return usageErrorf("invalid value %q for --events: must be one of message-accepted|message-processed|email-delivered|email-bounced|email-opened|link-clicked|recipient-complained|recipient-unsubscribed|recipient-resubscribed|domain-key-verified|domain-key-verification-failed|return-path-domain-verified|return-path-domain-verification-failed", v)
				}
			}
			id := args[0]

			params := &helo.UpdateWebhookRequest{}
			if _, err := decodeBody(requestData, params); err != nil {
				return err
			}
			if cmd.Flags().Changed("url") {
				params.URL = url
			}
			if cmd.Flags().Changed("events") {
				params.Events = nil
				for _, v := range events {
					params.Events = append(params.Events, helo.WebhookEvent(v))
				}
			}
			if cmd.Flags().Changed("channel-id") {
				params.ChannelID = channelID
			}
			if cmd.Flags().Changed("additional-headers") {
				var v []helo.WebhookHeader
				if err := decodeJSONFlag("additional-headers", additionalHeaders, &v); err != nil {
					return err
				}
				params.AdditionalHeaders = v
			}
			if cmd.Flags().Changed("enabled") {
				params.Enabled = &enabled
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Webhooks.Update(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "New delivery URL. Must be HTTPS")
	cmd.Flags().StringSliceVar(&events, "events", nil, "Replaces the full list of event type subscriptions (one of: message-accepted|message-processed|email-delivered|email-bounced|email-opened|link-clicked|recipient-complained|recipient-unsubscribed|recipient-resubscribed|domain-key-verified|domain-key-verification-failed|return-path-domain-verified|return-path-domain-verification-failed; comma-separated)")
	_ = cmd.RegisterFlagCompletionFunc("events", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"message-accepted", "message-processed", "email-delivered", "email-bounced", "email-opened", "link-clicked", "recipient-complained", "recipient-unsubscribed", "recipient-resubscribed", "domain-key-verified", "domain-key-verification-failed", "return-path-domain-verified", "return-path-domain-verification-failed"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&channelID, "channel-id", "", "Set to a UUID to scope to a channel, or null to remove channel scoping")
	cmd.Flags().StringVar(&additionalHeaders, "additional-headers", "", "Replaces all custom headers. Set to null to remove all custom headers (JSON: [{\"name\": string, \"value\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "Enables or disables the webhook")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}

func newWebhooksDeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a webhook",
		Long:  "Permanently removes a webhook.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
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
		Long:  "Regenerate the signing key used for the webhook signature. This operation replaces the old key.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Webhooks.RegenerateSigningKey(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}
