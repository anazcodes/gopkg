package goslack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
)

func TestNew(t *testing.T) {
	client := New("test-token", "test-channel")
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
	if client.channel != "test-channel" {
		t.Errorf("Expected channel 'test-channel', got '%s'", client.channel)
	}
	if client.severity != SeverityInfo {
		t.Errorf("Expected default severity 'info', got '%s'", client.severity)
	}
}

func TestSeverityChaining(t *testing.T) {
	client := New("test-token", "test-channel")

	infoClient := client.Info()
	if infoClient.severity != SeverityInfo {
		t.Errorf("Expected severity 'info', got '%s'", infoClient.severity)
	}

	warnClient := client.Warn()
	if warnClient.severity != SeverityWarning {
		t.Errorf("Expected severity 'warning', got '%s'", warnClient.severity)
	}

	errorClient := client.Error()
	if errorClient.severity != SeverityError {
		t.Errorf("Expected severity 'error', got '%s'", errorClient.severity)
	}

	successClient := client.Success()
	if successClient.severity != SeveritySuccess {
		t.Errorf("Expected severity 'success', got '%s'", successClient.severity)
	}

	// Verify original client is unchanged
	if client.severity != SeverityInfo {
		t.Errorf("Original client severity should remain 'info', got '%s'", client.severity)
	}
}

func TestSendMessage(t *testing.T) {
	// Create a mock Slack server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the request
		if err := r.ParseForm(); err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		// Verify the channel
		channel := r.FormValue("channel")
		if channel != "test-channel" {
			t.Errorf("Expected channel 'test-channel', got '%s'", channel)
		}

		// Verify attachments
		attachmentsJSON := r.FormValue("attachments")
		var attachments []slack.Attachment
		if err := json.Unmarshal([]byte(attachmentsJSON), &attachments); err != nil {
			t.Fatalf("Failed to unmarshal attachments: %v", err)
		}

		if len(attachments) != 1 {
			t.Fatalf("Expected 1 attachment, got %d", len(attachments))
		}

		if attachments[0].Text != "test message" {
			t.Errorf("Expected text 'test message', got '%s'", attachments[0].Text)
		}

		if attachments[0].Color != "#36a64f" { // Info color
			t.Errorf("Expected color '#36a64f', got '%s'", attachments[0].Color)
		}

		// Send mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":      true,
			"channel": "test-channel",
			"ts":      "1234567890.123456",
		})
	}))
	defer server.Close()

	// Create client with mock server
	client := New("test-token", "test-channel", slack.OptionAPIURL(server.URL+"/"))

	// Test SendMessage
	err := client.SendMessage(context.Background(), "test message")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestSendMessageError(t *testing.T) {
	// Create a mock Slack server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":    false,
			"error": "channel_not_found",
		})
	}))
	defer server.Close()

	client := New("test-token", "test-channel", slack.OptionAPIURL(server.URL+"/"))

	err := client.SendMessage(context.Background(), "test message")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestSendAttachment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		attachmentsJSON := r.FormValue("attachments")
		var attachments []slack.Attachment
		if err := json.Unmarshal([]byte(attachmentsJSON), &attachments); err != nil {
			t.Fatalf("Failed to unmarshal attachments: %v", err)
		}

		if len(attachments) != 2 {
			t.Fatalf("Expected 2 attachments, got %d", len(attachments))
		}

		// First attachment has explicit color
		if attachments[0].Color != "#ff0000" {
			t.Errorf("Expected color '#ff0000', got '%s'", attachments[0].Color)
		}

		// Second attachment should get default severity color (warning in this case)
		if attachments[1].Color != "#ffae42" {
			t.Errorf("Expected color '#ffae42', got '%s'", attachments[1].Color)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":      true,
			"channel": "test-channel",
			"ts":      "1234567890.123456",
		})
	}))
	defer server.Close()

	client := New("test-token", "test-channel", slack.OptionAPIURL(server.URL+"/"))

	err := client.Warn().SendAttachment(context.Background(),
		slack.Attachment{Text: "attachment 1", Color: "#ff0000"},
		slack.Attachment{Text: "attachment 2"}, // No color, should use default
	)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestColorForSeverity(t *testing.T) {
	client := New("test-token", "test-channel")

	tests := []struct {
		severity Severity
		expected string
	}{
		{SeverityInfo, "#36a64f"},
		{SeverityWarning, "#ffae42"},
		{SeverityError, "#ff0000"},
		{SeveritySuccess, "#2eb886"},
		{Severity("unknown"), "#dddddd"},
	}

	for _, tt := range tests {
		client.severity = tt.severity
		got := client.colorForSeverity()
		if got != tt.expected {
			t.Errorf("For severity '%s', expected color '%s', got '%s'", tt.severity, tt.expected, got)
		}
	}
}
