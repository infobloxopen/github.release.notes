package main

import (
	"log"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	gh "github.com/infobloxopen/github.release.notes/pkg/githubclient"
)

func main() {
	logger := NewLogger()

	logrus.Debugf("update.exist: %v", viper.GetString("update.exist"))
	logrus.Debugf("github.repository: %v", viper.GetString("github.repository"))
	logrus.Debugf("github.actor: %v", viper.GetString("github.actor"))
	logrus.Debugf("github.owner: %v", viper.GetString("github.owner"))
	logrus.Debugf("github.repository.owner: %v", viper.GetString("github.repository.owner"))
	logrus.Debugf("github.tag: %v", viper.GetString("github.tag"))
	logrus.Debugf("template.title: %v", viper.GetString("template.title"))
	logrus.Debugf("template.body: %v", viper.GetString("template.body"))

	err := publishNotes()

	if err != nil {
		logger.Fatal(err)
	}
}

// NewLogger sets log level for standart logger
func NewLogger() *logrus.Logger {
	logger := logrus.StandardLogger()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)

	// Set the log level on the default logger based on command line flag
	if level, err := logrus.ParseLevel(viper.GetString("logging.level")); err != nil {
		logger.Errorf("Invalid %q provided for log level", viper.GetString("logging.level"))
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}

	return logger
}

// publishNotes builds and runs the server that listens on ServerAddress and GatewayAddress
func publishNotes() error {
	repo := viper.GetString("github.repository")
	// check if repo variable contains owner and repository name
	repoSlice := strings.Split(repo, "/")
	if len(repoSlice) > 1 {
		repo = repoSlice[1]
		logrus.Infof("Repository variable is overridden because of the variable contains owner name")
	}

	owner := viper.GetString("github.owner")
	if owner == "" {
		owner = viper.GetString("github.repository.owner")
	}

	ghClient := gh.NewGithubClient(
		gh.GithubClientOptions{
			Token:      viper.GetString("github.token"),
			Owner:      owner,
			Repository: repo,
		},
	)

	logrus.Debugf("github.owner: %v", viper.GetString("github.owner"))
	logrus.Debugf("github.owner1: %v", viper.GetString("github.repository.owner"))
	logrus.Debugf("github.owner2: %v", owner)

	rndList, err := ghClient.GetReleaseNotesData(viper.GetString("github.tag"))
	if err != nil {
		return err
	}
	ghClient.PublishReleaseNotes(rndList)

	return nil
}

func init() {
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalln(err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(viper.GetString("config.source"))

	log.Printf("Serving from default values, environment variables, and/or flags")
}
