package models

// Athlete defines the base structure for an athlete.
type Athlete struct {
	ID        string // Unique identifier for the athlete
	FirstName string // Athlete's first name
	LastName  string // Athlete's last name
	Age       int    // Athlete's age
}

// NewAthlete creates a new athlete with the provided details.
func NewAthlete(id, firstName, lastName string, age int, gender, email, sport, bio string) *Athlete {
	return &Athlete{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

// FullName returns the athlete's full name.
func (a *Athlete) FullName() string {
	return a.FirstName + " " + a.LastName
}
