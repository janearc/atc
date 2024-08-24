package main

import (
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
	fs := http.FileServer(http.Dir("/app/web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))

	// Serve index.html on the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			log.Info("Serving index.html")
			http.ServeFile(w, r, "/app/web/index.html")
		} else {
			http.NotFound(w, r)
		}
	})

	// Handle the OAuth redirect to Strava
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Redirecting to Strava's OAuth page")
		http.Redirect(w, r, t.GetAuthURL(), http.StatusFound)
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

		log.Info("Successfully retrieved access token")

		// Make a request to OpenAI
		prompt := "Tell me about the importance of leg day."
		response, err := t.OpenAIRequest(prompt)
		if err != nil {
			log.WithError(err).Error("Failed to get response from OpenAI")
			http.Error(w, "Failed to get response from OpenAI", http.StatusInternalServerError)
			return
		}

		log.Info("OpenAI response: ", response)
		fmt.Fprintf(w, "OpenAI Response: %s\n", response)
	})

	// Start the server on the configured port
	log.Infof("Starting server on :%d", config.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil))
}
