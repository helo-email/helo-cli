package commands

import (
	"github.com/spf13/cobra"

	"fmt"
	"github.com/helo-email/helo-sdk-go"
)

// NewDomainsCmd builds the cobra command group for Domains.
func NewDomainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domains",
		Short: "Register, verify, and manage domains for email sending.",
	}
	cmd.AddCommand(newDomainsListCmd())
	cmd.AddCommand(newDomainsCreateCmd())
	cmd.AddCommand(newDomainsRetrieveCmd())
	cmd.AddCommand(newDomainsUpdateCmd())
	cmd.AddCommand(newDomainsDeleteCmd())
	cmd.AddCommand(newDomainsVerifyCmd())
	cmd.AddCommand(newDomainsRotateKeyCmd())
	return cmd
}

func newDomainsListCmd() *cobra.Command {
	var limit int
	var offset int
	var name string
	var channelIds []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}

			params := &helo.DomainsListParams{}
			params.Limit = limit
			params.Offset = offset
			params.Name = name
			params.ChannelIds = channelIds

			result, err := client.Domains.List(cmd.Context(), params)
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

	return cmd
}

func newDomainsCreateCmd() *cobra.Command {
	var name string
	var channelIds []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}

			params := &helo.CreateDomainRequest{}
			params.Name = name
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			params.ChannelIds = channelIds

			result, err := client.Domains.Create(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "name (required)")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "channelIds")

	return cmd
}

func newDomainsRetrieveCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "retrieve <id>",
		Short: "Retrieve a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			result, err := client.Domains.Retrieve(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}

func newDomainsUpdateCmd() *cobra.Command {
	var channelIds []string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			params := &helo.UpdateDomainRequest{}
			params.ChannelIds = channelIds

			result, err := client.Domains.Update(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "channelIds")

	return cmd
}

func newDomainsDeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			if err := client.Domains.Delete(cmd.Context(), id); err != nil {
				return err
			}
			return emit(map[string]string{"status": "ok"})
		},
	}

	return cmd
}

func newDomainsVerifyCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "verify <id>",
		Short: "Verify a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			result, err := client.Domains.Verify(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}

func newDomainsRotateKeyCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "rotate-key <id>",
		Short: "Rotate a domain key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			id := args[0]

			result, err := client.Domains.RotateKey(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}
