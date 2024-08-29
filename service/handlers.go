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
		s.Log.Info("Redirecting to Strava for OAuth...")
		http.Redirect(w, r, t.GetOAuthURL(), http.StatusFound)
	})
	return
}

// this exists to respond to oauth callbacks and isn't interactive
func (s *Service) oauthCallbackHandler() {
	// Handle the callback from Strava and store the token in a cookie
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		s.Log.Infof("Received callback from Strava with URL: %s", r.URL.String())

		code := r.URL.Query().Get("code")
		if code == "" {
			s.Log.Error("Code not found in callback")
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		s.Log.Info("Exchanging code for token...")
		if err := t.ExchangeCodeForToken(code); err != nil {
			s.Log.WithError(err).Error("Failed to exchange code for token")
			http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
			return
		}

		token := s.Backend.GetAccessToken()
		refreshToken := s.Backend.GetRefreshToken()

		s.Log.Info("Successfully retrieved access token")

		// Store the access token and refresh token in cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "strava_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "strava_refresh_token",
			Value:    refreshToken,
			Path:     "/",
			Expires:  time.Now().Add(30 * 24 * time.Hour), // Refresh token typically lasts longer
			HttpOnly: true,
			Secure:   true,
		})

		// Redirect to the home page after setting the cookies
		http.Redirect(w, r, "/", http.StatusFound)
	})
	return
}

// /activities is the endpoint that displays activities and data
func (s *Service) activitiesHandler() {
	// Handle requests to fetch activities and display CTL
	http.HandleFunc("/activities", func(w http.ResponseWriter, r *http.Request) {
		// Check for the OAuth token in cookies
		cookie, err := r.Cookie("strava_token")
		if err != nil {
			if err == http.ErrNoCookie {
				s.Log.Warn("No strava_token cookie found, redirecting to OAuth")
				http.Redirect(w, r, "/auth", http.StatusFound)
				return
			}
			s.Log.WithError(err).Error("Failed to retrieve cookie")
			http.Error(w, "Failed to retrieve cookie", http.StatusInternalServerError)
			return
		}

		token := cookie.Value

		// Check if the token is expired and refresh it if necessary
		if t.IsTokenExpired() {
			s.Log.Info("Token expired, attempting to refresh...")
			refreshCookie, err := r.Cookie("strava_refresh_token")
			if err != nil {
				s.Log.WithError(err).Error("Failed to retrieve refresh token cookie, redirecting to OAuth")
				http.Redirect(w, r, "/auth", http.StatusFound)
				return
			}

			refreshToken := refreshCookie.Value
			newAccessToken, err := t.RefreshAccessToken(refreshToken)
			if err != nil {
				s.Log.WithError(err).Error("Failed to refresh access token, redirecting to OAuth")
				http.Redirect(w, r, "/auth", http.StatusFound)
				return
			}

			s.Log.Info("Token refreshed successfully")

			// Update the access token cookie with the new token
			http.SetCookie(w, &http.Cookie{
				Name:     "strava_token",
				Value:    newAccessToken,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
				Secure:   true,
			})

			token = newAccessToken
		} else {
			s.Log.Info("Token is still valid, proceeding with fetching activities")
		}

		// Fetch the last six weeks of Swim, Bike, and Run activities
		s.Log.Info("Fetching activities...")
		stravaActivities, err := t.FetchActivities(token)
		if err != nil {
			s.Log.WithError(err).Error("Failed to fetch activities")
			http.Error(w, "Failed to fetch activities", http.StatusInternalServerError)
			return
		}

		s.Log.Infof("Fetched %d activities", len(stravaActivities))

		if len(stravaActivities) == 0 {
			s.Log.Warn("No activities found")
			fmt.Fprintf(w, "No activities found")
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
				s.Log.Warnf("Unknown activity type: %s", sa.Type)
				continue // Skip unknown activity types
			}

			activity := models.NewActivity(sa, thresholdHR)
			activities = append(activities, activity)
		}

		s.Log.Infof("Mapped to %d activities", len(activities))

		// Calculate CTL for Swim, Bike, and Run separately using models.CalculateCTL
		swimCTL := models.CalculateCTL(filterActivitiesByType(activities, "Swim"), 42)
		bikeCTL := models.CalculateCTL(filterActivitiesByType(activities, "Ride"), 42)
		runCTL := models.CalculateCTL(filterActivitiesByType(activities, "Run"), 42)

		// Render the activities in an HTML table and display CTL
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
