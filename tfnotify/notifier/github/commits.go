package github

import (
	"context"
	"errors"

	"github.com/google/go-github/github"
)

// CommitsService handles communication with the commits related
// methods of GitHub API
type CommitsService service

// List lists commits on a repository
func (g *CommitsService) List(owner, repo, revision string) ([]string, error) {
	var s []string
	commits, _, err := g.client.Repositories.ListCommits(
		context.Background(),
		owner,
		repo,
		&github.CommitsListOptions{SHA: revision},
	)
	if err != nil {
		return s, err
	}
	for _, commit := range commits {
		s = append(s, *commit.SHA)
	}
	return s, nil
}

// Last returns the hash of the previous commit of the given commit
func (g *CommitsService) Last(owner, repo, revision string) (string, error) {
	if revision == "" {
		return "", errors.New("no revision specified")
	}
	commits, err := g.List(owner, repo, revision)
	if err != nil {
		return "", err
	}
	return commits[1], nil
}
