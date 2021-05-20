package githubclient

import (
	"context"

	"github.com/google/go-github/v35/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type githubClient struct {
	GithubToken string
	RepoURL     string
	client      *github.Client
}

type GithubClientClient interface {
	PrepeareReleaseNotes(repo string) error
}

// New creates a client wrapper
func NewGithubClient(token string) GithubClientClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	return &githubClient{
		client: client,
	}
}

// GetFeatureFlagValue return feature flag value
func (gc *githubClient) PrepeareReleaseNotes(repo string) error {
	gc.RepoURL = repo
	repos, _, err := gc.client.Repositories.ListByOrg(context.Background(), "infobloxopen", nil)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("repos: %v", repos)
	return nil
}
