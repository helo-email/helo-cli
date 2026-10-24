package commands

import (
	"github.com/spf13/cobra"

	"github.com/helo-email/helo-sdk-go"
)

// NewTemplatesCmd builds the cobra command group for Templates.
func NewTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Manage templates for sending.",
	}
	cmd.AddCommand(newTemplatesListCmd())
	cmd.AddCommand(newTemplatesCreateCmd())
	cmd.AddCommand(newTemplatesRenderDraftCmd())
	cmd.AddCommand(newTemplatesRetrieveCmd())
	cmd.AddCommand(newTemplatesUpdateCmd())
	cmd.AddCommand(newTemplatesDeleteCmd())
	cmd.AddCommand(newTemplatesRenderCmd())
	return cmd
}

func newTemplatesListCmd() *cobra.Command {
	var limit int
	var offset int
	var name string
	var channelIds []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List operation",
		Long:  "Fetch a paginated list of templates.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.TemplatesListParams{}
			params.Limit = limit
			params.Offset = offset
			params.Name = name
			params.ChannelIds = channelIds

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Templates.List(cmd.Context(), params)
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

func newTemplatesCreateCmd() *cobra.Command {
	var name string
	var subject string
	var html string
	var mjml string
	var text string
	var channelIds []string
	var requestData string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create operation",
		Long:  "Create a new template.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.CreateTemplateRequest{}
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
			if cmd.Flags().Changed("subject") {
				params.Subject = subject
			}
			if !cmd.Flags().Changed("subject") && !dataFields.set["subject"] {
				return usageErrorf("--subject is required (or set \"subject\" in --data)")
			}
			if cmd.Flags().Changed("html") {
				params.Html = html
			}
			if cmd.Flags().Changed("mjml") {
				params.Mjml = mjml
			}
			if cmd.Flags().Changed("text") {
				params.Text = text
			}
			if cmd.Flags().Changed("channel-ids") {
				params.ChannelIds = channelIds
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Templates.Create(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "name (required)")
	cmd.Flags().StringVar(&subject, "subject", "", "subject (required)")
	cmd.Flags().StringVar(&html, "html", "", "html")
	cmd.Flags().StringVar(&mjml, "mjml", "", "mjml")
	cmd.Flags().StringVar(&text, "text", "", "text")
	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "channelIds (comma-separated)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}

func newTemplatesRenderDraftCmd() *cobra.Command {
	var subject string
	var html string
	var mjml string
	var text string
	var inlineStyles bool
	var data string
	var requestData string

	cmd := &cobra.Command{
		Use:   "render-draft",
		Short: "RenderDraft operation",
		Long:  "Render ad-hoc template content (Subject/Html/Text) with sample data, without persisting anything.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.RenderDraftRequest{}
			dataFields, err := decodeBody(requestData, params)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("subject") {
				params.Subject = subject
			}
			if !cmd.Flags().Changed("subject") && !dataFields.set["subject"] {
				return usageErrorf("--subject is required (or set \"subject\" in --data)")
			}
			if cmd.Flags().Changed("html") {
				params.Html = html
			}
			if cmd.Flags().Changed("mjml") {
				params.Mjml = mjml
			}
			if cmd.Flags().Changed("text") {
				params.Text = text
			}
			if cmd.Flags().Changed("inline-styles") {
				params.InlineStyles = &inlineStyles
			}
			if cmd.Flags().Changed("data") {
				var v map[string]any
				if err := decodeJSONFlag("data", data, &v); err != nil {
					return err
				}
				params.Data = v
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Templates.RenderDraft(cmd.Context(), params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&subject, "subject", "", "subject (required)")
	cmd.Flags().StringVar(&html, "html", "", "html")
	cmd.Flags().StringVar(&mjml, "mjml", "", "mjml")
	cmd.Flags().StringVar(&text, "text", "", "text")
	cmd.Flags().BoolVar(&inlineStyles, "inline-styles", false, "Inline CSS into the HTML, as a send does. Defaults to true, so the preview matches what recipients receive")
	cmd.Flags().StringVar(&data, "data", "", "Sample data used to interpolate Liquid variables (JSON: object; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}

func newTemplatesRetrieveCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "retrieve <id>",
		Short: "Retrieve operation",
		Long:  "Fetch a single template.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Templates.Retrieve(cmd.Context(), id)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	return cmd
}

func newTemplatesUpdateCmd() *cobra.Command {
	var name string
	var subject string
	var html string
	var mjml string
	var text string
	var channelIds []string
	var version int64
	var clearHtml bool
	var clearMjml bool
	var clearText bool
	var clearChannelIds bool
	var requestData string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update operation",
		Long:  "Update an existing template.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			params := &helo.UpdateTemplateRequest{}
			dataFields, err := decodeBody(requestData, params)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("name") {
				params.Name = name
			}
			if cmd.Flags().Changed("subject") {
				params.Subject = subject
			}
			if cmd.Flags().Changed("html") {
				params.Html = html
			}
			if cmd.Flags().Changed("mjml") {
				params.Mjml = mjml
			}
			if cmd.Flags().Changed("text") {
				params.Text = text
			}
			if cmd.Flags().Changed("channel-ids") {
				params.ChannelIds = channelIds
			}
			if cmd.Flags().Changed("version") {
				params.Version = version
			}
			// A null in --data, or a --clear-* flag, clears a field; setting the field's own flag wins over --data.
			if clearHtml && cmd.Flags().Changed("html") {
				return usageErrorf("--clear-html and --html can't be used together")
			}
			if clearHtml || (dataFields.null["html"] && !cmd.Flags().Changed("html")) {
				params.NullFields = append(params.NullFields, "html")
			}
			if clearMjml && cmd.Flags().Changed("mjml") {
				return usageErrorf("--clear-mjml and --mjml can't be used together")
			}
			if clearMjml || (dataFields.null["mjml"] && !cmd.Flags().Changed("mjml")) {
				params.NullFields = append(params.NullFields, "mjml")
			}
			if clearText && cmd.Flags().Changed("text") {
				return usageErrorf("--clear-text and --text can't be used together")
			}
			if clearText || (dataFields.null["text"] && !cmd.Flags().Changed("text")) {
				params.NullFields = append(params.NullFields, "text")
			}
			if clearChannelIds && cmd.Flags().Changed("channel-ids") {
				return usageErrorf("--clear-channel-ids and --channel-ids can't be used together")
			}
			if clearChannelIds || (dataFields.null["channelIds"] && !cmd.Flags().Changed("channel-ids")) {
				params.NullFields = append(params.NullFields, "channelIds")
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Templates.Update(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Renames the template. Cannot be cleared")
	cmd.Flags().StringVar(&subject, "subject", "", "Cannot be cleared")
	cmd.Flags().StringVar(&html, "html", "", "Replaces the html, and clears any stored mjml (switching the template back to plain HTML authoring). Set to null to remove the html - blank values are rejected (--clear-html sends null)")
	cmd.Flags().StringVar(&mjml, "mjml", "", "Replaces the mjml source, recompiling the html from it. Set to null to remove the mjml source while keeping the html already compiled from it - blank values are rejected (--clear-mjml sends null)")
	cmd.Flags().StringVar(&text, "text", "", "Replaces the plain text content. Set to null to remove it - blank values are rejected (--clear-text sends null)")
	cmd.Flags().StringSliceVar(&channelIds, "channel-ids", nil, "Replaces the channels the template is scoped to. Set to null or an empty array to remove channel scoping, making the template available to every channel in the account (comma-separated; --clear-channel-ids sends null)")
	cmd.Flags().Int64Var(&version, "version", 0, "The version of the template the update was based on, as returned by a previous read. If the template has since been modified the update is rejected with 409 Conflict, so concurrent edits aren't silently overwritten")
	cmd.Flags().BoolVar(&clearHtml, "clear-html", false, "Clear html by sending it as null")
	cmd.Flags().BoolVar(&clearMjml, "clear-mjml", false, "Clear mjml by sending it as null")
	cmd.Flags().BoolVar(&clearText, "clear-text", false, "Clear text by sending it as null")
	cmd.Flags().BoolVar(&clearChannelIds, "clear-channel-ids", false, "Clear channelIds by sending it as null")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}

func newTemplatesDeleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete operation",
		Long:  "Delete a template.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			client, err := newClient()
			if err != nil {
				return err
			}
			if err := client.Templates.Delete(cmd.Context(), id); err != nil {
				return err
			}
			return emit(map[string]string{"status": "ok"})
		},
	}

	return cmd
}

func newTemplatesRenderCmd() *cobra.Command {
	var inlineStyles bool
	var data string
	var requestData string

	cmd := &cobra.Command{
		Use:   "render <id>",
		Short: "Render operation",
		Long:  "Render a saved template with sample data.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			params := &helo.RenderTemplateRequest{}
			if _, err := decodeBody(requestData, params); err != nil {
				return err
			}
			if cmd.Flags().Changed("inline-styles") {
				params.InlineStyles = &inlineStyles
			}
			if cmd.Flags().Changed("data") {
				var v map[string]any
				if err := decodeJSONFlag("data", data, &v); err != nil {
					return err
				}
				params.Data = v
			}

			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Templates.Render(cmd.Context(), id, params)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().BoolVar(&inlineStyles, "inline-styles", false, "Inline CSS into the HTML, as a send does. Defaults to true, so the preview matches what recipients receive")
	cmd.Flags().StringVar(&data, "data", "", "Sample data used to interpolate Liquid variables (JSON: object; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")

	return cmd
}
