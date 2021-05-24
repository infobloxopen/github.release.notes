package githubclient

import (
	"bytes"
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
	CommitAuthor    string
	CommitAuthorURL string
	CommitDate      string
	CommitMessage   string
	CommitPR        string
	CommitURL       string
}

// PrepareReleaseNotesMessage prepares full information about tag or list of tags
func (rnd *ReleaseNotesData) PrepareReleaseNotesMessage() (string, string) {
	releaseTitle := rnd.prepareTitle()
	releaseBody := rnd.prepareBody()

	return releaseTitle, releaseBody
}

func (rnd *ReleaseNotesData) prepareTitle() string {
	titleTmpl := template.Must(template.ParseFiles(viper.GetString("template.title")))
	var title bytes.Buffer
	err := titleTmpl.Execute(&title, rnd)
	if err != nil {
		log.Errorf("Error while title template %s rendering: %v", viper.GetString("template.title"), err)
		return ""
	}
	return title.String()
}

func (rnd *ReleaseNotesData) prepareBody() string {
	bodyTmpl := template.Must(template.ParseFiles(viper.GetString("template.body")))
	var releaseBody bytes.Buffer
	err := bodyTmpl.Execute(&releaseBody, rnd)
	if err != nil {
		log.Errorf("Error while body template %s rendering: %v", viper.GetString("template.body"), err)
		return ""
	}
	return releaseBody.String()
}
