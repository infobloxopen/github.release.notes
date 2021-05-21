package githubclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v35/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	releases, _, err := gc.client.Repositories.ListReleases(context.Background(), gc.OrgName, gc.RepoName, nil)
	if err != nil {
		return nil, err
	}
	var rnd []ReleaseNotesData
	for i, tag := range tagsList {
		tagData, _, err := gc.client.Git.GetTag(context.Background(), gc.OrgName, gc.RepoName, tag.GetObject().GetSHA())
		if err != nil {
			return nil, err
		}
		var releaseID int64
		for _, release := range releases {
			if release.GetTagName() == tagData.GetTag() {
				releaseID = release.GetID()
				continue
			}
		}
		if !viper.GetBool("update.exist") && releaseID != 0 {
			log.Debugf("Skipping tag: %v", tagData.GetTag())
			continue
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
						commitMsg := i.GetCommit().GetMessage()
						commitMsg = strings.Split(commitMsg, "\n")[0]
						commitMsg = strings.ReplaceAll(commitMsg, "(", "")
						commitMsg = strings.ReplaceAll(commitMsg, ")", "")
						commits = append([]CommitData{
							{
								Author:  i.GetCommit().GetAuthor().GetName(),
								Message: commitMsg,
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
				releaseID:     releaseID,
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
		if v.releaseID != 0 {
			_, err := gc.client.Repositories.DeleteRelease(context.Background(), gc.OrgName, gc.RepoName, v.releaseID)
			if err != nil {
				log.Errorf("Error while deleting release notes: %v", err)
				continue
			}
		}
		_, _, err := gc.client.Repositories.CreateRelease(context.Background(), gc.OrgName, gc.RepoName, release)
		if err != nil {
			log.Errorf("Error while publishing release notes: %v", err)
			continue
		}
	}
}
