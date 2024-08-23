package main

import (
	"atc/backend"
	"atc/models"
	"atc/transport"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.New()

func main() {
	// Set up Logrus for structured logging
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Load the configuration
	config, err := transport.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the transport using NewTransport
	t, err := transport.NewTransport(config)
	if err != nil {
		log.Fatalf("Failed to initialize transport: %v", err)
	}

	// Serve static files from the "web" directory
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))

	// Handle the root path to serve the homepage or redirect to Strava's OAuth page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			log.Info("Serving index.html")
			http.ServeFile(w, r, "./web/index.html")
		} else {
			http.NotFound(w, r)
		}
	})

	// Handle the callback from Strava
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Received callback from Strava")

		code := r.URL.Query().Get("code")
		if code == "" {
			log.Error("Code not found in callback")
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		log.Info("Exchanging code for token...")
		if err := t.ExchangeCodeForToken(code); err != nil {
			log.WithError(err).Error("Failed to exchange code for token")
			http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
			return
		}

		token := t.GetAccessToken()
		log.Info("Successfully retrieved access token")

		// Fetch the last six weeks of activities
		log.Info("Fetching activities...")
		activities, err := backend.FetchActivities(config, token)
		if err != nil {
			log.WithError(err).Error("Failed to fetch activities")
			http.Error(w, "Failed to fetch activities", http.StatusInternalServerError)
			return
		}

		log.Info("Successfully fetched activities")

		// Create the Athlete struct and assign activities
		athlete := models.Athlete{
			Activities: activities,
		}

		// Display the athlete's activities in the response
		log.WithField("activities", len(athlete.Activities)).Info("Displaying athlete activities")
		fmt.Fprintf(w, "Athlete Activities: %+v\n", athlete)
	})

	// Start the server on the configured port
	log.Infof("Starting server on :%d", config.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil))
}
