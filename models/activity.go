package models

import "time"

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
	TSS                float64   `json:"tss"`               // Calculated TSS
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
		TSS:                hrTSS,
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
