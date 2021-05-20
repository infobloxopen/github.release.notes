package main

import "github.com/spf13/pflag"

const (
	// configuration defaults support local development (i.e. "go run ...")

	defaultConfigFile    = ""
	defaultSecretFile    = ""
	defaultApplicationID = "github.release.notes"

	// Logging
	defaultLoggingLevel = "debug"

	// github
	defaulGithubToken = ""
	defaulGithubOrg   = "infobloxopen"
	defaulGithubRepo  = "github.release.notes"
	defaulGithubUser  = "user"
	defaulGithubGit   = ""
)

var (
	_ = pflag.String("config.file", defaultConfigFile, "directory of the configuration file")
	_ = pflag.String("config.secret.file", defaultSecretFile, "directory of the secrets configuration file")
	_ = pflag.String("app.id", defaultApplicationID, "identifier for the application")

	_ = pflag.String("logging.level", defaultLoggingLevel, "log level of application")

	_ = pflag.String("github.token", defaulGithubToken, "github token")
	_ = pflag.String("github.org", defaulGithubOrg, "github organization")
	_ = pflag.String("github.repo", defaulGithubRepo, "github repository name")
	_ = pflag.String("github.user", defaulGithubUser, "github user name")
	_ = pflag.String("github.tag", defaulGithubUser, "github repository tag")
)
