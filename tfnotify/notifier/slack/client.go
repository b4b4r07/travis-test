package slack

import (
	"context"
	"errors"

	"github.com/kouzoh/tfnotify/terraform"
	"github.com/lestrrat-go/slack"
)

// Client is a API client for Slack
type Client struct {
	*slack.Client

	Config Config

	common service

	Notify *NotifyService
}

// Config is a configuration for GitHub client
type Config struct {
	Token    string
	Channel  string
	Context  context.Context
	Botname  string
	Message  string
	CI       string
	Parser   terraform.Parser
	Template terraform.Template
}

type service struct {
	client *Client
}

// NewClient returns Client initialized with Config
func NewClient(cfg Config) (*Client, error) {
	if cfg.Token == "" {
		return &Client{}, errors.New("token is missing")
	}
	c := &Client{
		Config: cfg,
		Client: slack.New(cfg.Token),
	}
	c.common.client = c
	c.Notify = (*NotifyService)(&c.common)
	return c, nil
}
