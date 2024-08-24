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
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

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

		// Calculate CTL for Swim, Bike, and Run separately using models.CalculateCTL
		swimCTL := models.CalculateCTL(filterActivitiesByType(activities, "Swim"), 42)
		bikeCTL := models.CalculateCTL(filterActivitiesByType(activities, "Ride"), 42)
		runCTL := models.CalculateCTL(filterActivitiesByType(activities, "Run"), 42)

		// Render the activities in an HTML table and display CTL
		renderActivitiesTableWithCTL(w, activities, swimCTL, bikeCTL, runCTL)
	})

	// Start the server on the configured port
	log.Infof("Starting server on :%d", config.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil))
}

// filterActivitiesByType filters the activities by type (Swim, Bike, Run)
func filterActivitiesByType(activities []models.Activity, activityType string) []models.Activity {
	var filtered []models.Activity
	for _, activity := range activities {
		if activity.Type == activityType {
			filtered = append(filtered, activity)
		}
	}
	return filtered
}

func renderActivitiesTableWithCTL(w http.ResponseWriter, activities []models.Activity, swimCTL, bikeCTL, runCTL float64) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Start the HTML document
	fmt.Fprintf(w, "<html><head><title>Activity Data</title></head><body>")
	fmt.Fprintf(w, "<h1>Your Activities</h1>")
	fmt.Fprintf(w, "<table border='1'><tr><th>Date</th><th>Type</th><th>Duration (min)</th><th>TSS</th></tr>")

	// Populate the table with activity data
	for _, activity := range activities {
		durationMinutes := activity.MovingTime / 60
		activityDate := activity.StartDate.Format("2006-01-02")
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%d</td><td>%d</td></tr>", activityDate, activity.Type, durationMinutes, activity.TSS)
	}

	// Display the CTL for each sport
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<h2>Chronic Training Load (CTL)</h2>")
	fmt.Fprintf(w, "<p>Swim CTL: %.2f</p>", swimCTL)
	fmt.Fprintf(w, "<p>Bike CTL: %.2f</p>", bikeCTL)
	fmt.Fprintf(w, "<p>Run CTL: %.2f</p>", runCTL) // Correctly passing runCTL

	// End the HTML document
	fmt.Fprintf(w, "</body></html>")
}
