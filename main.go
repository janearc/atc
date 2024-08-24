package main

import (
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

		token := t.GetAccessToken()
		log.Info("Successfully retrieved access token")

		// Fetch the last six weeks of Swim, Bike, and Run activities
		log.Info("Fetching activities...")
		stravaActivities, err := t.FetchActivities(token)
		if err != nil {
			log.WithError(err).Error("Failed to fetch activities")
			http.Error(w, "Failed to fetch activities", http.StatusInternalServerError)
			return
		}

		// Map Strava activities to your model's Activity struct and calculate TSS
		var activities []models.Activity
		for _, sa := range stravaActivities {
			var thresholdHR float64

			// Determine the correct threshold HR based on the activity type
			switch sa.Type {
			case "Run":
				thresholdHR = config.Athlete.Run.ThresholdHR
			case "Ride":
				thresholdHR = config.Athlete.Bike.ThresholdHR
			case "Swim":
				thresholdHR = config.Athlete.Swim.ThresholdHR
			default:
				log.Warnf("Unknown activity type: %s", sa.Type)
				continue // Skip unknown activity types
			}

			activity := models.NewActivity(sa, thresholdHR)
			activities = append(activities, activity)
		}

		// Create the Athlete struct and assign activities
		athlete := models.Athlete{
			Activities: activities,
		}

		// Render the activities in an HTML table
		renderActivitiesTable(w, athlete.Activities)
	})

	// Start the server on the configured port
	log.Infof("Starting server on :%d", config.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil))
}

// renderActivitiesTable generates an HTML table of activities and writes it to the response writer
func renderActivitiesTable(w http.ResponseWriter, activities []models.Activity) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Start the HTML document
	fmt.Fprintf(w, "<html><head><title>Activity Data</title></head><body>")
	fmt.Fprintf(w, "<h1>Your Activities</h1>")
	fmt.Fprintf(w, "<table border='1'><tr><th>Type</th><th>Duration (min)</th><th>TSS</th></tr>")

	// Populate the table with activity data
	for _, activity := range activities {
		durationMinutes := activity.MovingTime / 60
		fmt.Fprintf(w, "<tr><td>%s</td><td>%d</td><td>%.2f</td></tr>", activity.Type, durationMinutes, activity.TSS)
	}

	// End the HTML document
	fmt.Fprintf(w, "</table></body></html>")
}
