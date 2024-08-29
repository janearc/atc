package main

import (
	"atc/service"
)

func main() {

	// instantiate the service
	s := service.NewService()

	// build the service object. this will pop Fatal if it fails so we don't have to worry about that here.
	if s != nil {
		s.Log.Info("Service object instantiated")
	}

	// Start the server on the configured port
	s.Log.Infof("Starting service from wrapper on port %d", s.Config.Server.Port)

	// Start the server
	s.Start()

	// This is the end of the main function. The service will run until it is stopped.
	s.Log.Fatal("Service stopped unexpectedly.")

}
