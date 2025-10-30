package hub

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

type Client struct {
	config     config.HubConfig
	httpClient *http.Client
}

func NewClient(config config.HubConfig, httpClient *http.Client) *Client {
	return &Client{
		config:     config,
		httpClient: httpClient,
	}
}

func (c *Client) Enabled() bool {
	return c.config.Enabled()
}

type startRequest struct {
	NodeToken string    `json:"node_token"`
	Tag       string    `json:"tag,omitempty"`
	StartedAt time.Time `json:"started_at"`
}

func (c *Client) Start(ctx context.Context, startedAt time.Time) error {
	if !c.Enabled() {
		return nil
	}

	payload := startRequest{
		NodeToken: c.config.NodeToken,
		Tag:       c.config.Tag,
		StartedAt: startedAt,
	}

	return c.post(ctx, "/api/job-trigger/start", payload)
}

type finishRequest struct {
	NodeToken     string    `json:"node_token"`
	Status        string    `json:"status"`
	FinishedAt    time.Time `json:"finished_at"`
	DurationHours float64   `json:"duration_hours"`
}

func (c *Client) Finish(ctx context.Context, status string, finishedAt time.Time, duration time.Duration) error {
	if !c.Enabled() {
		return nil
	}

	payload := finishRequest{
		NodeToken:     c.config.NodeToken,
		Status:        status,
		FinishedAt:    finishedAt,
		DurationHours: duration.Hours(),
	}

	return c.post(ctx, "/api/job-trigger/finish", payload)
}

func (c *Client) post(ctx context.Context, path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpoint := strings.TrimRight(c.config.URL, "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		msg, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("hub request failed: %s: %s", resp.Status, strings.TrimSpace(string(msg)))
	}

	return nil
}
