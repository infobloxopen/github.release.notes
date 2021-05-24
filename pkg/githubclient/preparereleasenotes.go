package githubclient

import (
	"fmt"
	"time"
)

// ReleaseNotesData is a base struct for release notes
type ReleaseNotesData struct {
	Tag           string
	Branch        string
	Comment       string
	Date          time.Time
	ChangeLogLink string
	Commits       []CommitData
	releaseID     int64
}

// CommitData contains necessary information about commit data
type CommitData struct {
	Message string
	Author  string
	URL     string
}

// PrepareReleaseNotesMessage prepares full information about tag or list of tags
func (rnd *ReleaseNotesData) PrepareReleaseNotesMessage() (string, string) {
	releaseTitle := rnd.prepareTitle()
	releaseBody := rnd.prepareBody()

	return releaseTitle, releaseBody
}

func (rnd *ReleaseNotesData) prepareTitle() string {
	title := ""
	if rnd.Branch != "" {
		title = "[" + rnd.Branch + "] "
	}
	return title + "[" + rnd.Tag + "] " + rnd.Comment + " (" + rnd.Date.Format("2006-01-02") + ")"
}

// Body template:
//
// [Full Changelog]
//
// New commits and merged pull requests:
//
// - commit_N [#N] (author_N)
// - ...
// - commit_2 [#2] (author_2)
// - commit_1 [#1] (author_1)
func (rnd *ReleaseNotesData) prepareBody() string {
	resp := ""
	if rnd.ChangeLogLink != "" {
		resp = fmt.Sprintf("[Full Changelog](%s)", rnd.ChangeLogLink)
	}
	if rnd.Commits != nil {
		if resp != "" {
			resp += "\n\n"
		}
		resp += "**New commits and merged pull requests:**\n"
		for _, v := range rnd.Commits {
			// log.Debugf("%v", v)
			commit := v.Message
			// prLink, err := url.Parse("https://github.com")
			// if err != nil {
			// 	log.Errorf("error while PR link creation: %v", err)
			// } else {
			// 	prLink.Path = path.Join(prLink.Path, viper.GetString("github.org"))
			// 	prLink.Path = path.Join(prLink.Path, viper.GetString("github.repo"))
			// 	prLink.Path = path.Join(prLink.Path, "pull")
			// 	prLink.Path = path.Join(prLink.Path, strconv.Itoa(v.PullRequest))
			// 	commit += " [#" + strconv.Itoa(v.PullRequest) + "](" + prLink.String() + ")"
			// }
			commit = fmt.Sprintf("%s ([%s](%s))", commit, v.Author, v.URL)
			resp = fmt.Sprintf("%s\n- %s", resp, commit)
		}
	}
	return resp
}
