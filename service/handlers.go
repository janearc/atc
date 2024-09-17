package service

import (
	"atc/models"
	"fmt"
	"net/http"
	"time"
)

// this is a redirect to strava for oauth and to let the user know whats up
func (s *Service) oauthRedirectHandler() {
	// Redirect to Strava for OAuth
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		s.Log.Infof("[%s]: Redirecting to Strava for OAuth -> [%s]...", r.URL.Path, s.Backend.GetAuthURL())
		http.Redirect(w, r, s.Backend.GetAuthURL(), http.StatusFound)
	})

	return
}

// this exists to respond to oauth callbacks and isn't interactive
func (s *Service) oauthCallbackHandler() {
	// Handle the callback from Strava and store the token in a cookie
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		s.Log.Infof("[%s]: Received callback from Strava with URL: %s", r.URL.Path, r.URL.String())

		s.Backend.AuthGood()

		token := r.URL.Query().Get("code")
		if token == "" {
			s.Log.Warn("No token found in callback")
			http.Error(w, "No token found in callback", http.StatusBadRequest)
		} else {
			s.Log.Infof("Received token: %s", token)
			s.Backend.SetAccessToken(token)
		}

		// TODO: where is this found?
		// s.Backend.SetRefreshToken()

		// Cookie the access token
		http.SetCookie(w, &http.Cookie{
			Name:       "strava_token",
			Value:      s.Backend.GetAccessToken(),
			Path:       "/",
			Domain:     "",
			Expires:    time.Now().Add(24 * time.Hour),
			RawExpires: "",
			MaxAge:     0,
			Secure:     true,
			HttpOnly:   true,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		})

		// Cookie the refresh token
		http.SetCookie(w, &http.Cookie{
			Name:     "strava_refresh_token",
			Value:    s.Backend.GetRefreshToken(),
			Path:     "/",
			Expires:  time.Now().Add(30 * 24 * time.Hour), // Refresh token typically lasts longer
			HttpOnly: true,
			Secure:   true,
		})

		// Also need to set this in the transport object

		http.Redirect(w, r, "/activities", http.StatusFound)
	})

	return
}

// /activities is the endpoint that displays activities and data
func (s *Service) activitiesHandler() {
	// behind the scenes this is wrapping transport.FetchActivities

	// Handle requests to fetch activities and display CTL
	http.HandleFunc("/activities", func(w http.ResponseWriter, r *http.Request) {
		if s.Backend.Authenticated() == false {
			s.Log.Info("not authenticated, passing to /auth")
			http.Redirect(w, r, "/auth", http.StatusFound)
		} else {
			s.Log.Info("authenticated, attempting to fetch activities")
		}

		// Fetch the last six weeks of Swim, Bike, and Run activities
		s.Log.Info("Fetching activities...")
		stravaActivities, err := s.Backend.FetchActivities()
		if err != nil {
			s.Log.WithError(err).Error("Failed to fetch activities")
			// http.Error(w, "Failed to fetch activities", http.StatusInternalServerError)
			return
		}

		s.Log.Infof("Fetched %d activities", len(stravaActivities))

		if len(stravaActivities) == 0 {
			// send to both syslog and the browser to let them know what's happened
			s.Log.Warn("No activities found")
			_, perr := fmt.Fprintf(w, "No activities found")
			if perr != nil {
				s.Log.WithError(perr).Error("error writing to socket")
			}
			return
		}

		// Map Strava activities to native Activity struct and calculate TSS
		var activities []models.Activity
		for _, sa := range stravaActivities {
			var thresholdHR float64

			// Determine the correct threshold HR based on the activity type
			switch sa.Type {
			case "Run":
				thresholdHR = s.Config.Athlete.Run.ThresholdHR
			case "Ride":
				thresholdHR = s.Config.Athlete.Bike.ThresholdHR
			case "Swim":
				thresholdHR = s.Config.Athlete.Swim.ThresholdHR
			default:
				s.Log.Warnf("Unexpected/unknown activity type: %s", sa.Type)
				continue // Skip unwanted activity types
			}

			// this constructs our new native activity, which calculates
			//   tss, trimps, and hrtss
			// in the constructor (models/activity) so we don't have to.
			activity := models.NewActivity(sa, thresholdHR)
			activities = append(activities, activity)
		}

		s.Log.Infof("Mapped to %d activities", len(activities))

		// Calculate CTL for Swim, Bike, and Run separately using models.CalculateCTL
		swimCTL := models.CalculateCTL(models.FilterActivitiesByType(activities, "Swim"), 42)
		bikeCTL := models.CalculateCTL(models.FilterActivitiesByType(activities, "Ride"), 42)
		runCTL := models.CalculateCTL(models.FilterActivitiesByType(activities, "Run"), 42)

		// ask renderer to display the activities in a table with CTL and IF
		renderActivitiesTableWithCTL(w, activities, swimCTL, bikeCTL, runCTL)
	})

	return
}

// returns information about the service
func (s *Service) aboutHandler() {
	// handle the "about" request
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Activity Dashboard</title>
			<style>
				/* TODO: */
			</style>
		</head>
		<body>
			<div class="container">
				<h3>ATC</h3>
				<p>ATC is a web application that helps athletes track their performance and progress in swimming, biking, and running.</p>
				<p>author: Jane Arc</p>
				<p>Build Version: %s</p>
				<p>Build Date: %s</p>
				<p>source: <a href="http://github.com/janearc/atc">http://github.com/janearc/atc</a></p>
			</div>
		</body>
		</html>
		`, s.Config.Build.Build, s.Config.Build.BuildDate)

		w.Write([]byte(html))
	})

	return
}
