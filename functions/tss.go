package functions

// includes functions for tss, ctl, and atl calculations

// tss is the training stress score, a measure of the overall training load of a workout
// basically a measure of the duration of a workout spent at threshold.
// a tss of 100 is one hour at threshold, 50 is half an hour, etc.

// intensity ("if") is the intensity factor, a measure of how hard the workout was,
// with an intensity factor of 1 being a threshold workout.

// ctl (sometimes called "fitness") is chronic training load, or the rolling average of tss over 42 days.

// atl (sometimes called "fatigue") is acute training load, or the rolling average of tss over 7 days.

// tsb (sometimes called "form") is training stress balance, or the difference between ctl and atl.
