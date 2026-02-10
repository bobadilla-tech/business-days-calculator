package businessdayscalculator

import "time"

func CalculateBusinessDays(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	totalDays := int(end.Sub(start)/(24*time.Hour)) + 1

	fullWeeks := totalDays / 7
	remainingDays := totalDays % 7

	businessDays := fullWeeks * 5

	weekday := int(start.Weekday())

	for i := range remainingDays {
		day := (weekday + i) % 7
		if day != int(time.Saturday) && day != int(time.Sunday) {
			businessDays++
		}
	}

	return businessDays
}
