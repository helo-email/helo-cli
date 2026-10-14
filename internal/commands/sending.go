package commands

import (
	"github.com/spf13/cobra"

	"github.com/helo-email/helo-sdk-go"
)

// NewSendingCmd builds the cobra command group for Sending.
func NewSendingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sending",
		Short: "Send transactional and broadcast emails.",
	}
	cmd.AddCommand(newSendingSendTransactionalCmd())
	cmd.AddCommand(newSendingSendTransactionalBatchCmd())
	cmd.AddCommand(newSendingSendBroadcastCmd())
	cmd.AddCommand(newSendingSendBroadcastMessageCmd())
	return cmd
}

func newSendingSendTransactionalCmd() *cobra.Command {
	var from string
	var to string
	var cc string
	var bcc string
	var replyTo string
	var subject string
	var html string
	var text string
	var template string
	var tracking string
	var attachments string
	var tags []string
	var headers string
	var metadata string
	var channelID string
	var idempotencyKey string
	var requestData string

	cmd := &cobra.Command{
		Use:   "send-transactional",
		Short: "Send a transactional email",
		Long:  "Sends a single transactional email such as receipts, confirmations, or notifications.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.SendMessageRequest{}
			dataFields, err := decodeBody(requestData, params)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("from") {
				var v helo.MailAddress
				if err := decodeJSONFlag("from", from, &v); err != nil {
					return err
				}
				params.From = v
			}
			if !cmd.Flags().Changed("from") && !dataFields.set["from"] {
				return usageErrorf("--from is required (or set \"from\" in --data)")
			}
			if cmd.Flags().Changed("to") {
				var v []helo.MailAddress
				if err := decodeJSONFlag("to", to, &v); err != nil {
					return err
				}
				params.To = v
			}
			if !cmd.Flags().Changed("to") && !dataFields.set["to"] {
				return usageErrorf("--to is required (or set \"to\" in --data)")
			}
			if cmd.Flags().Changed("cc") {
				var v []helo.MailAddress
				if err := decodeJSONFlag("cc", cc, &v); err != nil {
					return err
				}
				params.Cc = v
			}
			if cmd.Flags().Changed("bcc") {
				var v []helo.MailAddress
				if err := decodeJSONFlag("bcc", bcc, &v); err != nil {
					return err
				}
				params.Bcc = v
			}
			if cmd.Flags().Changed("reply-to") {
				var v []helo.MailAddress
				if err := decodeJSONFlag("reply-to", replyTo, &v); err != nil {
					return err
				}
				params.ReplyTo = v
			}
			if cmd.Flags().Changed("subject") {
				params.Subject = subject
			}
			if cmd.Flags().Changed("html") {
				params.Html = html
			}
			if cmd.Flags().Changed("text") {
				params.Text = text
			}
			if cmd.Flags().Changed("template") {
				var v helo.SendMessageRequestTemplate
				if err := decodeJSONFlag("template", template, &v); err != nil {
					return err
				}
				params.Template = &v
			}
			if cmd.Flags().Changed("tracking") {
				var v helo.SendMessageRequestTracking
				if err := decodeJSONFlag("tracking", tracking, &v); err != nil {
					return err
				}
				params.Tracking = &v
			}
			if cmd.Flags().Changed("attachments") {
				var v []helo.Attachment
				if err := decodeJSONFlag("attachments", attachments, &v); err != nil {
					return err
				}
				params.Attachments = v
			}
			if cmd.Flags().Changed("tags") {
				params.Tags = tags
			}
			if cmd.Flags().Changed("headers") {
				var v map[string]any
				if err := decodeJSONFlag("headers", headers, &v); err != nil {
					return err
				}
				params.Headers = v
			}
			if cmd.Flags().Changed("metadata") {
				var v map[string]any
				if err := decodeJSONFlag("metadata", metadata, &v); err != nil {
					return err
				}
				params.Metadata = v
			}

			channelID = envDefault(channelID, channelIDEnv)
			opts := &helo.SendingSendTransactionalOptions{
				ChannelID:      channelID,
				IdempotencyKey: idempotencyKey,
			}
			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Sending.SendTransactional(cmd.Context(), params, opts)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&from, "from", "", "from (required; JSON: {\"email\": string, \"name\": string}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&to, "to", "", "to (required; JSON: [{\"email\": string, \"name\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&cc, "cc", "", "cc (JSON: [{\"email\": string, \"name\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&bcc, "bcc", "", "bcc (JSON: [{\"email\": string, \"name\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "replyTo (JSON: [{\"email\": string, \"name\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&subject, "subject", "", "subject")
	cmd.Flags().StringVar(&html, "html", "", "html")
	cmd.Flags().StringVar(&text, "text", "", "text")
	cmd.Flags().StringVar(&template, "template", "", "Either reference a stored template by 'id', or supply 'subject' and at least one of 'html' or 'text' inline (JSON: {\"id\": string, \"subject\": string, \"html\": string, \"text\": string, \"inlineStyles\": boolean, \"data\": object}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&tracking, "tracking", "", "tracking (JSON: {\"opens\": boolean, \"links\": boolean}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&attachments, "attachments", "", "attachments (JSON: [{\"content\": base64 string, \"contentId\": string, \"contentType\": string, \"fileName\": string, \"disposition\": \"attachment\"|\"inline\"}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags (comma-separated)")
	cmd.Flags().StringVar(&headers, "headers", "", "headers (JSON: {string: string}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "metadata (JSON: {string: string}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")
	cmd.Flags().StringVar(&channelID, "channel-id", "", "Used to specify a channel ID for sending when using an account-level API credential (default $HELO_CHANNEL_ID)")
	cmd.Flags().StringVar(&idempotencyKey, "idempotency-key", "", "A unique identifier used to prevent duplicate messages being sent when retrying failed requests (max length 36)")

	return cmd
}

func newSendingSendTransactionalBatchCmd() *cobra.Command {
	var requests string
	var channelID string
	var idempotencyKey string
	var requestData string

	cmd := &cobra.Command{
		Use:   "send-transactional-batch",
		Short: "Send transactional emails in batch",
		Long:  "Sends multiple transactional emails in a single API request for better performance.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.SendMessageBatchRequest{}
			dataFields, err := decodeBody(requestData, params)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("requests") {
				var v []helo.SendMessageRequest
				if err := decodeJSONFlag("requests", requests, &v); err != nil {
					return err
				}
				params.Requests = v
			}
			if !cmd.Flags().Changed("requests") && !dataFields.set["requests"] {
				return usageErrorf("--requests is required (or set \"requests\" in --data)")
			}

			channelID = envDefault(channelID, channelIDEnv)
			opts := &helo.SendingSendTransactionalBatchOptions{
				ChannelID:      channelID,
				IdempotencyKey: idempotencyKey,
			}
			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Sending.SendTransactionalBatch(cmd.Context(), params, opts)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&requests, "requests", "", "requests (required; JSON: [{\"from\": {…}, \"to\": [{…}], \"cc\": [{…}], \"bcc\": [{…}], \"replyTo\": [{…}], \"subject\": string, \"html\": string, \"text\": string, \"template\": {…}, \"tracking\": {…}, \"attachments\": [{…}], \"tags\": [string], \"headers\": {string: string}, \"metadata\": {string: string}}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")
	cmd.Flags().StringVar(&channelID, "channel-id", "", "Used to specify a channel ID for sending when using an account-level API credential (default $HELO_CHANNEL_ID)")
	cmd.Flags().StringVar(&idempotencyKey, "idempotency-key", "", "A unique identifier used to prevent duplicate messages being sent when retrying failed requests (max length 36)")

	return cmd
}

func newSendingSendBroadcastCmd() *cobra.Command {
	var from string
	var replyTo string
	var template string
	var tracking string
	var attachments string
	var tags []string
	var headers string
	var metadata string
	var messages string
	var channelID string
	var idempotencyKey string
	var requestData string

	cmd := &cobra.Command{
		Use:   "send-broadcast",
		Short: "Send a broadcast",
		Long:  "Sends a broadcast of multiple messages for marketing or announcement purposes.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.SendBroadcastRequest{}
			dataFields, err := decodeBody(requestData, params)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("from") {
				var v helo.MailAddress
				if err := decodeJSONFlag("from", from, &v); err != nil {
					return err
				}
				params.From = v
			}
			if !cmd.Flags().Changed("from") && !dataFields.set["from"] {
				return usageErrorf("--from is required (or set \"from\" in --data)")
			}
			if cmd.Flags().Changed("reply-to") {
				var v []helo.MailAddress
				if err := decodeJSONFlag("reply-to", replyTo, &v); err != nil {
					return err
				}
				params.ReplyTo = v
			}
			if cmd.Flags().Changed("template") {
				var v helo.SendBroadcastRequestTemplate
				if err := decodeJSONFlag("template", template, &v); err != nil {
					return err
				}
				params.Template = v
			}
			if !cmd.Flags().Changed("template") && !dataFields.set["template"] {
				return usageErrorf("--template is required (or set \"template\" in --data)")
			}
			if cmd.Flags().Changed("tracking") {
				var v helo.SendBroadcastRequestTracking
				if err := decodeJSONFlag("tracking", tracking, &v); err != nil {
					return err
				}
				params.Tracking = &v
			}
			if cmd.Flags().Changed("attachments") {
				var v []helo.Attachment
				if err := decodeJSONFlag("attachments", attachments, &v); err != nil {
					return err
				}
				params.Attachments = v
			}
			if cmd.Flags().Changed("tags") {
				params.Tags = tags
			}
			if cmd.Flags().Changed("headers") {
				var v map[string]any
				if err := decodeJSONFlag("headers", headers, &v); err != nil {
					return err
				}
				params.Headers = v
			}
			if cmd.Flags().Changed("metadata") {
				var v map[string]any
				if err := decodeJSONFlag("metadata", metadata, &v); err != nil {
					return err
				}
				params.Metadata = v
			}
			if cmd.Flags().Changed("messages") {
				var v []helo.SendBroadcastRequestMessage
				if err := decodeJSONFlag("messages", messages, &v); err != nil {
					return err
				}
				params.Messages = v
			}
			if !cmd.Flags().Changed("messages") && !dataFields.set["messages"] {
				return usageErrorf("--messages is required (or set \"messages\" in --data)")
			}

			channelID = envDefault(channelID, channelIDEnv)
			opts := &helo.SendingSendBroadcastOptions{
				ChannelID:      channelID,
				IdempotencyKey: idempotencyKey,
			}
			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Sending.SendBroadcast(cmd.Context(), params, opts)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&from, "from", "", "from (required; JSON: {\"email\": string, \"name\": string}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "Reply-to addresses (JSON: [{\"email\": string, \"name\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&template, "template", "", "Email template applied to every message in the broadcast. Either reference a stored template by 'id', or supply 'subject' and at least one of 'html' or 'text' inline (required; JSON: {\"id\": string, \"subject\": string, \"html\": string, \"text\": string, \"inlineStyles\": boolean, \"data\": object}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&tracking, "tracking", "", "Override channel-level open and link tracking settings (JSON: {\"opens\": boolean, \"links\": boolean}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&attachments, "attachments", "", "File attachments included with every message in the broadcast. Executable and potentially harmful file types are blocked (JSON: [{\"content\": base64 string, \"contentId\": string, \"contentType\": string, \"fileName\": string, \"disposition\": \"attachment\"|\"inline\"}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Up to 5 tags for filtering and analytics. Allowed characters are letters, numbers, hyphens, and underscores. Max 100 characters per tag (comma-separated)")
	cmd.Flags().StringVar(&headers, "headers", "", "Custom email headers applied to every message. Total size of all headers combined must not exceed 5,000 characters (JSON: {string: string}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "Custom key/value metadata stored with the broadcast. Max 10 fields; keys max 50 characters, values max 100 characters (JSON: {string: string}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&messages, "messages", "", "One entry per outbound email. At least one message is required (required; JSON: [{\"to\": [{\"email\": string, \"name\": string}], \"cc\": [{\"email\": string, \"name\": string}], \"bcc\": [{\"email\": string, \"name\": string}], \"tags\": [string], \"headers\": {string: string}, \"metadata\": {string: string}, \"data\": object}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")
	cmd.Flags().StringVar(&channelID, "channel-id", "", "Used to specify a channel ID for sending when using an account-level API credential (default $HELO_CHANNEL_ID)")
	cmd.Flags().StringVar(&idempotencyKey, "idempotency-key", "", "A unique identifier used to prevent duplicate messages being sent when retrying failed requests (max length 36)")

	return cmd
}

func newSendingSendBroadcastMessageCmd() *cobra.Command {
	var from string
	var to string
	var cc string
	var bcc string
	var replyTo string
	var subject string
	var html string
	var text string
	var template string
	var tracking string
	var attachments string
	var tags []string
	var headers string
	var metadata string
	var channelID string
	var idempotencyKey string
	var requestData string

	cmd := &cobra.Command{
		Use:   "send-broadcast-message",
		Short: "Send a single broadcast email",
		Long:  "Sends a single broadcast email message.\n\nSet the request body's fields with flags, or pass the whole body as JSON with --data (inline, @file, or @- for stdin); flags override fields from --data.",
		RunE: func(cmd *cobra.Command, args []string) error {

			params := &helo.SendMessageRequest{}
			dataFields, err := decodeBody(requestData, params)
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("from") {
				var v helo.MailAddress
				if err := decodeJSONFlag("from", from, &v); err != nil {
					return err
				}
				params.From = v
			}
			if !cmd.Flags().Changed("from") && !dataFields.set["from"] {
				return usageErrorf("--from is required (or set \"from\" in --data)")
			}
			if cmd.Flags().Changed("to") {
				var v []helo.MailAddress
				if err := decodeJSONFlag("to", to, &v); err != nil {
					return err
				}
				params.To = v
			}
			if !cmd.Flags().Changed("to") && !dataFields.set["to"] {
				return usageErrorf("--to is required (or set \"to\" in --data)")
			}
			if cmd.Flags().Changed("cc") {
				var v []helo.MailAddress
				if err := decodeJSONFlag("cc", cc, &v); err != nil {
					return err
				}
				params.Cc = v
			}
			if cmd.Flags().Changed("bcc") {
				var v []helo.MailAddress
				if err := decodeJSONFlag("bcc", bcc, &v); err != nil {
					return err
				}
				params.Bcc = v
			}
			if cmd.Flags().Changed("reply-to") {
				var v []helo.MailAddress
				if err := decodeJSONFlag("reply-to", replyTo, &v); err != nil {
					return err
				}
				params.ReplyTo = v
			}
			if cmd.Flags().Changed("subject") {
				params.Subject = subject
			}
			if cmd.Flags().Changed("html") {
				params.Html = html
			}
			if cmd.Flags().Changed("text") {
				params.Text = text
			}
			if cmd.Flags().Changed("template") {
				var v helo.SendMessageRequestTemplate
				if err := decodeJSONFlag("template", template, &v); err != nil {
					return err
				}
				params.Template = &v
			}
			if cmd.Flags().Changed("tracking") {
				var v helo.SendMessageRequestTracking
				if err := decodeJSONFlag("tracking", tracking, &v); err != nil {
					return err
				}
				params.Tracking = &v
			}
			if cmd.Flags().Changed("attachments") {
				var v []helo.Attachment
				if err := decodeJSONFlag("attachments", attachments, &v); err != nil {
					return err
				}
				params.Attachments = v
			}
			if cmd.Flags().Changed("tags") {
				params.Tags = tags
			}
			if cmd.Flags().Changed("headers") {
				var v map[string]any
				if err := decodeJSONFlag("headers", headers, &v); err != nil {
					return err
				}
				params.Headers = v
			}
			if cmd.Flags().Changed("metadata") {
				var v map[string]any
				if err := decodeJSONFlag("metadata", metadata, &v); err != nil {
					return err
				}
				params.Metadata = v
			}

			channelID = envDefault(channelID, channelIDEnv)
			opts := &helo.SendingSendBroadcastMessageOptions{
				ChannelID:      channelID,
				IdempotencyKey: idempotencyKey,
			}
			client, err := newClient()
			if err != nil {
				return err
			}
			result, err := client.Sending.SendBroadcastMessage(cmd.Context(), params, opts)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&from, "from", "", "from (required; JSON: {\"email\": string, \"name\": string}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&to, "to", "", "to (required; JSON: [{\"email\": string, \"name\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&cc, "cc", "", "cc (JSON: [{\"email\": string, \"name\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&bcc, "bcc", "", "bcc (JSON: [{\"email\": string, \"name\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "replyTo (JSON: [{\"email\": string, \"name\": string}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&subject, "subject", "", "subject")
	cmd.Flags().StringVar(&html, "html", "", "html")
	cmd.Flags().StringVar(&text, "text", "", "text")
	cmd.Flags().StringVar(&template, "template", "", "Either reference a stored template by 'id', or supply 'subject' and at least one of 'html' or 'text' inline (JSON: {\"id\": string, \"subject\": string, \"html\": string, \"text\": string, \"inlineStyles\": boolean, \"data\": object}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&tracking, "tracking", "", "tracking (JSON: {\"opens\": boolean, \"links\": boolean}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&attachments, "attachments", "", "attachments (JSON: [{\"content\": base64 string, \"contentId\": string, \"contentType\": string, \"fileName\": string, \"disposition\": \"attachment\"|\"inline\"}]; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags (comma-separated)")
	cmd.Flags().StringVar(&headers, "headers", "", "headers (JSON: {string: string}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "metadata (JSON: {string: string}; @file or @- reads it from a file or stdin)")
	cmd.Flags().StringVar(&requestData, "data", "", "The whole request body as JSON (inline, @file, or @- for stdin); other flags override its fields")
	cmd.Flags().StringVar(&channelID, "channel-id", "", "Used to specify a channel ID for sending when using an account-level API credential (default $HELO_CHANNEL_ID)")
	cmd.Flags().StringVar(&idempotencyKey, "idempotency-key", "", "A unique identifier used to prevent duplicate messages being sent when retrying failed requests (max length 36)")

	return cmd
}
