package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultHubURL = "http://localhost:8080"

	envHubURL          = "JOBBOARD_HUB_URL"
	envNodeToken       = "JOBBOARD_NODE_TOKEN"
	envSlackWebhookURL = "JOBBOARD_SLACK_WEBHOOK_URL"
)

type Config struct {
	Hub     HubConfig
	Slack   SlackConfig
	Command CommandConfig
}

type HubConfig struct {
	Enabled   bool
	BaseURL   string
	NodeToken string
	Tag       *string
	Timeout   time.Duration
}

type SlackConfig struct {
	WebhookURL string
}

type CommandConfig struct {
	Args []string
}

func Load(args []string) (Config, []string, error) {
	var warnings []string

	fs := flag.NewFlagSet("jobboard", flag.ContinueOnError)
	hubURL := fs.String("hub-url", "", "Hub base URL (fallback: $"+envHubURL+")")
	nodeToken := fs.String("node-token", "", "Node token for hub trigger (fallback: $"+envNodeToken+")")
	tag := fs.String("tag", "", "Optional tag sent to the hub")
	slackWebhook := fs.String("slack-webhook", "", "Slack webhook URL (fallback: $"+envSlackWebhookURL+")")
	timeout := fs.Duration("timeout", 60*time.Second, "Timeout for hub requests")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: jobboard [flags] -- <command> [args...]\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return Config{}, warnings, err
	}

	cmdArgs := fs.Args()
	if len(cmdArgs) == 0 {
		return Config{}, warnings, errors.New("missing command to execute (provide arguments after --)")
	}

	slack := firstNonEmpty(*slackWebhook, os.Getenv(envSlackWebhookURL))
	if strings.TrimSpace(slack) == "" {
		warnings = append(warnings, "warning: slack webhook is not provided; aborting execution")
	}

	baseURL := firstNonEmpty(*hubURL, os.Getenv(envHubURL), defaultHubURL)
	nodeToken := firstNonEmpty(*nodeToken, os.Getenv(envNodeToken))

	hubCfg := HubConfig{
		Enabled:   true,
		BaseURL:   baseURL,
		NodeToken: nodeToken,
		Timeout:   *timeout,
	}

	if trimmed := strings.TrimSpace(hubCfg.NodeToken); trimmed == "" {
		hubCfg.Enabled = false
		warnings = append(warnings, "warning: node token is not provided; hub integration is disabled")
	}

	cfg := Config{
		Hub: hubCfg,
		Slack: SlackConfig{
			WebhookURL: slack,
		},
		Command: CommandConfig{
			Args: cmdArgs,
		},
	}
	return cfg, warnings, nil
}
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
