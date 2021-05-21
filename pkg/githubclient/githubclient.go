package githubclient

import (
	"context"
	"fmt"
	"strings"

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
	tagsList, _, err := gc.client.Git.ListMatchingRefs(context.Background(), gc.OrgName, gc.RepoName, &github.ReferenceListOptions{Ref: "tags"})
	if err != nil {
		return nil, err
	}
	if len(tagsList) == 0 {
		return nil, fmt.Errorf("no tags were found")
	}
	var rnd []ReleaseNotesData
	for i, tag := range tagsList {
		tagData, _, err := gc.client.Git.GetTag(context.Background(), gc.OrgName, gc.RepoName, tag.GetObject().GetSHA())
		if err != nil {
			return nil, err
		}
		changeLogLink := ""
		var previousTag *github.Tag
		var commits []CommitData
		if i > 0 {
			previousTag, _, err = gc.client.Git.GetTag(context.Background(), gc.OrgName, gc.RepoName, tagsList[i-1].GetObject().GetSHA())
			if err == nil {
				changeLogLink = fmt.Sprintf("https://github.com/%s/%s/compare/%s...%s",
					gc.OrgName, gc.RepoName, previousTag.GetTag(), tagData.GetTag())
				tagCompare, _, _ := gc.client.Repositories.CompareCommits(context.Background(), gc.OrgName, gc.RepoName, previousTag.GetTag(), tagData.GetTag())
				if tagCompare != nil {
					for _, i := range tagCompare.Commits {
						commits = append([]CommitData{
							{
								Author:  i.GetCommit().GetAuthor().GetName(),
								Message: strings.Replace(strings.Replace(i.GetCommit().GetMessage(), "(", "", -1), ")", "", -1),
								URL:     i.GetAuthor().GetURL(),
							},
						}, commits...)
					}
				}
			}

		}
		rnd = append(rnd,
			ReleaseNotesData{Tag: tagData.GetTag(),
				Comment:       tagData.GetMessage(),
				Date:          tagData.GetTagger().GetDate(),
				ChangeLogLink: changeLogLink,
				Commits:       commits,
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
		_, _, err := gc.client.Repositories.CreateRelease(context.Background(), gc.OrgName, gc.RepoName, release)
		if err != nil {
			log.Errorf("Error while publishing release notes: %v", err)
			continue
		}
	}
}
