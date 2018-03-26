package github

import (
	"github.com/kouzoh/tfnotify/terraform"
)

// NotifyService handles communication with the notification related
// methods of GitHub API
type NotifyService service

// Notify posts comment optimized for notifications
func (g *NotifyService) Notify(body string) (exit int, err error) {
	cfg := g.client.Config
	parser := g.client.Config.Parser
	template := g.client.Config.Template

	result := parser.Parse(body)
	if result.Error != nil {
		return result.ExitCode, result.Error
	}
	if result.Result == "" {
		return result.ExitCode, result.Error
	}

	template.SetValue(terraform.CommonTemplate{
		Message: cfg.PR.Message,
		Result:  result.Result,
		Body:    body,
	})
	body, err = template.Execute()
	if err != nil {
		return result.ExitCode, err
	}

	value := template.GetValue()

	if cfg.PR.IsNumber() {
		if err := g.client.Comment.DeleteDuplicates(value.Title); err != nil {
			return result.ExitCode, err
		}
	}

	_, isApply := parser.(*terraform.ApplyParser)
	if !cfg.PR.IsNumber() && isApply {
		lastRevision, _ := g.client.Commits.Last(cfg.Owner, cfg.Repo, cfg.PR.Revision)
		cfg.PR.Revision = lastRevision
	}

	return result.ExitCode, g.client.Comment.Post(cfg.Owner, cfg.Repo, body, PostOptions{
		Number:   cfg.PR.Number,
		Revision: cfg.PR.Revision,
	})
}
