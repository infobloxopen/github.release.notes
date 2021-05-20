package githubclient

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ReleaseNotesData struct {
	Tag           string
	Branch        string
	Comment       string
	Date          time.Time
	ChangeLogLink string
	Commits       []Commit
}

type Commit struct {
	Title       string
	Author      string
	AuthorLink  string
	PullRequest int
}

func (rnd *ReleaseNotesData) PrepareReleaseNotesMessage(logger *logrus.Logger) (string, string, error) {
	if rnd == nil {
		return "", "", fmt.Errorf("error while preparation release notes message: release notes are (nil)")
	}

	releaseTitle := rnd.prepareTitle()
	releaseBody, err := rnd.prepareBody(logger)
	if err != nil {
		return "", "", fmt.Errorf("error while release notes body preparation: %v for data %v", err, *rnd)
	}

	return releaseTitle, releaseBody, nil
}

func (rnd *ReleaseNotesData) prepareTitle() string {
	return "[" + rnd.Branch + "] " + rnd.Tag + " (" + rnd.Date.Format("2006-01-02") + ")"
}

func (rnd *ReleaseNotesData) prepareBody(logger *logrus.Logger) (string, error) {
	resp := fmt.Sprintf("[Full Changelog](%s)\n\n**New commits and merged pull requests:**", rnd.ChangeLogLink)
	if rnd.Commits != nil {
		for _, v := range rnd.Commits {
			commit := v.Title
			if v.PullRequest != 0 {
				prLink, err := url.Parse("https://github.com")
				if err != nil {
					logger.Errorf("error while PR link creation: %v", err)
				} else {
					prLink.Path = path.Join(prLink.Path, viper.GetString("github.org"))
					prLink.Path = path.Join(prLink.Path, viper.GetString("github.repo"))
					prLink.Path = path.Join(prLink.Path, "pull")
					prLink.Path = path.Join(prLink.Path, strconv.Itoa(v.PullRequest))
					commit += " [#" + strconv.Itoa(v.PullRequest) + "](" + prLink.String() + ")"
				}
			}
			if v.Author != "" && v.AuthorLink != "" {
				commit += " ([" + v.Author + "](" + v.AuthorLink + "))"
			}
			resp = fmt.Sprintf("%s\n- %s", resp, commit)
		}
	}
	return resp, nil
}
