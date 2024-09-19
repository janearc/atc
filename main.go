package main

import (
	"atc/service"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func main() {
	// You must have ATC_ROOT set in your environment. Sorry about that.
	root := os.Getenv("ATC_ROOT")

	defaultConfigFileName := filepath.Join(root, "config/config.yml")
	defaultVersionFileName := filepath.Join(root, "config/version.yml")
	defaultSecretsFileName := filepath.Join(root, "config/secrets.yml")

	// instantiate the service
	s := service.NewService(defaultConfigFileName, defaultVersionFileName, defaultSecretsFileName)

	// build the service object. this will pop Fatal if it fails so we don't have to worry about that here.
	if s != nil {
		s.Log.Info("Service object instantiated")

		// Start the server on the configured port
		s.Log.Infof("Starting service from wrapper on port %d", s.Config.Server.Port)

		// Start the server
		s.Start()

		// The service will run until it is stopped.
		s.Log.Fatal("Service stopped unexpectedly.")
	} else {
		logrus.Fatalf("Failed to instantiate service object")
	}

	// This should be unreachable code
	logrus.Fatalf("Unexpected failure in main.go")
}
