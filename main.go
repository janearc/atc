package main

import (
	"atc/models"
	"atc/transport"
	"fmt"
	"net/http"
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

// XXX: move to models/activity.go
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

// XXX: move to service/handlers.go
// renderActivitiesTableWithCTL generates an HTML table of activities and writes it to the response writer
// Also displays the CTL for Swim, Bike, and Run
func renderActivitiesTableWithCTL(w http.ResponseWriter, activities []models.Activity, swimCTL, bikeCTL, runCTL float64) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Start the HTML document
	fmt.Fprintf(w, "<html><head><title>Activity Data</title></head><body>")
	fmt.Fprintf(w, "<h1>Activities (42 days)</h1>")
	fmt.Fprintf(w, "<table border='1'><tr><th>Date</th><th>Type</th><th>Duration (min)</th><th>TSS</th><th>IF</th></tr>")

	// Populate the table with activity data
	for _, activity := range activities {
		durationMinutes := activity.MovingTime / 60
		activityDate := activity.StartDate.Format("2006-01-02")
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%.2f</td></tr>",
			activityDate,
			activity.Type,
			durationMinutes,
			activity.TSS,
			activity.IntensityFactor)
	}

	// Display the CTL for each sport
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<h2>Chronic Training Load (CTL)</h2>")
	fmt.Fprintf(w, "<p>Swim CTL: %.2f</p>", swimCTL)
	fmt.Fprintf(w, "<p>Bike CTL: %.2f</p>", bikeCTL)
	fmt.Fprintf(w, "<p>Run CTL: %.2f</p>", runCTL)

	// End the HTML document
	fmt.Fprintf(w, "</body></html>")
}
