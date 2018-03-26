package github

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/github"
)

// CommentService handles communication with the comment related
// methods of GitHub API
type CommentService service

// PostOptions specifies the optional parameters to post comments to a pull request
type PostOptions struct {
	Number   int
	Revision string
}

// Post posts comment
func (g *CommentService) Post(owner, repo, body string, opt PostOptions) error {
	if opt.Number != 0 {
		_, _, err := g.client.Issues.CreateComment(
			context.Background(),
			owner,
			repo,
			opt.Number,
			&github.IssueComment{Body: &body},
		)
		return err
	}
	if opt.Revision != "" {
		_, _, err := g.client.Repositories.CreateComment(
			context.Background(),
			owner,
			repo,
			opt.Revision,
			&github.RepositoryComment{Body: &body},
		)
		return err
	}
	return fmt.Errorf("PR number or Revision is required")
}

// List lists comments on GitHub issues/pull requests
func (g *CommentService) List(owner, repo string, number int) ([]*github.IssueComment, error) {
	comments, _, err := g.client.Issues.ListComments(
		context.Background(),
		owner,
		repo,
		number,
		&github.IssueListCommentsOptions{},
	)
	return comments, err
}

// Delete deletes comment on GitHub issues/pull requests
func (g *CommentService) Delete(owner, repo string, id int) error {
	_, err := g.client.Issues.DeleteComment(
		context.Background(),
		owner,
		repo,
		id,
	)
	return err
}

// DeleteDuplicates deletes duplicate comments containing arbitrary character strings
func (g *CommentService) DeleteDuplicates(title string) error {
	cfg := g.client.Config
	re := regexp.MustCompile(`(?m)^(\n+)?` + title + `\n+` + cfg.PR.Message + `\n+`)
	comments, err := g.client.Comment.List(cfg.Owner, cfg.Repo, cfg.PR.Number)
	if err != nil {
		return err
	}
	var ids []int64
	for _, comment := range comments {
		if re.MatchString(*comment.Body) {
			ids = append(ids, *comment.ID)
		}
	}
	for _, id := range ids {
		err := g.client.Comment.Delete(cfg.Owner, cfg.Repo, int(id))
		if err != nil {
			// TODO do not return error
			return err
		}
	}
	return nil
}
