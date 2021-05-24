package githubclient

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v35/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	MaxPerPage = 1000
)

// GithubClientOptions structure holds necessary options for Github Client
type GithubClientOptions struct {
	Token      string
	Owner      string
	Repository string
}
type githubClient struct {
	OrgName  string
	RepoName string
	client   *github.Client
}

type GithubClientClient interface {
	GetReleaseNotesData(tag string) ([]ReleaseNotesData, error)
	PublishReleaseNotes(rndList []ReleaseNotesData)
}

// New creates a client wrapper
func NewGithubClient(opt GithubClientOptions) GithubClientClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: opt.Token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	return &githubClient{
		OrgName:  opt.Owner,
		RepoName: opt.Repository,
		client:   client,
	}
}

// GetReleaseNotesData return release notes data collected
func (gc *githubClient) GetReleaseNotesData(githubTag string) ([]ReleaseNotesData, error) {
	rlo := &github.ReferenceListOptions{
		Ref:         "tags",
		ListOptions: github.ListOptions{PerPage: MaxPerPage},
	}
	var tagsList []*github.Reference
	for {
		tags, resp, err := gc.client.Git.ListMatchingRefs(context.Background(), gc.OrgName, gc.RepoName, rlo)
		if err != nil {
			return nil, err
		}
		tagsList = append(tagsList, tags...)
		if resp.NextPage == 0 {
			break
		}
		rlo.Page = resp.NextPage
	}
	if len(tagsList) == 0 {
		return nil, fmt.Errorf("no tags were found")
	}

	lo := &github.ListOptions{PerPage: MaxPerPage}
	var releases []*github.RepositoryRelease
	for {
		rls, resp, err := gc.client.Repositories.ListReleases(context.Background(), gc.OrgName, gc.RepoName, lo)
		if err != nil {
			return nil, err
		}
		releases = append(releases, rls...)
		if resp.NextPage == 0 {
			break
		}
		lo.Page = resp.NextPage
	}

	var rnd []ReleaseNotesData
	for i, tag := range tagsList {
		if tag.GetObject().GetType() != "tag" {
			log.Errorf("Tag %v is not annotated: %v", tag.GetRef(), tag.GetObject().GetType())
			continue
		}
		tagData, _, err := gc.client.Git.GetTag(context.Background(), gc.OrgName, gc.RepoName, tag.GetObject().GetSHA())
		if err != nil {
			log.Errorf("Error while tag processing: %v", err)
			continue
		}

		// check if the tag is set and will release just the given tag
		if *tagData.Tag != githubTag && githubTag != "" {
			continue
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
						re := regexp.MustCompile(`\(#\d+\)`)
						commitMsg = re.ReplaceAllStringFunc(commitMsg, repl)
						commits = append([]CommitData{
							{
								Author:  i.GetAuthor().GetLogin(),
								Message: commitMsg,
								URL:     i.GetAuthor().GetURL(),
							},
						}, commits...)
					}
				}
			}
		} else {
			clo := &github.CommitsListOptions{
				SHA:         tagData.GetObject().GetSHA(),
				ListOptions: github.ListOptions{PerPage: MaxPerPage},
			}
			var commitsList []*github.RepositoryCommit
			for {
				comtsLst, resp, err := gc.client.Repositories.ListCommits(context.Background(), gc.OrgName, gc.RepoName, clo)
				if err != nil {
					log.Errorf("Error getting commits for the tag %v: %v", tagData.GetTag(), err)
					commitsList = nil
					break
				}
				commitsList = append(commitsList, comtsLst...)
				if resp.NextPage == 0 {
					break
				}
				clo.Page = resp.NextPage
			}
			for _, i := range commitsList {
				commitMsg := i.GetCommit().GetMessage()
				commitMsg = strings.Split(commitMsg, "\n")[0]
				re := regexp.MustCompile(`\(#\d+\)`)
				commitMsg = re.ReplaceAllStringFunc(commitMsg, repl)
				commits = append(commits, CommitData{
					Author:  i.GetAuthor().GetLogin(),
					Message: commitMsg,
					URL:     i.GetAuthor().GetURL(),
				})
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

func repl(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, ")", ""), "(", "")
}
