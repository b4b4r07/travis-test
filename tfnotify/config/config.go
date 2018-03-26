package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config is for tfnotify config structure
type Config struct {
	Repository Repository `yaml:"repository"`
	CI         string     `yaml:"ci"`
	Notifier   Notifier   `yaml:"notifier"`
	Terraform  Terraform  `yaml:"terraform"`

	path string
}

// Repository is a general repository structure such as GitHub
type Repository struct {
	Owner string `yaml:"owner"`
	Name  string `yaml:"name"`
}

// Notifier is a notification notifier
type Notifier struct {
	Github GithubNotifier `yaml:"github"`
	Slack  SlackNotifier  `yaml:"slack"`
}

// GithubNotifier is a notifier for GitHub
type GithubNotifier struct {
	Token string `yaml:"token"`
}

// SlackNotifier is a notifier for Slack
type SlackNotifier struct {
	Token   string `yaml:"token"`
	Channel string `yaml:"channel"`
	Bot     string `yaml:"bot"`
}

// Terraform represents terraform configurations
type Terraform struct {
	Default Default `yaml:"default"`
	Fmt     Fmt     `yaml:"fmt"`
	Plan    Plan    `yaml:"plan"`
	Apply   Apply   `yaml:"apply"`
}

// Default is a default setting for terraform commands
type Default struct {
	Template string `yaml:"template"`
}

// Fmt is a terraform fmt config
type Fmt struct {
	Template string `yaml:"template"`
}

// Plan is a terraform plan config
type Plan struct {
	Template string `yaml:"template"`
}

// Apply is a terraform apply config
type Apply struct {
	Template string `yaml:"template"`
}

// LoadFile binds the config file to Config structure
func (cfg *Config) LoadFile(path string) error {
	cfg.path = path
	_, err := os.Stat(cfg.path)
	if err != nil {
		return fmt.Errorf("%s: no config file", cfg.path)
	}
	raw, _ := ioutil.ReadFile(cfg.path)
	return yaml.Unmarshal(raw, cfg)
}

// Validation validates config file
func (cfg *Config) Validation() error {
	if cfg.Repository.Owner == "" {
		return fmt.Errorf("repository owner is missing")
	}
	if cfg.Repository.Name == "" {
		return fmt.Errorf("repository name is missing")
	}
	switch strings.ToLower(cfg.CI) {
	case "":
		return errors.New("ci: need to be set")
	case "circleci", "circle-ci":
		// ok pattern
	default:
		return fmt.Errorf("%s: not supported yet", cfg.CI)
	}
	notifier := cfg.GetNotifierType()
	if notifier == "" {
		return fmt.Errorf("notifier is missing")
	}
	return nil
}

func (cfg *Config) isDefinedGithub() bool {
	// not empty
	return cfg.Notifier.Github != (GithubNotifier{})
}

func (cfg *Config) isDefinedSlack() bool {
	// not empty
	return cfg.Notifier.Slack != (SlackNotifier{})
}

// GetNotifierType return notifier type described in Config
func (cfg *Config) GetNotifierType() string {
	if cfg.isDefinedGithub() {
		return "github"
	}
	if cfg.isDefinedSlack() {
		return "slack"
	}
	return ""
}
