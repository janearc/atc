package models

import "time"

// activity model, currently just ganked from strava's activity model

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
}

// StravaActivity represents the detailed activity data returned by the Strava API.
type StravaActivity struct {
	Id            int64       `json:"id"`
	ResourceState int         `json:"resource_state"`
	ExternalId    interface{} `json:"external_id"`
	UploadId      interface{} `json:"upload_id"`
	Athlete       struct {
		Id            int64 `json:"id"`
		ResourceState int   `json:"resource_state"`
	} `json:"athlete"`
	Name               string    `json:"name"`
	Distance           float64   `json:"distance"` // Changed to float64
	MovingTime         int       `json:"moving_time"`
	ElapsedTime        int       `json:"elapsed_time"`
	TotalElevationGain float64   `json:"total_elevation_gain"` // Changed to float64
	Type               string    `json:"type"`
	SportType          string    `json:"sport_type"`
	StartDate          time.Time `json:"start_date"`
	StartDateLocal     time.Time `json:"start_date_local"`
	Timezone           string    `json:"timezone"`
	UtcOffset          float64   `json:"utc_offset"`
	AchievementCount   int       `json:"achievement_count"`
	KudosCount         int       `json:"kudos_count"`
	CommentCount       int       `json:"comment_count"`
	AthleteCount       int       `json:"athlete_count"`
	PhotoCount         int       `json:"photo_count"`
	Map                struct {
		Id            string      `json:"id"`
		Polyline      interface{} `json:"polyline"`
		ResourceState int         `json:"resource_state"`
	} `json:"map"`
	Trainer         bool          `json:"trainer"`
	Commute         bool          `json:"commute"`
	Manual          bool          `json:"manual"`
	Private         bool          `json:"private"`
	Flagged         bool          `json:"flagged"`
	GearId          string        `json:"gear_id"`
	FromAcceptedTag interface{}   `json:"from_accepted_tag"`
	AverageSpeed    float64       `json:"average_speed"` // Changed to float64
	MaxSpeed        float64       `json:"max_speed"`     // Changed to float64
	DeviceWatts     bool          `json:"device_watts"`
	HasHeartrate    bool          `json:"has_heartrate"`
	PrCount         int           `json:"pr_count"`
	TotalPhotoCount int           `json:"total_photo_count"`
	HasKudoed       bool          `json:"has_kudoed"`
	WorkoutType     interface{}   `json:"workout_type"`
	Description     interface{}   `json:"description"`
	Calories        int           `json:"calories"`
	SegmentEfforts  []interface{} `json:"segment_efforts"`
}

// public functions
// NewActivity creates a new Activity from a StravaActivity.
func NewActivity(sa StravaActivity) Activity {
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
	}
}
