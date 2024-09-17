package models

// Athlete defines the base structure for an athlete.
type Athlete struct {
	Id         string     `json:"id"`         // Unique identifier for the athlete
	FirstName  string     `json:"first_name"` // Athlete's first name
	LastName   string     `json:"last_name"`  // Athlete's last name
	Age        int        // Athlete's age
	Sex        string     // not sure if this is strictly required but strava's gonna tell us anyways
	Activities []Activity `json:"activities"`
	Thresholds Thresholds
}

type Thresholds struct {
	Run struct {
		ThresholdHR float64 `yaml:"threshold_hr"`
	} `yaml:"run"`
	Swim struct {
		ThresholdHR float64 `yaml:"threshold_hr"`
	} `yaml:"swim"`
	Bike struct {
		ThresholdHR float64 `yaml:"threshold_hr"`
	} `yaml:"bike"`
}

// NewAthlete creates a new athlete with the provided details.
func NewAthlete(id, firstName, lastName string, sex string, thresholds *Thresholds) *Athlete {
	// the threshold stuff is a little bit of a hack because strava
	// doesn't seem to want to give us this data. so we're going to
	// hack this together from service config

	runThreshold := thresholds.Run.ThresholdHR
	swimThreshold := thresholds.Swim.ThresholdHR
	bikeThreshold := thresholds.Bike.ThresholdHR

	t := Thresholds{}

	t.Run.ThresholdHR = runThreshold
	t.Swim.ThresholdHR = swimThreshold
	t.Bike.ThresholdHR = bikeThreshold

	return &Athlete{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Sex:       sex,

		Thresholds: t,
	}
}

// FullName returns the athlete's full name.
func (a *Athlete) FullName() string {
	return a.FirstName + " " + a.LastName
}

func (a *Athlete) GetRunThreshold() float64 {
	return a.Thresholds.Run.ThresholdHR
}

func (a *Athlete) GetSwimThreshold() float64 {
	return a.Thresholds.Swim.ThresholdHR
}

func (a *Athlete) GetBikeThreshold() float64 {
	return a.Thresholds.Bike.ThresholdHR
}
