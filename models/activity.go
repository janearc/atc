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
	// Calculate hrTSS based on heart rate data
	hrTSS := calculateHrTSS(sa.MovingTime, sa.AverageHeartRate, thresholdHR)

	return Activity{
		Id:                 sa.Id,
		Name:               sa.Name,
		Distance:           sa.Distance,
		MovingTime:         sa.MovingTime,
		ElapsedTime:        sa.ElapsedTime,
		TotalElevationGain: sa.TotalElevationGain,
		Type:               sa.Type,
		StartDate:          sa.StartDate,
		Calories:           sa.Calories,
		TSS:                int(math.Round(hrTSS)), // Rounded TSS
		AverageHeartRate:   sa.AverageHeartRate,
		MaxHeartRate:       sa.MaxHeartRate,
	}
}

// calculateHrTSS calculates TSS based on heart rate data.
func calculateHrTSS(movingTime int, averageHeartRate, thresholdHR float64) float64 {
	// Convert moving time from seconds to hours
	durationHours := float64(movingTime) / 3600.0

	// Calculate the Intensity Factor (IF)
	IF := averageHeartRate / thresholdHR

	// Calculate hrTSS
	hrTSS := durationHours * IF * IF * 100

	return hrTSS
}

// calculateCTL calculates the Chronic Training Load (CTL) based on TSS values.
func CalculateCTL(activities []Activity, days int) float64 {
	decayFactor := 2.0 / float64(days+1)
	var ctl float64
	for i, activity := range activities {
		if i == 0 {
			ctl = float64(activity.TSS)
		} else {
			ctl = ctl*(1-decayFactor) + float64(activity.TSS)*decayFactor
		}
	}
	return ctl
}
