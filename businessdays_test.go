package businessdayscalculator

import (
	"testing"
	"time"
)

func TestCountBusinessDays(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int
	}{
		{
			name:     "Same day (Monday)",
			start:    time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC), // Monday
			end:      time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC), // Monday
			expected: 1,
		},
		{
			name:     "Two consecutive business days",
			start:    time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC),  // Monday
			end:      time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC), // Tuesday
			expected: 2,
		},
		{
			name:     "Full business week",
			start:    time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC),  // Monday
			end:      time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC), // Friday
			expected: 5,
		},
		{
			name:     "Across weekend",
			start:    time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC), // Friday
			end:      time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC), // Monday
			expected: 2,
		},
		{
			name:     "Two full weeks",
			start:    time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC),  // Monday
			end:      time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC), // Friday (next week)
			expected: 10,
		},
		{
			name:     "Start on weekend Saturday",
			start:    time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC), // Saturday
			end:      time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC), // Monday
			expected: 1,
		},
		{
			name:     "Start on weekend Sunday",
			start:    time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC), // Sunday
			end:      time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC), // Tuesday
			expected: 2,
		},
		{
			name:     "Both days on weekdays (Thursday to Monday)",
			start:    time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC), // Thursday
			end:      time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC), // Monday
			expected: 3,
		},
		{
			name:     "Dates in reverse order",
			start:    time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC), // Monday
			end:      time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC),  // Previous Monday
			expected: 6,
		},
		{
			name:     "Range starting on Saturday",
			start:    time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC), // Saturday
			end:      time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC), // Friday
			expected: 5,
		},
		{
			name:     "Single weekend day (Saturday)",
			start:    time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC), // Saturday
			end:      time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC), // Saturday
			expected: 0,
		},
		{
			name:     "Single weekend day (Sunday)",
			start:    time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC), // Sunday
			end:      time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC), // Sunday
			expected: 0,
		},
		{
			name:     "Entire weekend",
			start:    time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC), // Saturday
			end:      time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC), // Sunday
			expected: 0,
		},
		{
			name:     "Three weeks plus partial",
			start:    time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC),  // Monday
			end:      time.Date(2026, 2, 27, 0, 0, 0, 0, time.UTC), // Friday
			expected: 15,
		},
		{
			name:     "Starting mid-week",
			start:    time.Date(2026, 2, 11, 0, 0, 0, 0, time.UTC), // Wednesday
			end:      time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC), // Friday
			expected: 3,
		},
		{
			name:     "With time component (should ignore time)",
			start:    time.Date(2026, 2, 9, 15, 30, 45, 0, time.UTC),
			end:      time.Date(2026, 2, 13, 8, 15, 30, 0, time.UTC),
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountBusinessDays(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("CountBusinessDays(%s, %s) = %d, want %d",
					tt.start.Format("2006-01-02"), tt.end.Format("2006-01-02"), result, tt.expected)
			}
		})
	}
}

// ── Holiday-aware tests ────────────────────────────────────────────────────────
// All tests use the "US" country and 2026 holidays confirmed from the library:
//   2026-01-01 Thu  New Year's Day          (nationwide)
//   2026-01-19 Mon  MLK Day                 (nationwide)
//   2026-02-16 Mon  Washington's Birthday   (nationwide)
//   2026-05-25 Mon  Memorial Day            (nationwide)
//   2026-06-19 Fri  Juneteenth              (nationwide)
//   2026-07-04 Sat  Independence Day        (nationwide, falls on Saturday)
//   2026-09-07 Mon  Labor Day               (nationwide)
//   2026-11-11 Wed  Veterans Day            (nationwide)
//   2026-11-26 Thu  Thanksgiving Day        (nationwide)
//   2026-12-25 Fri  Christmas Day           (nationwide)

var usOpts = HolidayOptions{CountryCode: "US"}

func TestIsBusinessDayWithHolidays(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		opts     HolidayOptions
		expected bool
	}{
		{
			name:     "New Year's Day (Thursday) is not a business day",
			date:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: false,
		},
		{
			name:     "Day after New Year's (Friday) is a business day",
			date:     time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: true,
		},
		{
			name:     "Regular Monday is a business day",
			date:     time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: true,
		},
		{
			name:     "Saturday is not a business day",
			date:     time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: false,
		},
		{
			name:     "Sunday is not a business day",
			date:     time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: false,
		},
		{
			name:     "MLK Day (Monday) is not a business day",
			date:     time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: false,
		},
		{
			name:     "Christmas Day (Friday) is not a business day",
			date:     time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: false,
		},
		{
			name:     "Ignores time-of-day component",
			date:     time.Date(2026, 1, 2, 15, 30, 0, 0, time.UTC),
			opts:     usOpts,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBusinessDayWithHolidays(tt.date, tt.opts)
			if result != tt.expected {
				t.Errorf("IsBusinessDayWithHolidays(%s) = %v, want %v",
					tt.date.Format("2006-01-02"), result, tt.expected)
			}
		})
	}
}

func TestCountBusinessDaysWithHolidays(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		opts     HolidayOptions
		expected int
	}{
		{
			// Jan 1 (Thu, holiday), Jan 2 (Fri), Jan 5–9 (Mon–Fri) = 6
			name:     "First week of January skips New Year's Day",
			start:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2026, 1, 9, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: 6,
		},
		{
			// Jan 19 (Mon, MLK Day holiday), Jan 20–23 (Tue–Fri) = 4
			name:     "MLK Day week has only 4 business days",
			start:    time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: 4,
		},
		{
			// Feb 16 (Mon, Presidents' Day holiday), Feb 17–20 (Tue–Fri) = 4
			name:     "Presidents' Day week has only 4 business days",
			start:    time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: 4,
		},
		{
			// Jul 4 falls on Saturday so no weekday is displaced; week is full
			name:     "Independence Day on Saturday does not reduce business days",
			start:    time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC),  // Monday
			end:      time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC), // Friday
			opts:     usOpts,
			expected: 5,
		},
		{
			// Nov 26 (Thu, Thanksgiving) is the only holiday; Mon–Fri = 4 business days
			name:     "Thanksgiving week Mon–Fri has 4 business days",
			start:    time.Date(2026, 11, 23, 0, 0, 0, 0, time.UTC), // Monday
			end:      time.Date(2026, 11, 27, 0, 0, 0, 0, time.UTC), // Friday
			opts:     usOpts,
			expected: 4,
		},
		{
			// Reverse order should give same result
			name:     "Reverse date order gives same count",
			start:    time.Date(2026, 1, 9, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			opts:     usOpts,
			expected: 6,
		},
		{
			// Columbus Day (Oct 12) applies to US-CT but not US-FL
			name:     "Columbus Day reduces count for US-CT",
			start:    time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC), // Monday
			end:      time.Date(2026, 10, 16, 0, 0, 0, 0, time.UTC), // Friday
			opts:     HolidayOptions{CountryCode: "US", Subdivision: "US-CT"},
			expected: 4,
		},
		{
			// Columbus Day does NOT apply to US-FL
			name:     "Columbus Day does not reduce count for US-FL",
			start:    time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC), // Monday
			end:      time.Date(2026, 10, 16, 0, 0, 0, 0, time.UTC), // Friday
			opts:     HolidayOptions{CountryCode: "US", Subdivision: "US-FL"},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountBusinessDaysWithHolidays(tt.start, tt.end, tt.opts)
			if result != tt.expected {
				t.Errorf("CountBusinessDaysWithHolidays(%s, %s) = %d, want %d",
					tt.start.Format("2006-01-02"), tt.end.Format("2006-01-02"), result, tt.expected)
			}
		})
	}
}

func TestAddBusinessDaysWithHolidays(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		days     int
		opts     HolidayOptions
		expected time.Time
	}{
		{
			// Dec 24 (Thu) + 1 business day: Dec 25 (Fri) is Christmas → skip to Dec 28 (Mon)
			name:     "Skips Christmas when adding 1 day from Dec 24",
			start:    time.Date(2026, 12, 24, 0, 0, 0, 0, time.UTC),
			days:     1,
			opts:     usOpts,
			expected: time.Date(2026, 12, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			// Jan 2 (Fri) - 1 business day: Jan 1 (Thu) is New Year's → lands on Dec 31 2025 (Wed)
			name:     "Skips New Year's Day going backwards from Jan 2",
			start:    time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			days:     -1,
			opts:     usOpts,
			expected: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			// 0 days → same date
			name:     "Zero days returns same date",
			start:    time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
			days:     0,
			opts:     usOpts,
			expected: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			// Jan 15 (Thu) + 5 business days, crossing MLK Day (Jan 19 Mon):
			// Jan 16 (Fri) ✓1, Jan 17 (Sat) skip, Jan 18 (Sun) skip,
			// Jan 19 (Mon) holiday skip, Jan 20 (Tue) ✓2, Jan 21 (Wed) ✓3,
			// Jan 22 (Thu) ✓4, Jan 23 (Fri) ✓5 → Jan 23
			name:     "Adding 5 days across MLK Day lands one day later",
			start:    time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			days:     5,
			opts:     usOpts,
			expected: time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddBusinessDaysWithHolidays(tt.start, tt.days, tt.opts)
			if !result.Equal(tt.expected) {
				t.Errorf("AddBusinessDaysWithHolidays(%s, %d) = %s, want %s",
					tt.start.Format("2006-01-02"), tt.days,
					result.Format("2006-01-02"), tt.expected.Format("2006-01-02"))
			}
		})
	}
}

func TestCountBusinessDaysProperties(t *testing.T) {
	t.Run("Symmetry: order shouldn't matter", func(t *testing.T) {
		date1 := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)  // Monday
		date2 := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC) // Friday (next week)

		result1 := CountBusinessDays(date1, date2)
		result2 := CountBusinessDays(date2, date1)

		if result1 != result2 {
			t.Errorf("CountBusinessDays should be symmetric, got %d and %d", result1, result2)
		}
	})

	t.Run("Single day always returns 1 for business days", func(t *testing.T) {
		for i := range 7 {
			date := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i)
			result := CountBusinessDays(date, date)

			// Weekdays should return 1, weekends should return 0
			if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
				if result != 1 {
					t.Errorf("Single business day should return 1, got %d for %s", result, date.Weekday())
				}
			} else {
				if result != 0 {
					t.Errorf("Single weekend day should return 0, got %d for %s", result, date.Weekday())
				}
			}
		}
	})

	t.Run("Result is always non-negative", func(t *testing.T) {
		startDates := []time.Time{
			time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		}

		for _, start := range startDates {
			for days := 0; days < 100; days += 7 {
				end := start.AddDate(0, 0, days)
				result := CountBusinessDays(start, end)
				if result < 0 {
					t.Errorf("CountBusinessDays should never be negative, got %d", result)
				}
			}
		}
	})
}
