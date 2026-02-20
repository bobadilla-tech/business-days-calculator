package businessdayscalculator

import (
	"time"

	holidays "github.com/bobadilla-tech/holidays-per-country"
)

// HolidayOptions configures which country (and optional subdivision) to use
// when determining public holidays.
type HolidayOptions struct {
	// CountryCode is the ISO 3166-1 alpha-2 country code (e.g. "US", "CA", "GB").
	CountryCode string
	// Subdivision is the optional ISO 3166-2 subdivision code (e.g. "US-NY", "GB-ENG").
	// When set, only nationwide holidays and holidays that apply to this subdivision
	// are considered. When empty, all holidays for the country are used.
	Subdivision string
}

// CountBusinessDaysWithHolidays counts the number of business days between start
// and end (both inclusive), excluding weekends and public holidays for the given
// country/subdivision.
func CountBusinessDaysWithHolidays(start, end time.Time, opts HolidayOptions) int {
	start = truncateToDay(start)
	end = truncateToDay(end)

	if start.After(end) {
		start, end = end, start
	}

	hs := buildHolidaySet(opts.CountryCode, opts.Subdivision, start, end)

	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if isWeekday(d) {
			if _, isHoliday := hs[dateKey(d)]; !isHoliday {
				count++
			}
		}
	}
	return count
}

// AddBusinessDaysWithHolidays returns the date that is n business days after
// start, skipping weekends and public holidays for the given country/subdivision.
// Negative values of n move backwards in time.
func AddBusinessDaysWithHolidays(start time.Time, n int, opts HolidayOptions) time.Time {
	if n == 0 {
		return start
	}

	step := 1
	if n < 0 {
		step = -1
		n = -n
	}

	// Pre-fetch holidays for an estimated range (worst case ~260 business days/year).
	estimatedCalendarDays := n*2 + 30
	rangeStart := start
	rangeEnd := start.AddDate(0, 0, estimatedCalendarDays*step)
	if step < 0 {
		rangeStart, rangeEnd = rangeEnd, start
	}

	hs := buildHolidaySet(opts.CountryCode, opts.Subdivision, rangeStart, rangeEnd)

	current := truncateToDay(start)
	counted := 0
	for counted < n {
		current = current.AddDate(0, 0, step)

		// Lazily expand the holiday set if we walk past the pre-fetched range.
		if step > 0 && current.After(rangeEnd) {
			newEnd := rangeEnd.AddDate(1, 0, 0)
			for k, v := range buildHolidaySet(opts.CountryCode, opts.Subdivision, rangeEnd.AddDate(0, 0, 1), newEnd) {
				hs[k] = v
			}
			rangeEnd = newEnd
		} else if step < 0 && current.Before(rangeStart) {
			newStart := rangeStart.AddDate(-1, 0, 0)
			for k, v := range buildHolidaySet(opts.CountryCode, opts.Subdivision, newStart, rangeStart.AddDate(0, 0, -1)) {
				hs[k] = v
			}
			rangeStart = newStart
		}

		if isWeekday(current) {
			if _, isHoliday := hs[dateKey(current)]; !isHoliday {
				counted++
			}
		}
	}
	return current
}

// IsBusinessDayWithHolidays reports whether date is a business day, taking
// weekends and public holidays for the given country/subdivision into account.
func IsBusinessDayWithHolidays(date time.Time, opts HolidayOptions) bool {
	d := truncateToDay(date)
	if !isWeekday(d) {
		return false
	}
	hs := buildHolidaySet(opts.CountryCode, opts.Subdivision, d, d)
	_, isHoliday := hs[dateKey(d)]
	return !isHoliday
}

// buildHolidaySet fetches holidays for the given range and returns a set keyed
// by "YYYY-MM-DD" for O(1) lookup. When subdivision is non-empty, subdivision-
// specific holidays that do not match it are filtered out; nationwide holidays
// (empty Subdivisions slice) always apply.
func buildHolidaySet(countryCode, subdivision string, start, end time.Time) map[string]struct{} {
	raw := holidays.GetHolidaysInRange(countryCode, start, end)
	set := make(map[string]struct{}, len(raw))
	for _, h := range raw {
		if subdivision != "" && len(h.Subdivisions) > 0 {
			if !containsSubdivision(h.Subdivisions, subdivision) {
				continue
			}
		}
		set[dateKey(h.Date)] = struct{}{}
	}
	return set
}

// isWeekday reports whether t falls on Monday–Friday.
func isWeekday(t time.Time) bool {
	wd := t.Weekday()
	return wd != time.Saturday && wd != time.Sunday
}

// truncateToDay strips the time-of-day component from t.
func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// dateKey returns a string key in "YYYY-MM-DD" format.
func dateKey(t time.Time) string {
	return t.Format("2006-01-02")
}

// containsSubdivision reports whether target is in the subdivisions list.
func containsSubdivision(subdivisions []string, target string) bool {
	for _, s := range subdivisions {
		if s == target {
			return true
		}
	}
	return false
}

func CountBusinessDays(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}

	start = truncateToDay(start)
	end = truncateToDay(end)

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
