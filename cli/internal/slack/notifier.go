package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kanaya/jobboard-cli/internal/config"
)

type Notifier struct {
	config     config.SlackConfig
	httpClient *http.Client
}

type Payload struct {
	Command    string
	Tag        string
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   time.Duration
	Status     string
	ExitCode   int
	Error      string
}

func NewNotifier(config config.SlackConfig, httpClient *http.Client) *Notifier {
	return &Notifier{
		config:     config,
		httpClient: httpClient,
	}
}

func (n *Notifier) Enabled() bool {
	return n.config.Enabled()
}

func (n *Notifier) Notify(ctx context.Context, payload Payload) error {
	if !n.Enabled() {
		return nil
	}

	icon := ":white_check_mark:"
	if payload.Status == "failed" {
		icon = ":x:"
	}

	statusLabel := strings.ToUpper(payload.Status)
	var text strings.Builder
	text.WriteString(fmt.Sprintf("%s *jobboard* `%s`\n", icon, payload.Command))

	if payload.Tag != "" {
		text.WriteString(fmt.Sprintf("*Tag:* %s\n", payload.Tag))
	}

	text.WriteString(fmt.Sprintf("*Status:* %s (exit code %d)\n", statusLabel, payload.ExitCode))
	text.WriteString(fmt.Sprintf("*Started:* %s\n", payload.StartedAt.Format(time.RFC3339)))
	text.WriteString(fmt.Sprintf("*Finished:* %s\n", payload.FinishedAt.Format(time.RFC3339)))
	text.WriteString(fmt.Sprintf("*DurationHours:* %f\n", payload.Duration.Hours()))

	if trimmed := strings.TrimSpace(payload.Error); trimmed != "" {
		text.WriteString("*Error:*\n```")
		text.WriteString(truncateHeadTail(trimmed, 900))
		text.WriteString("```\n")
	}

	body, err := json.Marshal(map[string]string{
		"text": text.String(),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.config.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		msg, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("slack notification failed: %s: %s", resp.Status, strings.TrimSpace(string(msg)))
	}

	return nil
}

func truncateHeadTail(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 6 {
		return s[:max]
	}

	head := s[:(max-3)/2]
	tail := s[len(s)-((max-3)/2):]
	return head + "..." + tail
}
