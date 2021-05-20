package githubclient

import (
	"context"

	"github.com/google/go-github/v35/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type githubClient struct {
	GithubToken string
	OrgName     string
	RepoName    string
	client      *github.Client
}

type GithubClientClient interface {
	GetReleaseNotesData() ([]ReleaseNotesData, error)
	PublishReleaseNotes(rndList []ReleaseNotesData)
}

// New creates a client wrapper
func NewGithubClient(token, org, repo string) GithubClientClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	return &githubClient{
		OrgName:  org,
		RepoName: repo,
		client:   client,
	}
}

// GetReleaseNotesData return release notes data collected
func (gc *githubClient) GetReleaseNotesData() ([]ReleaseNotesData, error) {
	tagsList, _, err := gc.client.Repositories.ListTags(context.Background(), gc.OrgName, gc.RepoName, &github.ListOptions{Page: 0, PerPage: 1000})
	//repos, _, err := gc.client.Repositories.ListByOrg(context.Background(), org, nil)
	if err != nil {
		log.Error(err)
	}
	for _, tag := range tagsList {
		log.Debugf("tag: %v", *tag.Name)
	}
	rnd := make([]ReleaseNotesData, 5)
	return rnd, nil
}

// PublishReleaseNotes publishes release notes to GitHub
func (gc *githubClient) PublishReleaseNotes(rndList []ReleaseNotesData) {
	for _, v := range rndList {
		title, body := v.PrepareReleaseNotesMessage()

		release := &github.RepositoryRelease{
			TagName:         &v.Tag,
			TargetCommitish: &v.Branch,
			Name:            &title,
			Body:            &body,
		}
		log.Debugf("release: %v", release)
		// _, _, err := gc.client.Repositories.CreateRelease(context.Background(), viper.GetString("github.org"), viper.GetString("github.repo"), release)
		// if err != nil {
		// 	log.Errorf("Error while publishing release notes: %v", err)
		// 	continue
		// }
	}
}
