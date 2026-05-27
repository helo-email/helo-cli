package commands

import (
	"github.com/spf13/cobra"

	"encoding/json"
	"fmt"
	"github.com/balvig/helo-sdk-go"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			if deliveryType != "" {
				switch deliveryType {
				case "live", "sandbox":
				default:
					return fmt.Errorf("invalid value %q for --delivery-type: must be one of live|sandbox", deliveryType)
				}
			}

			params := &helo.ChannelsListParams{}
			params.Limit = limit
			params.Offset = offset
			params.Name = name
			params.ChannelIds = channelIds
			params.DeliveryType = deliveryType

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
	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "channelIds")
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

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			if deliveryType != "" {
				switch deliveryType {
				case "live", "sandbox":
				default:
					return fmt.Errorf("invalid value %q for --delivery-type: must be one of live|sandbox", deliveryType)
				}
			}

			params := &helo.CreateChannelRequest{}
			params.Name = name
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			params.DeliveryType = helo.DeliveryType(deliveryType)
			if deliveryType == "" {
				return fmt.Errorf("--delivery-type is required")
			}
			if tracking != "" {
				var v helo.CreateChannelTracking
				if err := json.Unmarshal([]byte(tracking), &v); err != nil {
					return fmt.Errorf("invalid JSON for --tracking: %w", err)
				}
				params.Tracking = v
			}

			result, err := client.Channels.Create(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "name (required)")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&deliveryType, "delivery-type", "", "deliveryType (required, one of: live|sandbox)")
	_ = cmd.MarkFlagRequired("delivery-type")
	_ = cmd.RegisterFlagCompletionFunc("delivery-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"live", "sandbox"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&tracking, "tracking", "", "tracking (JSON)")

	return cmd
}

func newChannelsRetrieveCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "retrieve <id>",
		Short: "Retrieve a channel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

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

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a channel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			if deliveryType != "" {
				switch deliveryType {
				case "live", "sandbox":
				default:
					return fmt.Errorf("invalid value %q for --delivery-type: must be one of live|sandbox", deliveryType)
				}
			}
			id := args[0]

			params := &helo.UpdateChannelRequest{}
			params.Name = name
			params.DeliveryType = helo.DeliveryType(deliveryType)
			if tracking != "" {
				var v helo.UpdateChannelTracking
				if err := json.Unmarshal([]byte(tracking), &v); err != nil {
					return fmt.Errorf("invalid JSON for --tracking: %w", err)
				}
				params.Tracking = v
			}

			result, err := client.Channels.Update(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "name")
	cmd.Flags().StringVar(&deliveryType, "delivery-type", "", "deliveryType (one of: live|sandbox)")
	_ = cmd.RegisterFlagCompletionFunc("delivery-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"live", "sandbox"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVar(&tracking, "tracking", "", "tracking (JSON)")

	return cmd
}

func newChannelsDeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a channel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			if err := client.Channels.Delete(cmd.Context(), id); err != nil {
				return err
			}
			return emit(map[string]string{"status": "ok"})
		},
	}

	return cmd
}
