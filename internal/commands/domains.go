package commands

import (
	"github.com/spf13/cobra"

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
		Long:  "Retrieves all domains associated with the current account, including their verification status.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.DomainsListParams{}
			params.Limit = limit
			params.Offset = offset
			params.Name = name
			params.ChannelIds = channelIds

			client, err := newClient()
			if err != nil {
				return err
			}
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
	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "channelIds (comma-separated)")

	return cmd
}

func newDomainsCreateCmd() *cobra.Command {
	var name string
	var channelIds []string
	var requestData string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a domain",
		Long:  "Registers a new domain for sending emails. The domain must be verified before it can be used.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.CreateDomainRequest{}
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
			if cmd.Flags().Changed("channel-ids") {
				params.ChannelIds = channelIds
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Domains.Create(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "name (required)")
	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "channelIds (comma-separated)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}

func newDomainsRetrieveCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "retrieve <id>",
		Short: "Retrieve a domain",
		Long:  "Gets detailed information about a specific domain, including verification status and configuration.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
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
	var requestData string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a domain",
		Long:  "Modifies an existing domain by ID.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			params := &helo.UpdateDomainRequest{}
			if _, err := decodeBody(requestData, params); err != nil {
				return err
			}
			if cmd.Flags().Changed("channel-ids") {
				params.ChannelIds = channelIds
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Domains.Update(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "channelIds (comma-separated)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}

func newDomainsDeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a domain",
		Long:  "Removes a domain from the account. This will stop all email sending from this domain.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
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
		Long:  "Initiates the domain verification process by checking DNS records.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
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
		Long:  "Generates new DKIM keys for the domain. This is recommended for security best practices.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Domains.RotateKey(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}
