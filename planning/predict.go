package planning

import (
	"atc/models"
	"math"
)

// deltaTSS calculates the necessary change in TSS to reach a target CTL and volume.
func deltaTSS(athlete models.Athlete, targetVolume float64, targetCTL float64, activityType string) (float64, error) {
	filteredActivities := models.FilterActivitiesByType(athlete.Activities, activityType)

	filteredCTL := models.CalculateCTL(filteredActivities, 42)
	filteredDuration := models.CalculateDurationHrs(filteredActivities)

	// Calculate the necessary change in CTL
	deltaCTL := targetCTL - filteredCTL

	// Estimate the duration in hours (assuming linear relation with volume)
	duration := targetVolume / filteredCTL * filteredDuration

	// Calculate the new Intensity Factor based on target TSS
	newIntensityFactor := math.Sqrt((targetTSS * athlete.Threshold) / (duration * 100))

	// Return the target volume and intensity factor
	return targetVolume, newIntensityFactor

}
