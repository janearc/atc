package models

// Athlete defines the base structure for an athlete.
type Athlete struct {
	Id         string     `json:"id"`         // Unique identifier for the athlete
	FirstName  string     `json:"first_name"` // Athlete's first name
	LastName   string     `json:"last_name"`  // Athlete's last name
	Age        int        // Athlete's age
	Activities []Activity `json:"activities"`
}

// NewAthlete creates a new athlete with the provided details.
func NewAthlete(id, firstName, lastName string, age int, gender, email, sport, bio string) *Athlete {
	return &Athlete{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

// FullName returns the athlete's full name.
func (a *Athlete) FullName() string {
	return a.FirstName + " " + a.LastName
}
