package commands

import (
	"github.com/spf13/cobra"

	"encoding/json"
	"fmt"
	"github.com/balvig/helo-sdk-go"
)

// NewWebhookEndpointsCmd builds the cobra command group for WebhookEndpoints.
func NewWebhookEndpointsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook-endpoints",
		Short: "Create and manage webhook endpoints for event notifications.",
	}
	cmd.AddCommand(newWebhookEndpointsListCmd())
	cmd.AddCommand(newWebhookEndpointsCreateCmd())
	cmd.AddCommand(newWebhookEndpointsRetrieveCmd())
	cmd.AddCommand(newWebhookEndpointsUpdateCmd())
	cmd.AddCommand(newWebhookEndpointsDeleteCmd())
	cmd.AddCommand(newWebhookEndpointsRegenerateSigningKeyCmd())
	return cmd
}

func newWebhookEndpointsListCmd() *cobra.Command {
	var limit int
	var offset int
	var channelIds []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all webhook endpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}

			params := &helo.WebhookEndpointsListParams{}
			params.Limit = limit
			params.Offset = offset
			params.ChannelIds = channelIds

			result, err := client.WebhookEndpoints.List(cmd.Context(), params)
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

func newWebhookEndpointsCreateCmd() *cobra.Command {
	var url string
	var events []string
	var channelID string
	var additionalHeaders string
	var enabled bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a webhook endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			for _, v := range events {
				switch v {
				case "accepted", "processed", "bounced", "delivered", "opened", "clicked", "complained", "unsubscribed", "resubscribed":
				default:
					return fmt.Errorf("invalid value %q for --events: must be one of accepted|processed|bounced|delivered|opened|clicked|complained|unsubscribed|resubscribed", v)
				}
			}

			params := &helo.CreateWebhookEndpointRequest{}
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

			result, err := client.WebhookEndpoints.Create(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "url (required)")
	_ = cmd.MarkFlagRequired("url")
	cmd.Flags().StringSliceVar(&events, "events", nil, "events (required, one of: accepted|processed|bounced|delivered|opened|clicked|complained|unsubscribed|resubscribed)")
	_ = cmd.MarkFlagRequired("events")
	_ = cmd.RegisterFlagCompletionFunc("events", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"accepted", "processed", "bounced", "delivered", "opened", "clicked", "complained", "unsubscribed", "resubscribed"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&additionalHeaders, "additional-headers", "", "additionalHeaders (JSON)")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "enabled")

	return cmd
}

func newWebhookEndpointsRetrieveCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "retrieve <id>",
		Short: "Retrieve a webhook endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			result, err := client.WebhookEndpoints.Retrieve(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}

func newWebhookEndpointsUpdateCmd() *cobra.Command {
	var url string
	var events []string
	var channelID string
	var additionalHeaders string
	var enabled bool

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a webhook endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			for _, v := range events {
				switch v {
				case "accepted", "processed", "bounced", "delivered", "opened", "clicked", "complained", "unsubscribed", "resubscribed":
				default:
					return fmt.Errorf("invalid value %q for --events: must be one of accepted|processed|bounced|delivered|opened|clicked|complained|unsubscribed|resubscribed", v)
				}
			}
			id := args[0]

			params := &helo.UpdateWebhookEndpointRequest{}
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

			result, err := client.WebhookEndpoints.Update(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "url")
	cmd.Flags().StringSliceVar(&events, "events", nil, "events (one of: accepted|processed|bounced|delivered|opened|clicked|complained|unsubscribed|resubscribed)")
	_ = cmd.RegisterFlagCompletionFunc("events", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"accepted", "processed", "bounced", "delivered", "opened", "clicked", "complained", "unsubscribed", "resubscribed"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&channelID, "channel-id", "", "channelId")
	cmd.Flags().StringVar(&additionalHeaders, "additional-headers", "", "additionalHeaders (JSON)")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "enabled")

	return cmd
}

func newWebhookEndpointsDeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a webhook endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			if err := client.WebhookEndpoints.Delete(cmd.Context(), id); err != nil {
				return err
			}
			return emit(map[string]string{"status": "ok"})
		},
	}

	return cmd
}

func newWebhookEndpointsRegenerateSigningKeyCmd() *cobra.Command {

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

			result, err := client.WebhookEndpoints.RegenerateSigningKey(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}
