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
	doneC := make(chan error)
	logger := NewLogger()

	logrus.Debugf("update.exist: %v", viper.GetString("update.exist"))
	logrus.Debugf("github.repo: %v", viper.GetString("github.repo"))
	logrus.Debugf("github.user: %v", viper.GetString("github.user"))
	logrus.Debugf("github.org: %v", viper.GetString("github.org"))
	logrus.Debugf("github.tag: %v", viper.GetString("github.tag"))

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
	ghClient := gh.NewGithubClient(viper.GetString("github.token"), viper.GetString("github.org"), viper.GetString("github.repo"))

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
