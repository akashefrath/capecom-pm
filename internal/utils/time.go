package utils

import "time"

func GetTodayRange() (time.Time, time.Time) {
	now := time.Now()

	// Get 00:00:00 of the current day in the local timezone
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Add 24 hours to get the start of tomorrow
	startOfTomorrow := startOfToday.AddDate(0, 0, 1)

	return startOfToday, startOfTomorrow
}
