package models_test

import (
	"atc/models"
	"atc/transport"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStravaActivity(t *testing.T) {
	/*
		// StravaActivity represents the detailed activity data returned by the Strava API.
		type StravaActivity struct {
			Id                 int64     `json:"id"`
			Name               string    `json:"name"`
			Distance           float64   `json:"distance"`             // in meters
			MovingTime         int       `json:"moving_time"`          // in seconds
			ElapsedTime        int       `json:"elapsed_time"`         // in seconds
			TotalElevationGain float64   `json:"total_elevation_gain"` // in meters
			Type               string    `json:"type"`
			StartDate          time.Time `json:"start_date"`
			Calories           int       `json:"calories"`
			AverageHeartRate   float64   `json:"average_heartrate"` // in bpm
			MaxHeartRate       float64   `json:"max_heartrate"`     // in bpm
		}

	*/

	startTime := "2024-09-05"
	layout := "2006-01-02"

	// get a time.Time object for our tests
	parsedTime, err := time.Parse(layout, startTime)

	// we're not testing the time package but we need this to move forward
	assert.Nil(t, err)

	var fake = models.NewStravaActivity(
		1234,
		"today i swam",
		float64(1234), // meters
		60*60,         // seconds, one hour
		60*60-30,      // seconds, gotta rest
		float64(150),  // meters

		// need an anti-test for activities that are like "Donut" type
		"Swim",

		parsedTime,
		750,
		115,
		117,
	)

	// just determine whether the constructor shit the bed
	assert.NotNil(t, fake)
}

func TestLoadConfig(t *testing.T) {
	root := os.Getenv("ATC_ROOT")

	configFileName := filepath.Join(root, "config/config.yml")
	versionFileName := filepath.Join(root, "config/version.yml")

	// fire up that constructor and read some yaml
	config, err := transport.LoadConfig(configFileName, versionFileName)

	assert.Nil(t, err)
	assert.NotNil(t, config)
}
