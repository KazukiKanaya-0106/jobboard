package config

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Hub       HubConfig
	Slack     SlackConfig
	Execution ExecutionConfig
	Time      TimeConfig
}

type HubConfig struct {
	URL       string
	NodeToken string
	Tag       string
	Timeout   time.Duration
}

type SlackConfig struct {
	WebhookURL string
	Timeout    time.Duration
}

type ExecutionConfig struct {
	Command []string
}

type TimeConfig struct {
	Name     string
	Location *time.Location
}

func Load(args []string) (*Config, []string, error) {
	fs := flag.NewFlagSet("jobboard", flag.ContinueOnError)
	var parseErr bytes.Buffer
	fs.SetOutput(&parseErr)

	hubURL := fs.String("hub-url", envString("JOBBOARD_HUB_URL", "http://localhost:8080"), "Hub base URL")
	nodeToken := fs.String("node-token", envString("JOBBOARD_NODE_TOKEN", ""), "Token for Hub node trigger API")
	tag := fs.String("tag", "", "Optional tag forwarded to Hub")
	slackWebhook := fs.String("slack-webhook", envString("JOBBOARD_SLACK_WEBHOOK", ""), "Slack incoming webhook URL")
	hubTimeout := fs.Duration("hub-timeout", envDuration("JOBBOARD_HUB_TIMEOUT", 60*time.Second), "Timeout for Hub API requests")
	slackTimeout := fs.Duration("slack-timeout", envDuration("JOBBOARD_SLACK_TIMEOUT", 10*time.Second), "Timeout for Slack API requests")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: jobboard [flags] -- <command> [args...]\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		msg := strings.TrimSpace(parseErr.String())
		if msg == "" {
			msg = err.Error()
		}
		if err == flag.ErrHelp {
			return nil, nil, errors.New(msg)
		}
		return nil, nil, fmt.Errorf("failed to parse flags: %s", msg)
	}

	command := fs.Args()
	if len(command) == 0 {
		return nil, nil, errors.New("execution command is required; pass it after `--`")
	}

	tzName, tzLocation := loadLocation()

	cfg := &Config{
		Hub: HubConfig{
			URL:       *hubURL,
			NodeToken: *nodeToken,
			Tag:       *tag,
			Timeout:   *hubTimeout,
		},
		Slack: SlackConfig{
			WebhookURL: *slackWebhook,
			Timeout:    *slackTimeout,
		},
		Execution: ExecutionConfig{
			Command: command,
		},
		Time: TimeConfig{
			Name:     tzName,
			Location: tzLocation,
		},
	}

	warnings := cfg.collectWarnings()
	if !cfg.Hub.Enabled() && !cfg.Slack.Enabled() {
		return nil, warnings, errors.New("either Slack webhook or Hub node token must be provided")
	}

	return cfg, warnings, nil
}

func (c HubConfig) Enabled() bool {
	return c.NodeToken != ""
}

func (c SlackConfig) Enabled() bool {
	return c.WebhookURL != ""
}

func (c *Config) collectWarnings() []string {
	var warnings []string
	if !c.Slack.Enabled() {
		warnings = append(warnings, "Slack webhook is not configured; skipping Slack notifications")
	}
	if !c.Hub.Enabled() {
		warnings = append(warnings, "Hub node token is not configured; skipping Hub integration")
	}
	return warnings
}

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func loadLocation() (string, *time.Location) {
	tz := envString("TIMEZONE", "")
	if tz == "" {
		tz = envString("TZ", "")
	}
	if tz == "" {
		tz = "Asia/Tokyo"
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("warn: failed to load timezone %q: %v; falling back to UTC", tz, err)
		return tz, time.UTC
	}
	return tz, loc
}
