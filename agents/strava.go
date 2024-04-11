package strava

import "time"
import "net/http"

// public classes

// Strava is a parent/abstraction class for strava interaction
type Strava struct {
	BaseUrl      string        // probably https://www.strava.com/api/v3/activities
	ClientId     string        // this is probably an integer but let's call it a string for now
	secret       string        // this is probably a string
	accessToken  string        // this is probably a string
	refreshToken string        // this is probably a string
	athlete      stravaAthlete // this refers to our stravaathlete struct
}

// private classes

// imported from strava api docs
type stravaAthlete struct {
	Id                    int64         `json:"id"`
	Username              string        `json:"username"`
	ResourceState         int           `json:"resource_state"`
	Firstname             string        `json:"firstname"`
	Lastname              string        `json:"lastname"`
	City                  string        `json:"city"`
	State                 string        `json:"state"`
	Country               string        `json:"country"`
	Sex                   string        `json:"sex"`
	Premium               bool          `json:"premium"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	BadgeTypeId           int           `json:"badge_type_id"`
	ProfileMedium         string        `json:"profile_medium"`
	Profile               string        `json:"profile"`
	Friend                interface{}   `json:"friend"`
	Follower              interface{}   `json:"follower"`
	FollowerCount         int           `json:"follower_count"`
	FriendCount           int           `json:"friend_count"`
	MutualFriendCount     int           `json:"mutual_friend_count"`
	AthleteType           int           `json:"athlete_type"`
	DatePreference        string        `json:"date_preference"`
	MeasurementPreference string        `json:"measurement_preference"`
	Clubs                 []interface{} `json:"clubs"`
	Ftp                   interface{}   `json:"ftp"`
	Weight                int           `json:"weight"`
	Bikes                 []struct {
		Id            string `json:"id"`
		Primary       bool   `json:"primary"`
		Name          string `json:"name"`
		ResourceState int    `json:"resource_state"`
		Distance      int    `json:"distance"`
	} `json:"bikes"`
	Shoes []struct {
		Id            string `json:"id"`
		Primary       bool   `json:"primary"`
		Name          string `json:"name"`
		ResourceState int    `json:"resource_state"`
		Distance      int    `json:"distance"`
	} `json:"shoes"`
}

type backend struct {
}

// public methods

// NewStrava creates a new strava object with the provided details.
function NewStrava(baseUrl, clientId, secret, accessToken, refreshToken string) *Strava {
	return &Strava{
		BaseUrl:      baseUrl,
		ClientId:     clientId,
		Secret:       secret,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// GetAthlete returns the athlete object for the current user.
func (s *Strava) GetAthlete() stravaAthlete {
	return s.athlete
}

function IsHealthy() bool {
	// this is a placeholder for now

}

function Connect() {
	// this is a placeholder for now

}

// private methods

function pingStrava() {
	// this is a placeholder for now
}

