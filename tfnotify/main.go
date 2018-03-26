package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kouzoh/tfnotify/config"
	"github.com/kouzoh/tfnotify/notifier"
	"github.com/kouzoh/tfnotify/notifier/github"
	"github.com/kouzoh/tfnotify/notifier/slack"
	"github.com/kouzoh/tfnotify/terraform"

	"github.com/urfave/cli"
)

const (
	defaultConfigPath = "tfnotify.yaml"
	defaultNotifier   = "github"
	defaultCI         = "circleci"
)

type tfnotify struct {
	config   config.Config
	context  *cli.Context
	parser   terraform.Parser
	template terraform.Template
}

// Run sends the notification with notifier
func (t *tfnotify) Run() error {
	ciname := t.config.CI
	if t.context.GlobalString("ci") != "" {
		ciname = t.context.GlobalString("ci")
	}
	ciname = strings.ToLower(ciname)
	var ci CI
	var err error
	switch ciname {
	case "circleci", "circle-ci":
		ci, err = circleci()
		if err != nil {
			return err
		}
	case "":
		return fmt.Errorf("CI service: required (e.g. circleci)")
	default:
		return fmt.Errorf("CI service %v: not supported yet", ci)
	}

	selectedNotifier := t.config.GetNotifierType()
	if t.context.GlobalString("notifier") != "" {
		selectedNotifier = t.context.GlobalString("notifier")
	}

	var notifier notifier.Notifier
	switch selectedNotifier {
	case "github":
		client, err := github.NewClient(github.Config{
			Token: t.config.Notifier.Github.Token,
			Owner: t.config.Repository.Owner,
			Repo:  t.config.Repository.Name,
			PR: github.PullRequest{
				Revision: ci.PR.Revision,
				Number:   ci.PR.Number,
				Message:  t.context.String("message"),
			},
			CI:       ci.URL,
			Parser:   t.parser,
			Template: t.template,
		})
		if err != nil {
			return err
		}
		notifier = client.Notify
	case "slack":
		client, err := slack.NewClient(slack.Config{
			Token:    t.config.Notifier.Slack.Token,
			Channel:  t.config.Notifier.Slack.Channel,
			Context:  context.Background(),
			Botname:  t.config.Notifier.Slack.Bot,
			Message:  t.context.String("message"),
			CI:       ci.URL,
			Parser:   t.parser,
			Template: t.template,
		})
		if err != nil {
			return err
		}
		notifier = client.Notify
	case "":
		return fmt.Errorf("notifier is missing")
	default:
		return fmt.Errorf("%s: not supported notifier yet", t.context.GlobalString("notifier"))
	}

	if notifier == nil {
		return fmt.Errorf("no notifier specified at all")
	}

	return NewExitError(notifier.Notify(tee()))
}

func main() {
	app := cli.NewApp()
	app.Name = "tfnotify"
	app.Usage = "Notify the execution result of terraform command"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "ci", Usage: "name of CI to run tfnotify"},
		cli.StringFlag{Name: "config", Usage: "config path"},
		cli.StringFlag{Name: "notifier", Usage: "notification destination"},
	}
	app.Commands = []cli.Command{
		{
			Name:   "fmt",
			Usage:  "Parse stdin as a fmt result",
			Action: cmdFmt,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "message, m",
					Usage: "Specify the message to use for notification",
				},
			},
		},
		{
			Name:   "plan",
			Usage:  "Parse stdin as a plan result",
			Action: cmdPlan,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "message, m",
					Usage: "Specify the message to use for notification",
				},
			},
		},
		{
			Name:   "apply",
			Usage:  "Parse stdin as a apply result",
			Action: cmdApply,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "message, m",
					Usage: "Specify the message to use for notification",
				},
			},
		},
	}

	err := app.Run(os.Args)
	os.Exit(RunExitCoder(err))
}

func newConfig(ctx *cli.Context) (cfg config.Config, err error) {
	confPath := ctx.String("config")
	if confPath == "" {
		confPath = defaultConfigPath
	}
	if err := cfg.LoadFile(confPath); err != nil {
		return cfg, err
	}
	if err := cfg.Validation(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func cmdFmt(ctx *cli.Context) error {
	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}
	t := &tfnotify{
		config:   cfg,
		context:  ctx,
		parser:   terraform.NewFmtParser(),
		template: terraform.NewFmtTemplate(cfg.Terraform.Fmt.Template),
	}
	return t.Run()
}

func cmdPlan(ctx *cli.Context) error {
	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}
	t := &tfnotify{
		config:   cfg,
		context:  ctx,
		parser:   terraform.NewPlanParser(),
		template: terraform.NewPlanTemplate(cfg.Terraform.Plan.Template),
	}
	return t.Run()
}

func cmdApply(ctx *cli.Context) error {
	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}
	t := &tfnotify{
		config:   cfg,
		context:  ctx,
		parser:   terraform.NewApplyParser(),
		template: terraform.NewApplyTemplate(cfg.Terraform.Apply.Template),
	}
	return t.Run()
}
