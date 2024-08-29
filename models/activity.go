package models

import (
	"math"
	"time"
)

// Activity is a simplified version of StravaActivity for internal use.
type Activity struct {
	Id                 int64     `json:"id"`
	Name               string    `json:"name"`
	Distance           float64   `json:"distance"`             // in meters
	MovingTime         int       `json:"moving_time"`          // in seconds
	ElapsedTime        int       `json:"elapsed_time"`         // in seconds
	TotalElevationGain float64   `json:"total_elevation_gain"` // in meters
	Type               string    `json:"type"`
	StartDate          time.Time `json:"start_date"`
	Calories           int       `json:"calories"`
	TSS                int       `json:"tss"`               // Rounded TSS
	IntensityFactor    float64   `json:"intensity_factor"`  // IF
	AverageHeartRate   float64   `json:"average_heartrate"` // in bpm
	MaxHeartRate       float64   `json:"max_heartrate"`     // in bpm
}

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

// NewActivity creates a new Activity from a StravaActivity and calculates TSS.
func NewActivity(sa StravaActivity, thresholdHR float64) Activity {

	// calculate our normalized values for fitness
	hrTSS := calculateHrTSS(sa.MovingTime, sa.AverageHeartRate, thresholdHR)
	intensityFactor := calculateIntensityFactor(sa.AverageHeartRate, thresholdHR)

	return Activity{
		// these values are ganked from the strava object
		Id:                 sa.Id,
		Name:               sa.Name,
		Distance:           sa.Distance,
		MovingTime:         sa.MovingTime,
		ElapsedTime:        sa.ElapsedTime,
		TotalElevationGain: sa.TotalElevationGain,
		Type:               sa.Type,
		StartDate:          sa.StartDate,
		Calories:           sa.Calories,
		MaxHeartRate:       sa.MaxHeartRate,
		AverageHeartRate:   sa.AverageHeartRate,

		// these are our values
		TSS:             hrTSS,
		IntensityFactor: intensityFactor,
	}
}

func calculateIntensityFactor(averageHeartRate float64, thresholdHR float64) float64 {
	return averageHeartRate / thresholdHR
}

// calculateHrTSS calculates TSS based on heart rate data.
func calculateHrTSS(movingTime int, averageHeartRate float64, thresholdHR float64) int {
	// Convert moving time from seconds to hours
	durationHours := float64(movingTime) / 3600.0

	// Calculate IF
	IF := calculateIntensityFactor(averageHeartRate, thresholdHR)

	// Calculate hrTSS (this is suspicious but also the number seems accurate so)
	hrTSS := durationHours * IF * IF * 100

	return int(math.Round(hrTSS))
}

// calculateCTL calculates the Chronic Training Load (CTL) based on TSS values.
func CalculateCTL(activities []Activity, days int) float64 {
	decayFactor := 2.0 / float64(days+1)
	var ctl float64

	// Issue #11
	// https://github.com/janearc/atc/issues/11
	for i, activity := range activities {
		if i == 0 {
			ctl = float64(activity.TSS)
		} else {
			ctl = ctl*(1-decayFactor) + float64(activity.TSS)*decayFactor
		}
	}
	return ctl
}

// FilterActivitiesByType filters the activities by supplied type, e.g., Swim, Ride, Run
func FilterActivitiesByType(activities []Activity, activityType string) []Activity {
	var filtered []Activity
	for _, activity := range activities {
		if activity.Type == activityType {
			filtered = append(filtered, activity)
		}
	}
	return filtered
}
