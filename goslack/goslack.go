package goslack

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/slack-go/slack"
)

type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
	SeveritySuccess Severity = "success"
)

var GlobalClient *Client

func (e *Client) SetGlobalClient() {
	GlobalClient = e
}

type Client struct {
	client   *slack.Client
	channel  string
	severity Severity
}

func New(token, channel string, options ...slack.Option) *Client {
	return &Client{
		client:   slack.New(token, options...),
		channel:  channel,
		severity: SeverityInfo,
	}
}

// clone creates a shallow copy of the emitter (safe for chaining).
func (e *Client) clone() *Client {
	c := *e
	return &c
}

// Info returns a cloned emitter with "info" severity.
func (e *Client) Info() *Client {
	c := e.clone()
	c.severity = SeverityInfo
	return c
}

// Warn returns a cloned emitter with "warning" severity.
func (e *Client) Warn() *Client {
	c := e.clone()
	c.severity = SeverityWarning
	return c
}

// Error returns a cloned emitter with "error" severity.
func (e *Client) Error() *Client {
	c := e.clone()
	c.severity = SeverityError
	return c
}

// Success returns a cloned emitter with "success" severity.
func (e *Client) Success() *Client {
	c := e.clone()
	c.severity = SeveritySuccess
	return c
}

// SendMessage sends a message to the configured channel.
func (e *Client) SendMessage(ctx context.Context, message string) error {
	color := e.colorForSeverity()
	attachment := slack.Attachment{
		Text:  message,
		Color: color,
	}

	_, _, err := e.client.PostMessageContext(ctx, e.channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("goslack: %w", err)
	}
	return nil
}

// SendAttachment sends attachments with severity color applied if missing.
func (e *Client) SendAttachment(ctx context.Context, attachments ...slack.Attachment) error {
	color := e.colorForSeverity()
	for i := range attachments {
		if attachments[i].Color == "" {
			attachments[i].Color = color
		}
	}
	_, _, err := e.client.PostMessageContext(ctx, e.channel, slack.MsgOptionAttachments(attachments...))
	if err != nil {
		return fmt.Errorf("goslack: %w", err)
	}
	return nil
}

// UploadFile uploads a file to Slack using a file reader
func (e *Client) UploadFile(ctx context.Context, filePath, title, comment string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("goslack: failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info to check size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("goslack: failed to get file info: %w", err)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("goslack: file is empty: %s", filePath)
	}

	_, err = e.client.UploadFileV2Context(ctx, slack.UploadFileV2Parameters{
		Channel:        e.channel,
		Filename:       filepath.Base(filePath),
		File:           filePath,
		FileSize:       int(fileInfo.Size()),
		Title:          title,
		InitialComment: comment,
	})
	if err != nil {
		return fmt.Errorf("goslack: slack file upload error: %w", err)
	}

	return nil
}

// colorForSeverity maps severity to Slack color hex.
func (e *Client) colorForSeverity() string {
	switch e.severity {
	case SeverityInfo:
		return "#36a64f"
	case SeverityWarning:
		return "#ffae42"
	case SeverityError:
		return "#ff0000"
	case SeveritySuccess:
		return "#2eb886"
	default:
		return "#dddddd"
	}
}
