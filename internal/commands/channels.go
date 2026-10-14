package commands

import (
	"github.com/spf13/cobra"

	"github.com/helo-email/helo-sdk-go"
)

// NewChannelsCmd builds the cobra command group for Channels.
func NewChannelsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channels",
		Short: "Create and manage communication channels for organizing messages.",
	}
	cmd.AddCommand(newChannelsListCmd())
	cmd.AddCommand(newChannelsCreateCmd())
	cmd.AddCommand(newChannelsRetrieveCmd())
	cmd.AddCommand(newChannelsUpdateCmd())
	cmd.AddCommand(newChannelsDeleteCmd())
	return cmd
}

func newChannelsListCmd() *cobra.Command {
	var limit int
	var offset int
	var name string
	var channelIds []string
	var deliveryType string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all channels",
		Long:  "Retrieves a list of all channels accessible to the current user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if deliveryType != "" {
				switch deliveryType {
				case "live", "sandbox":
				default:
					return usageErrorf("invalid value %q for --delivery-type: must be one of live|sandbox", deliveryType)
				}
			}

			params := &helo.ChannelsListParams{}
			params.Limit = limit
			params.Offset = offset
			params.Name = name
			params.ChannelIds = channelIds
			params.DeliveryType = deliveryType

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Channels.List(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 0, "limit")
	cmd.Flags().IntVar(&offset, "offset", 0, "offset")
	cmd.Flags().StringVar(&name, "name", "", "name")
	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "channelIds (comma-separated)")
	cmd.Flags().StringVar(&deliveryType, "delivery-type", "", "deliveryType (one of: live|sandbox)")
	_ = cmd.RegisterFlagCompletionFunc("delivery-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"live", "sandbox"}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func newChannelsCreateCmd() *cobra.Command {
	var name string
	var deliveryType string
	var tracking string
	var requestData string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a channel",
		Long:  "Creates a new communication channel for organizing and routing messages.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if deliveryType != "" {
				switch deliveryType {
				case "live", "sandbox":
				default:
					return usageErrorf("invalid value %q for --delivery-type: must be one of live|sandbox", deliveryType)
				}
			}

			params := &helo.CreateChannelRequest{}
			dataFields, err := decodeBody(requestData, params)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("name") {
				params.Name = name
			}
			if !cmd.Flags().Changed("name") && !dataFields.set["name"] {
				return usageErrorf("--name is required (or set \"name\" in --data)")
			}
			if cmd.Flags().Changed("delivery-type") {
				params.DeliveryType = helo.DeliveryType(deliveryType)
			}
			if !cmd.Flags().Changed("delivery-type") && !dataFields.set["deliveryType"] {
				return usageErrorf("--delivery-type is required (or set \"deliveryType\" in --data)")
			}
			if cmd.Flags().Changed("tracking") {
				var v helo.CreateChannelTracking
				if err := decodeJSONFlag("tracking", tracking, &v); err != nil {
					return err
				}
				params.Tracking = &v
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Channels.Create(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "The display name for the channel (required)")
	cmd.Flags().StringVar(&deliveryType, "delivery-type", "", "Delivery mode for a channel. 'live' sends real emails, whereas 'sandbox' accepts and processes messages without delivering them to recipients (required; one of: live|sandbox)")
	_ = cmd.RegisterFlagCompletionFunc("delivery-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"live", "sandbox"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&tracking, "tracking", "", "tracking (JSON: {\"links\": boolean, \"opens\": boolean}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}

func newChannelsRetrieveCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "retrieve <id>",
		Short: "Retrieve a channel",
		Long:  "Fetches the details and configuration of a specific channel.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Channels.Retrieve(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}

func newChannelsUpdateCmd() *cobra.Command {
	var name string
	var deliveryType string
	var tracking string
	var requestData string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a channel",
		Long:  "Modifies an existing channel by ID.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if deliveryType != "" {
				switch deliveryType {
				case "live", "sandbox":
				default:
					return usageErrorf("invalid value %q for --delivery-type: must be one of live|sandbox", deliveryType)
				}
			}
			id := args[0]

			params := &helo.UpdateChannelRequest{}
			if _, err := decodeBody(requestData, params); err != nil {
				return err
			}
			if cmd.Flags().Changed("name") {
				params.Name = name
			}
			if cmd.Flags().Changed("delivery-type") {
				params.DeliveryType = helo.DeliveryType(deliveryType)
			}
			if cmd.Flags().Changed("tracking") {
				var v helo.UpdateChannelTracking
				if err := decodeJSONFlag("tracking", tracking, &v); err != nil {
					return err
				}
				params.Tracking = &v
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Channels.Update(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "The new display name for the channel. Omit to keep the current value")
	cmd.Flags().StringVar(&deliveryType, "delivery-type", "", "Delivery mode for a channel. 'live' sends real emails, whereas 'sandbox' accepts and processes messages without delivering them to recipients (one of: live|sandbox)")
	_ = cmd.RegisterFlagCompletionFunc("delivery-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"live", "sandbox"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&tracking, "tracking", "", "tracking (JSON: {\"links\": boolean, \"opens\": boolean}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}

func newChannelsDeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a channel",
		Long:  "Permanently removes a channel and all associated data.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
			if err := client.Channels.Delete(cmd.Context(), id); err != nil {
				return err
			}
			return emit(map[string]string{"status": "ok"})
		},
	}

	return cmd
}
