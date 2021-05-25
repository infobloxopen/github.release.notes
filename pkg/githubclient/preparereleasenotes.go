package githubclient

import (
	"bytes"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ReleaseNotesData is a base struct for release notes
type ReleaseNotesData struct {
	Branch        string
	ChangeLogLink string
	Commits       []CommitData
	Tag           string
	TagComment    string
	TagDate       string
	releaseID     int64
}

// CommitData contains necessary information about commit data
type CommitData struct {
	CommitAdditions int
	CommitAuthor    string
	CommitAuthorURL string
	CommitDate      string
	CommitDeletions int
	CommitMessage   string
	CommitPR        string
	CommitURL       string
	SquashCommits   []string
}

// PrepareReleaseNotesMessage prepares full information about tag or list of tags
func (rnd *ReleaseNotesData) PrepareReleaseNotesMessage() (string, string) {
	releaseTitle := rnd.prepareTitle()
	releaseBody := rnd.prepareBody()

	return releaseTitle, releaseBody
}

func (rnd *ReleaseNotesData) prepareTitle() string {
	var titleTmpl *template.Template
	_, err := os.Stat(viper.GetString("template.title"))
	if err != nil {
		log.Errorf("Error with template title %s: %v. Will use default \"in memory\" template.",
			viper.GetString("template.title"), err)
		titleTmpl = template.Must(template.New("title").Parse(defaultTitle))
	} else {
		titleTmpl = template.Must(template.ParseFiles(viper.GetString("template.title")))
	}
	var title bytes.Buffer
	err = titleTmpl.Execute(&title, rnd)
	if err != nil {
		log.Errorf("Error while title template rendering: %v", err)
		return ""
	}
	return title.String()
}

func (rnd *ReleaseNotesData) prepareBody() string {
	var bodyTmpl *template.Template
	_, err := os.Stat(viper.GetString("template.body"))
	if err != nil {
		log.Errorf("Error with template body %s: %v. Will use default \"in memory\" template.",
			viper.GetString("template.body"), err)
		bodyTmpl = template.Must(template.New("body").Parse(defaultBody))
	} else {
		bodyTmpl = template.Must(template.ParseFiles(viper.GetString("template.body")))
	}
	var releaseBody bytes.Buffer
	err = bodyTmpl.Execute(&releaseBody, rnd)
	if err != nil {
		log.Errorf("Error while body template rendering: %v", err)
		return ""
	}
	return releaseBody.String()
}
