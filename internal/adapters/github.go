package adapters

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubAdapter is a struct to hold the GitHub client
type GitHubAdapter struct {
	client *github.Client
}

// NewGitHubAdapter initializes a new GitHubAdapter instance
func NewGitHubAdapter(accessToken string) *GitHubAdapter {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &GitHubAdapter{
		client: github.NewClient(tc),
	}
}

// GetLatestRelease fetches the latest release from the specified repository
func (gh *GitHubAdapter) GetLatestRelease(owner, repo string) (*github.RepositoryRelease, error) {
	release, _, err := gh.client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		return nil, err
	}
	return release, nil
}
