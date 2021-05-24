package main

import (
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	gh "github.com/infobloxopen/github.release.notes/pkg/githubclient"
)

func main() {
	doneC := make(chan error)
	logger := NewLogger()

	logrus.Debugf("update.exist: %v", viper.GetString("update.exist"))
	logrus.Debugf("github.repository: %v", viper.GetString("github.repository"))
	logrus.Debugf("github.actor: %v", viper.GetString("github.actor"))
	logrus.Debugf("github.owner: %v", viper.GetString("github.owner"))
	logrus.Debugf("github.tag: %v", viper.GetString("github.tag"))

	logrus.Infof("github.tag2: %v", os.Getenv("GITHUB_REPOSITORY_OWNER"))
	logrus.Infof("github.actor2: %v", os.Getenv("GITHUB_ACTOR"))

	go func() { doneC <- ServeExternal() }()

	if err := <-doneC; err != nil {
		logger.Fatal(err)
	}
}

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

// ServeExternal builds and runs the server that listens on ServerAddress and GatewayAddress
func ServeExternal() error {
	// check if repo variable contains owner and repository name
	slice := strings.Split(viper.GetString("github.repository"), "/")
	if len(slice) > 1 {
		viper.Set("github.repository", slice[1])
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
			Repository: viper.GetString("github.repository"),
		},
	)

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
