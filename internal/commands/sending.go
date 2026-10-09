package commands

import (
	"github.com/spf13/cobra"

	"encoding/json"
	"fmt"
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

	cmd := &cobra.Command{
		Use:   "send-transactional",
		Short: "Send a transactional email",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}

			params := &helo.SendMessageRequest{}
			if from != "" {
				var v helo.MailAddress
				if err := json.Unmarshal([]byte(from), &v); err != nil {
					return fmt.Errorf("invalid JSON for --from: %w", err)
				}
				params.From = v
			}
			if from == "" {
				return fmt.Errorf("--from is required")
			}
			if to != "" {
				var v []helo.MailAddress
				if err := json.Unmarshal([]byte(to), &v); err != nil {
					return fmt.Errorf("invalid JSON for --to: %w", err)
				}
				params.To = v
			}
			if to == "" {
				return fmt.Errorf("--to is required")
			}
			if cc != "" {
				var v []helo.MailAddress
				if err := json.Unmarshal([]byte(cc), &v); err != nil {
					return fmt.Errorf("invalid JSON for --cc: %w", err)
				}
				params.Cc = v
			}
			if bcc != "" {
				var v []helo.MailAddress
				if err := json.Unmarshal([]byte(bcc), &v); err != nil {
					return fmt.Errorf("invalid JSON for --bcc: %w", err)
				}
				params.Bcc = v
			}
			if replyTo != "" {
				var v []helo.MailAddress
				if err := json.Unmarshal([]byte(replyTo), &v); err != nil {
					return fmt.Errorf("invalid JSON for --reply-to: %w", err)
				}
				params.ReplyTo = v
			}
			params.Subject = subject
			params.Html = html
			params.Text = text
			if template != "" {
				var v helo.SendMessageRequestTemplate
				if err := json.Unmarshal([]byte(template), &v); err != nil {
					return fmt.Errorf("invalid JSON for --template: %w", err)
				}
				params.Template = &v
			}
			if tracking != "" {
				var v helo.SendMessageRequestTracking
				if err := json.Unmarshal([]byte(tracking), &v); err != nil {
					return fmt.Errorf("invalid JSON for --tracking: %w", err)
				}
				params.Tracking = &v
			}
			if attachments != "" {
				var v []helo.Attachment
				if err := json.Unmarshal([]byte(attachments), &v); err != nil {
					return fmt.Errorf("invalid JSON for --attachments: %w", err)
				}
				params.Attachments = v
			}
			params.Tags = tags
			if headers != "" {
				var v map[string]any
				if err := json.Unmarshal([]byte(headers), &v); err != nil {
					return fmt.Errorf("invalid JSON for --headers: %w", err)
				}
				params.Headers = v
			}
			if metadata != "" {
				var v map[string]any
				if err := json.Unmarshal([]byte(metadata), &v); err != nil {
					return fmt.Errorf("invalid JSON for --metadata: %w", err)
				}
				params.Metadata = v
			}

			opts := &helo.SendingSendTransactionalOptions{
				ChannelID:      channelID,
				IdempotencyKey: idempotencyKey,
			}
			result, err := client.Sending.SendTransactional(cmd.Context(), params, opts)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&from, "from", "", "from (JSON, required)")
	_ = cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&to, "to", "", "to (JSON, required)")
	_ = cmd.MarkFlagRequired("to")
	cmd.Flags().StringVar(&cc, "cc", "", "cc (JSON)")
	cmd.Flags().StringVar(&bcc, "bcc", "", "bcc (JSON)")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "replyTo (JSON)")
	cmd.Flags().StringVar(&subject, "subject", "", "subject")
	cmd.Flags().StringVar(&html, "html", "", "html")
	cmd.Flags().StringVar(&text, "text", "", "text")
	cmd.Flags().StringVar(&template, "template", "", "template (JSON)")
	cmd.Flags().StringVar(&tracking, "tracking", "", "tracking (JSON)")
	cmd.Flags().StringVar(&attachments, "attachments", "", "attachments (JSON)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags")
	cmd.Flags().StringVar(&headers, "headers", "", "headers (JSON)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "metadata (JSON)")
	cmd.Flags().StringVar(&channelID, "channel-id", "", "Used to specify a channel ID for sending when using an account-level API credential.")
	cmd.Flags().StringVar(&idempotencyKey, "idempotency-key", "", "A unique identifier used to prevent duplicate messages being sent when retrying failed requests.")

	return cmd
}

func newSendingSendTransactionalBatchCmd() *cobra.Command {
	var requests string
	var channelID string
	var idempotencyKey string

	cmd := &cobra.Command{
		Use:   "send-transactional-batch",
		Short: "Send transactional emails in batch",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}

			params := &helo.SendMessageBatchRequest{}
			if requests != "" {
				var v []helo.SendMessageRequest
				if err := json.Unmarshal([]byte(requests), &v); err != nil {
					return fmt.Errorf("invalid JSON for --requests: %w", err)
				}
				params.Requests = v
			}
			if requests == "" {
				return fmt.Errorf("--requests is required")
			}

			opts := &helo.SendingSendTransactionalBatchOptions{
				ChannelID:      channelID,
				IdempotencyKey: idempotencyKey,
			}
			result, err := client.Sending.SendTransactionalBatch(cmd.Context(), params, opts)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&requests, "requests", "", "requests (JSON, required)")
	_ = cmd.MarkFlagRequired("requests")
	cmd.Flags().StringVar(&channelID, "channel-id", "", "Used to specify a channel ID for sending when using an account-level API credential.")
	cmd.Flags().StringVar(&idempotencyKey, "idempotency-key", "", "A unique identifier used to prevent duplicate messages being sent when retrying failed requests.")

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

	cmd := &cobra.Command{
		Use:   "send-broadcast",
		Short: "Send a broadcast",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}

			params := &helo.SendBroadcastRequest{}
			if from != "" {
				var v helo.MailAddress
				if err := json.Unmarshal([]byte(from), &v); err != nil {
					return fmt.Errorf("invalid JSON for --from: %w", err)
				}
				params.From = v
			}
			if from == "" {
				return fmt.Errorf("--from is required")
			}
			if replyTo != "" {
				var v []helo.MailAddress
				if err := json.Unmarshal([]byte(replyTo), &v); err != nil {
					return fmt.Errorf("invalid JSON for --reply-to: %w", err)
				}
				params.ReplyTo = v
			}
			if template != "" {
				var v helo.SendBroadcastRequestTemplate
				if err := json.Unmarshal([]byte(template), &v); err != nil {
					return fmt.Errorf("invalid JSON for --template: %w", err)
				}
				params.Template = v
			}
			if template == "" {
				return fmt.Errorf("--template is required")
			}
			if tracking != "" {
				var v helo.SendBroadcastRequestTracking
				if err := json.Unmarshal([]byte(tracking), &v); err != nil {
					return fmt.Errorf("invalid JSON for --tracking: %w", err)
				}
				params.Tracking = &v
			}
			if attachments != "" {
				var v []helo.Attachment
				if err := json.Unmarshal([]byte(attachments), &v); err != nil {
					return fmt.Errorf("invalid JSON for --attachments: %w", err)
				}
				params.Attachments = v
			}
			params.Tags = tags
			if headers != "" {
				var v map[string]any
				if err := json.Unmarshal([]byte(headers), &v); err != nil {
					return fmt.Errorf("invalid JSON for --headers: %w", err)
				}
				params.Headers = v
			}
			if metadata != "" {
				var v map[string]any
				if err := json.Unmarshal([]byte(metadata), &v); err != nil {
					return fmt.Errorf("invalid JSON for --metadata: %w", err)
				}
				params.Metadata = v
			}
			if messages != "" {
				var v []helo.SendBroadcastRequestMessage
				if err := json.Unmarshal([]byte(messages), &v); err != nil {
					return fmt.Errorf("invalid JSON for --messages: %w", err)
				}
				params.Messages = v
			}
			if messages == "" {
				return fmt.Errorf("--messages is required")
			}

			opts := &helo.SendingSendBroadcastOptions{
				ChannelID:      channelID,
				IdempotencyKey: idempotencyKey,
			}
			result, err := client.Sending.SendBroadcast(cmd.Context(), params, opts)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&from, "from", "", "from (JSON, required)")
	_ = cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "replyTo (JSON)")
	cmd.Flags().StringVar(&template, "template", "", "template (JSON, required)")
	_ = cmd.MarkFlagRequired("template")
	cmd.Flags().StringVar(&tracking, "tracking", "", "tracking (JSON)")
	cmd.Flags().StringVar(&attachments, "attachments", "", "attachments (JSON)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags")
	cmd.Flags().StringVar(&headers, "headers", "", "headers (JSON)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "metadata (JSON)")
	cmd.Flags().StringVar(&messages, "messages", "", "messages (JSON, required)")
	_ = cmd.MarkFlagRequired("messages")
	cmd.Flags().StringVar(&channelID, "channel-id", "", "Used to specify a channel ID for sending when using an account-level API credential.")
	cmd.Flags().StringVar(&idempotencyKey, "idempotency-key", "", "A unique identifier used to prevent duplicate messages being sent when retrying failed requests.")

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

	cmd := &cobra.Command{
		Use:   "send-broadcast-message",
		Short: "Send a single broadcast email",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}

			params := &helo.SendMessageRequest{}
			if from != "" {
				var v helo.MailAddress
				if err := json.Unmarshal([]byte(from), &v); err != nil {
					return fmt.Errorf("invalid JSON for --from: %w", err)
				}
				params.From = v
			}
			if from == "" {
				return fmt.Errorf("--from is required")
			}
			if to != "" {
				var v []helo.MailAddress
				if err := json.Unmarshal([]byte(to), &v); err != nil {
					return fmt.Errorf("invalid JSON for --to: %w", err)
				}
				params.To = v
			}
			if to == "" {
				return fmt.Errorf("--to is required")
			}
			if cc != "" {
				var v []helo.MailAddress
				if err := json.Unmarshal([]byte(cc), &v); err != nil {
					return fmt.Errorf("invalid JSON for --cc: %w", err)
				}
				params.Cc = v
			}
			if bcc != "" {
				var v []helo.MailAddress
				if err := json.Unmarshal([]byte(bcc), &v); err != nil {
					return fmt.Errorf("invalid JSON for --bcc: %w", err)
				}
				params.Bcc = v
			}
			if replyTo != "" {
				var v []helo.MailAddress
				if err := json.Unmarshal([]byte(replyTo), &v); err != nil {
					return fmt.Errorf("invalid JSON for --reply-to: %w", err)
				}
				params.ReplyTo = v
			}
			params.Subject = subject
			params.Html = html
			params.Text = text
			if template != "" {
				var v helo.SendMessageRequestTemplate
				if err := json.Unmarshal([]byte(template), &v); err != nil {
					return fmt.Errorf("invalid JSON for --template: %w", err)
				}
				params.Template = &v
			}
			if tracking != "" {
				var v helo.SendMessageRequestTracking
				if err := json.Unmarshal([]byte(tracking), &v); err != nil {
					return fmt.Errorf("invalid JSON for --tracking: %w", err)
				}
				params.Tracking = &v
			}
			if attachments != "" {
				var v []helo.Attachment
				if err := json.Unmarshal([]byte(attachments), &v); err != nil {
					return fmt.Errorf("invalid JSON for --attachments: %w", err)
				}
				params.Attachments = v
			}
			params.Tags = tags
			if headers != "" {
				var v map[string]any
				if err := json.Unmarshal([]byte(headers), &v); err != nil {
					return fmt.Errorf("invalid JSON for --headers: %w", err)
				}
				params.Headers = v
			}
			if metadata != "" {
				var v map[string]any
				if err := json.Unmarshal([]byte(metadata), &v); err != nil {
					return fmt.Errorf("invalid JSON for --metadata: %w", err)
				}
				params.Metadata = v
			}

			opts := &helo.SendingSendBroadcastMessageOptions{
				ChannelID:      channelID,
				IdempotencyKey: idempotencyKey,
			}
			result, err := client.Sending.SendBroadcastMessage(cmd.Context(), params, opts)
			if err != nil {
				return err
			}
			return emit(result)
		},
	}

	cmd.Flags().StringVar(&from, "from", "", "from (JSON, required)")
	_ = cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&to, "to", "", "to (JSON, required)")
	_ = cmd.MarkFlagRequired("to")
	cmd.Flags().StringVar(&cc, "cc", "", "cc (JSON)")
	cmd.Flags().StringVar(&bcc, "bcc", "", "bcc (JSON)")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "replyTo (JSON)")
	cmd.Flags().StringVar(&subject, "subject", "", "subject")
	cmd.Flags().StringVar(&html, "html", "", "html")
	cmd.Flags().StringVar(&text, "text", "", "text")
	cmd.Flags().StringVar(&template, "template", "", "template (JSON)")
	cmd.Flags().StringVar(&tracking, "tracking", "", "tracking (JSON)")
	cmd.Flags().StringVar(&attachments, "attachments", "", "attachments (JSON)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags")
	cmd.Flags().StringVar(&headers, "headers", "", "headers (JSON)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "metadata (JSON)")
	cmd.Flags().StringVar(&channelID, "channel-id", "", "Used to specify a channel ID for sending when using an account-level API credential.")
	cmd.Flags().StringVar(&idempotencyKey, "idempotency-key", "", "A unique identifier used to prevent duplicate messages being sent when retrying failed requests.")

	return cmd
}
