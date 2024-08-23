package backend

import (
	"atc/models"
	"atc/transport"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

var log = logrus.New()

// FetchActivities retrieves the last six weeks of activities from Strava.
func FetchActivities(config *transport.Config, token string) ([]models.Activity, error) {
	log.Info("Starting to fetch activities from Strava")

	// Prepend the correct API path to the base URL
	url := fmt.Sprintf("%s/api/v3/athlete/activities", config.Strava.Url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithError(err).Error("Failed to create new request")
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	// Set the time frame for the last six weeks
	sixWeeksAgo := time.Now().AddDate(0, 0, -42).Unix()
	q := req.URL.Query()
	q.Add("after", fmt.Sprintf("%d", sixWeeksAgo))
	req.URL.RawQuery = q.Encode()

	// Log the full request URL
	log.Infof("Making request to: %s", req.URL.String())

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to make request to Strava API")
		return nil, err
	}
	defer resp.Body.Close()

	// Log the response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read response body")
		return nil, err
	}
	log.Infof("Strava API response: %s", string(body))

	// Decode the response into a slice of StravaActivity
	var stravaActivities []models.StravaActivity
	if err := json.Unmarshal(body, &stravaActivities); err != nil {
		log.WithError(err).Error("Failed to decode Strava activities response")
		return nil, err
	}

	log.Infof("Fetched %d activities", len(stravaActivities))

	// Convert StravaActivity to internal Activity type
	activities := make([]models.Activity, len(stravaActivities))
	for i, sa := range stravaActivities {
		activities[i] = models.NewActivity(sa)
	}

	return activities, nil
}
