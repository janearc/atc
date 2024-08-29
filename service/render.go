package service

import (
	"atc/models"
	"fmt"
	"net/http"
)

// TODO: this could probably be broken into html renders and json renders per github #21

// this file contains helper functions which render html for the web service

// renderActivitiesTableWithCTL generates an HTML table of activities with CTL and IF values
// and writes it back to the http writer
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
