package githubclient

import (
	"context"
	"fmt"

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
	//tagsList, _, err := gc.client.Repositories.ListTags(context.Background(), gc.OrgName, gc.RepoName, nil)
	tagRef, _, err := gc.client.Git.ListMatchingRefs(context.Background(), gc.OrgName, gc.RepoName, &github.ReferenceListOptions{Ref: "tags"})
	//repos, _, err := gc.client.Repositories.ListByOrg(context.Background(), org, nil)
	if err != nil {
		return nil, err
	}
	if len(tagRef) == 0 {
		return nil, fmt.Errorf("no tags were found")
	}
	rnd := make([]ReleaseNotesData, 5)
	for i, tag := range tagRef {
		log.Infoln(tag)
		log.Infoln(tag.GetObject().GetSHA())
		log.Infoln(tag.GetRef())
		changeLogLink := ""
		if i < len(tagRef)-1 {
			previousTag := tagRef[i+1].GetRef()
			changeLogLink = fmt.Sprintf("https://github.com/%s/%s/compare/%s...%s", gc.OrgName, gc.RepoName, previousTag, tag.GetRef())
		}
		//tagData, _, err := gc.client.Git.GetTag(context.Background(), gc.OrgName, gc.RepoName, tag.Commit.GetSHA())
		//if err != nil {
		//	return nil, err
		//}
		//log.Infof("!!! %s !!!", tag.GetCommit().GetAuthor().GetDate().String)
		//log.Infof("!!! %s !!!", tagData.GetTagger().GetDate().String)
		rnd = append(rnd,
			ReleaseNotesData{Tag: tag.GetRef(),
				Comment: "", //tag.GetCommit().GetMessage(),
				//Date:          "",//tag.GetCommit().GetAuthor().GetDate(),
				ChangeLogLink: changeLogLink,
			})
	}
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
