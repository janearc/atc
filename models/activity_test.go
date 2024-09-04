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

func TestNewActivity(t *testing.T) {
	// func NewActivity(sa StravaActivity, thresholdHR float64) Activity {
	startTime := "2024-09-05"
	layout := "2006-01-02"

	// get a time.Time object for our tests
	parsedTime, err := time.Parse(layout, startTime)

	// we're not testing the time package but we need this to move forward
	assert.Nil(t, err)

	// we just need this strava activity to pass to the constructor.
	var fake = models.NewStravaActivity(
		// Id
		1234,

		// Name
		"today i swam",

		// Distance
		float64(1234), // meters

		// MovingTime
		60*60-30, // seconds, one hour

		// ElapsedTime
		60*60, // seconds, gotta rest

		// Elev Gain
		float64(150), // meters

		// need an anti-test for activities that are like "Donut" type
		// Type
		"Swim",

		// StartDate
		parsedTime,

		// Calories
		750,

		// Avg HR
		115,

		// Max HR
		117,
	)

	// just determine whether the constructor shit the bed
	assert.NotNil(t, fake)

	activity := models.NewActivity(fake, 145)

	assert.NotNil(t, activity)

	// walk the fields to make sure it looks healthy
	assert.Equalf(t, fake.Id, activity.Id, "id matches")
	assert.Equalf(t, fake.Name, activity.Name, "name matches")
	assert.Equalf(t, fake.Distance, activity.Distance, "distance matches")
	assert.Equalf(t, fake.MovingTime, activity.MovingTime, "moving time matches")
	assert.Equalf(t, fake.ElapsedTime, activity.ElapsedTime, "elapsed time matches")
	assert.Equalf(t, fake.Distance, activity.Distance, "distance matches")
	assert.Equalf(t, fake.Type, activity.Type, "type matches")
	assert.Equalf(t, fake.StartDate, activity.StartDate, "start date matches")
	assert.Equalf(t, fake.Calories, activity.Calories, "calories match")
	assert.Equalf(t, fake.AverageHeartRate, activity.AverageHeartRate, "avg hr matches")
	assert.Equalf(t, fake.MaxHeartRate, activity.MaxHeartRate, "max hr matches")

	// these are calculated fields
	assert.NotNil(t, activity.Trimps)
	assert.NotNil(t, activity.TSS)
	assert.NotNil(t, activity.IntensityFactor)

	// TODO: verify those calculated fields are correct
}

func TestNewStravaActivity(t *testing.T) {

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
		60*60-30,      // seconds, moving time
		60*60,         // seconds, elapsed time
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

	assert.Equalf(t, fake.Id, int64(1234), "Id is correct")
	assert.Equalf(t, fake.Name, "today i swam", "Name is correct")
	assert.Equalf(t, fake.Distance, float64(1234), "distance matches")
	assert.Equalf(t, fake.MovingTime, 60*60-30, "movingtime is correct")
	assert.Equalf(t, fake.ElapsedTime, 60*60, "elapsed time is correct")
	assert.Equalf(t, fake.TotalElevationGain, float64(150), "total elevation gain is correct")
	assert.Equalf(t, fake.Type, "Swim", "type is correct")
	assert.Equalf(t, fake.StartDate, parsedTime, "start date is correct")
	assert.Equalf(t, fake.Calories, 750, "calories is correct")
	assert.Equalf(t, fake.AverageHeartRate, float64(115), "avg hr is correct")
	assert.Equalf(t, fake.MaxHeartRate, float64(117), "max hr is correct")
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
