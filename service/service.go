package service

import (
	"atc/transport"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

// abstracting away the various backend-y type things the app uses

type Service struct {
	Web     WebService
	Backend *transport.Transport
	Config  *transport.Config
	Log     *logrus.Logger
}

type WebService struct {
	// the web server
	Handle http.Handler
}

func NewService(configFileName string, versionFileName string, secretsFileName string) *Service {
	// create a new service

	//
	// Set up Logrus for structured logging
	//

	var log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	//
	// Load configuration yml
	//

	config, err := transport.LoadConfig(configFileName, versionFileName)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	//
	// create the backend transport
	//

	backend, err := transport.NewTransport(config, secretsFileName)
	if err != nil {
		log.Fatalf("Failed to initialize transport: %v", err)
	}

	s := &Service{
		Config:  config,
		Log:     log,
		Backend: backend,
		Web: WebService{
			// NOTE: this creates the http listener
			Handle: instantiateWebService(),
		},
	}

	// Set up the http request handlers ("endpoints")
	s.oauthRedirectHandler()
	s.oauthCallbackHandler()
	s.activitiesHandler()
	s.aboutHandler()

	// All you gotta do now is s.Start()
	return s
}

func (s *Service) Start() {
	// Start the server on the configured port
	s.Log.Infof("Starting server on :%d", s.Config.Server.Port)
	s.Log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Server.Port), nil))
}
