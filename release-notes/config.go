package main

import "github.com/spf13/pflag"

const (
	// configuration defaults support local development (i.e. "go run ...")

	defaultApplicationID = ""
	defaultUpdateIfExist = false

	// Logging
	defaultLoggingLevel = "info"

	// github
	defaulGithubToken      = ""
	defaulGithubOwner      = ""
	defaulGithubRepository = ""
	defaulGithubActor      = ""
	defaulGithubTag        = ""

	defaultTemplateTitle = "defaultTitle.tmpl"
	defaultTemplateBody  = "defaultBody.tmpl"
)

var (
	_ = pflag.String("app.id", defaultApplicationID, "identifier for the application")
	_ = pflag.Bool("update.exist", defaultUpdateIfExist, "id of the release notes update, if it already exists")

	_ = pflag.String("logging.level", defaultLoggingLevel, "log level of application")

	_ = pflag.String("github.token", defaulGithubToken, "github token")
	_ = pflag.String("github.owner", defaulGithubOwner, "github owner")
	_ = pflag.String("github.repository.owner", defaulGithubOwner, "github owner")
	_ = pflag.String("github.repository", defaulGithubRepository, "github repository name")
	_ = pflag.String("github.actor", defaulGithubActor, "github user name")
	_ = pflag.String("github.tag", defaulGithubTag, "github repository tag")

	_ = pflag.String("template.title", defaultTemplateTitle, "default template for release notes' title")
	_ = pflag.String("template.body", defaultTemplateBody, "default template for release notes' body")
)
