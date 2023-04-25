package release_fetcher

import (
	"github.com/mrlutik/autoflowhub/internal/adapters"
	"github.com/mrlutik/autoflowhub/internal/models"
)

type ReleaseFetcher interface {
	GetLatestRelease(owner, repo string) (models.Release, error)
}

type GitHubReleaseFetcher struct {
	Adapter adapters.GitHubAdapter
}

func NewGitHubReleaseFetcher(adapter adapters.GitHubAdapter) *GitHubReleaseFetcher {
	return &GitHubReleaseFetcher{
		Adapter: adapter,
	}
}

func (g *GitHubReleaseFetcher) GetLatestRelease(owner, repo string) (models.Release, error) {
	release, err := g.Adapter.GetLatestRelease(owner, repo)
	if err != nil {
		return models.Release{}, err
	}

	tagName := ""
	if release.TagName != nil {
		tagName = *release.TagName
	}

	return models.Release{
		TagName: tagName,
	}, nil

}
