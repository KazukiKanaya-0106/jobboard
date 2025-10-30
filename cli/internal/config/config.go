package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
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

	slack := firstNonEmpty(*slackWebhook, getEnv(envSlackWebhookURL))

}

func getEnvOr(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
