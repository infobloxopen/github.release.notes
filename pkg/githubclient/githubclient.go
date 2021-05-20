package githubclient

import (
	"context"

	"github.com/google/go-github/v35/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type githubClient struct {
	GithubToken string
	RepoURL     string
	client      *github.Client
}

type GithubClientClient interface {
	GetReleaseNotesData(repo string) ([]ReleaseNotesData, error)
	PublishReleaseNotes(rndList []ReleaseNotesData)
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
func (gc *githubClient) GetReleaseNotesData(repo string) ([]ReleaseNotesData, error) {
	gc.RepoURL = repo
	repos, _, err := gc.client.Repositories.ListByOrg(context.Background(), "infobloxopen", nil)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("repos: %v", repos)
	rnd := make([]ReleaseNotesData, 5)
	return rnd, nil
}

// GetFeatureFlagValue return feature flag value
func (gc *githubClient) PublishReleaseNotes(rndList []ReleaseNotesData) {
	for _, v := range rndList {
		title, body := v.PrepareReleaseNotesMessage()

		release := &github.RepositoryRelease{
			TagName:         &v.Tag,
			TargetCommitish: &v.Branch,
			Name:            &title,
			Body:            &body,
		}
		_, _, err := gc.client.Repositories.CreateRelease(context.Background(), viper.GetString("github.org"), viper.GetString("github.repo"), release)
		if err != nil {
			log.Errorf("Error while publishing release notes: %v", err)
			continue
		}
	}
}
