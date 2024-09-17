package planning

import (
	"atc/models"
	"errors"
	"math"
)

// deltaTSS calculates the necessary change in TSS to reach a target CTL and volume.
func deltaTSS(athlete models.Athlete, targetVolume float64, targetCTL float64, activityType string) (float64, error) {
	filteredActivities := models.FilterActivitiesByType(athlete.Activities, activityType)

	filteredCTL := models.CalculateCTL(filteredActivities, 42)
	filteredDuration := models.CalculateDurationHrs(filteredActivities)

	// Estimate the duration in hours (assuming linear relation with volume)
	duration := targetVolume / filteredCTL * filteredDuration

	if activityType == "Run" {
		// Calculate the new Intensity Factor based on target CTL
		newIntensityFactor := math.Sqrt((targetCTL * athlete.GetRunThreshold()) / (duration * 100))

		// Calculate the new TSS based on the new Intensity Factor
		newTSS := duration * newIntensityFactor * newIntensityFactor * 100

		// Calculate the delta TSS
		deltaTSS := newTSS - filteredCTL

		return deltaTSS, nil
	} else if activityType == "Bike" {
		// Calculate the new Intensity Factor based on target CTL
		newIntensityFactor := math.Sqrt((targetCTL * athlete.GetBikeThreshold()) / (duration * 100))

		// Calculate the new TSS based on the new Intensity Factor
		newTSS := duration * newIntensityFactor * newIntensityFactor * 100

		// Calculate the delta TSS
		deltaTSS := newTSS - filteredCTL

		return deltaTSS, nil
	} else if activityType == "Swim" {
		// Calculate the new Intensity Factor based on target CTL
		newIntensityFactor := math.Sqrt((targetCTL * athlete.GetSwimThreshold()) / (duration * 100))

		// Calculate the new TSS based on the new Intensity Factor
		newTSS := duration * newIntensityFactor * newIntensityFactor * 100

		// Calculate the delta TSS
		deltaTSS := newTSS - filteredCTL

		return deltaTSS, nil
	}
	return 0, errors.New("Invalid activity type")
}
