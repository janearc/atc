package models

import "time"

// activity model, currently just ganked from strava's activity model

// stravaActivity is the strava blob, but we don't need all this stuff so we're just going to create a new activity object that's smaller and more useful for internal use.
type stravaActivity struct {
	Id            int64       `json:"id"`
	ResourceState int         `json:"resource_state"`
	ExternalId    interface{} `json:"external_id"`
	UploadId      interface{} `json:"upload_id"`
	Athlete       struct {
		Id            int64 `json:"id"`
		ResourceState int   `json:"resource_state"`
	} `json:"athlete"`
	Name               string    `json:"name"`
	Distance           int       `json:"distance"`
	MovingTime         int       `json:"moving_time"`
	ElapsedTime        int       `json:"elapsed_time"`
	TotalElevationGain int       `json:"total_elevation_gain"`
	Type               string    `json:"type"`
	SportType          string    `json:"sport_type"`
	StartDate          time.Time `json:"start_date"`
	StartDateLocal     time.Time `json:"start_date_local"`
	Timezone           string    `json:"timezone"`
	UtcOffset          int       `json:"utc_offset"`
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
	AverageSpeed    int           `json:"average_speed"`
	MaxSpeed        int           `json:"max_speed"`
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

// NewActivity creates a new activity with the provided details.
func NewActivity(id int64, resourceState int, externalId, uploadId interface{}, athleteId int64, athleteResourceState int, name string, distance, movingTime, elapsedTime, totalElevationGain int, activityType, sportType string, startDate, startDateLocal time.Time, timezone string, utcOffset, achievementCount, kudosCount, commentCount, athleteCount, photoCount int, mapId string, mapPolyline interface{}, mapResourceState int, trainer, commute, manual, private, flagged bool, gearId string, fromAcceptedTag interface{}, averageSpeed, maxSpeed int, deviceWatts, hasHeartrate bool, prCount, totalPhotoCount int, hasKudoed bool, workoutType interface{}, description interface{}, calories int, segmentEfforts []interface{}) *Activity {
	return &Activity{
		// ...
	}
}

func NewActivityFromStravaActivity { sa stravaActivity } *Activity {
	a := NewActivity(sa.Id, sa.ResourceState, sa.ExternalId, sa.UploadId, sa.Athlete.Id, sa.Athlete.ResourceState, sa.Name, sa.Distance, sa.MovingTime, sa.ElapsedTime, sa.TotalElevationGain, sa.Type, sa.SportType, sa.StartDate, sa.StartDateLocal, sa.Timezone, sa.UtcOffset, sa.AchievementCount, sa.KudosCount, sa.CommentCount, sa.AthleteCount, sa.PhotoCount, sa.Map.Id, sa.Map.Polyline, sa.Map.ResourceState, sa.Trainer, sa.Commute, sa.Manual, sa.Private, sa.Flagged, sa.GearId, sa.FromAcceptedTag, sa.AverageSpeed, sa.MaxSpeed, sa.DeviceWatts, sa.HasHeartrate, sa.PrCount, sa.TotalPhotoCount, sa.HasKudoed, sa.WorkoutType, sa.Description, sa.Calories, sa.SegmentEfforts)

	return a
}