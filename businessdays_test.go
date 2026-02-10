package businessdayscalculator

import (
	"testing"
	"time"
)

func TestCalculateBusinessDays(t *testing.T) {
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
			result := CalculateBusinessDays(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("CalculateBusinessDays(%s, %s) = %d, want %d",
					tt.start.Format("2006-01-02"), tt.end.Format("2006-01-02"), result, tt.expected)
			}
		})
	}
}

func TestCalculateBusinessDaysProperties(t *testing.T) {
	t.Run("Symmetry: order shouldn't matter", func(t *testing.T) {
		date1 := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)  // Monday
		date2 := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC) // Friday (next week)

		result1 := CalculateBusinessDays(date1, date2)
		result2 := CalculateBusinessDays(date2, date1)

		if result1 != result2 {
			t.Errorf("CalculateBusinessDays should be symmetric, got %d and %d", result1, result2)
		}
	})

	t.Run("Single day always returns 1 for business days", func(t *testing.T) {
		for i := range 7 {
			date := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i)
			result := CalculateBusinessDays(date, date)

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
				result := CalculateBusinessDays(start, end)
				if result < 0 {
					t.Errorf("BusinessDays should never be negative, got %d", result)
				}
			}
		}
	})
}
