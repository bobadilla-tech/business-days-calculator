# Business Days Calculator

A Go package to calculate the number of business days (Monday to Friday) between two dates, with optional public holiday awareness for 60+ countries.

[![Go Reference](https://pkg.go.dev/badge/github.com/bobadilla-tech/business-days-calculator.svg)](https://pkg.go.dev/github.com/bobadilla-tech/business-days-calculator)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobadilla-tech/business-days-calculator)](https://goreportcard.com/report/github.com/bobadilla-tech/business-days-calculator)
[![License](https://img.shields.io/github/license/bobadilla-tech/business-days-calculator)](LICENSE)

## Features

- 📅 **Business day calculation** - Counts only Monday to Friday
- 🌍 **Holiday awareness** - Exclude public holidays for 60+ countries and their subdivisions
- ⚡ **Fast and efficient** - Holiday lookups are O(1) via pre-built date sets
- 🧪 **Well tested** - Comprehensive test coverage
- 💻 **Simple API** - Easy to integrate

## Installation

```bash
go get github.com/bobadilla-tech/business-days-calculator
```

## Usage

### Basic — weekends only

```go
package main

import (
    "fmt"
    "time"

    businessdayscalculator "github.com/bobadilla-tech/business-days-calculator"
)

func main() {
    start := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)  // Monday
    end   := time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC) // Next Monday

    days := businessdayscalculator.CountBusinessDays(start, end)
    fmt.Printf("Business days: %d\n", days) // 6
}
```

### With public holidays

```go
package main

import (
    "fmt"
    "time"

    businessdayscalculator "github.com/bobadilla-tech/business-days-calculator"
)

func main() {
    opts := businessdayscalculator.HolidayOptions{
        CountryCode: "US",
        Subdivision: "US-NY", // optional; omit for nationwide holidays only
    }

    start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
    end   := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

    // Count business days, skipping weekends and US public holidays
    count := businessdayscalculator.CountBusinessDaysWithHolidays(start, end, opts)
    fmt.Printf("Business days in January 2026: %d\n", count)

    // Add 10 business days to a date, skipping weekends and holidays
    deadline := businessdayscalculator.AddBusinessDaysWithHolidays(start, 10, opts)
    fmt.Printf("10 business days after Jan 1: %s\n", deadline.Format("2006-01-02"))

    // Check if a specific date is a business day
    isOpen := businessdayscalculator.IsBusinessDayWithHolidays(
        time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC), // MLK Day
        opts,
    )
    fmt.Printf("Jan 19 is a business day: %v\n", isOpen) // false
}
```

## API Reference

### Types

#### `HolidayOptions`

Configures the country and optional subdivision used for holiday lookups.

```go
type HolidayOptions struct {
    CountryCode string // ISO 3166-1 alpha-2 code, e.g. "US", "CA", "GB"
    Subdivision string // ISO 3166-2 code, e.g. "US-NY", "GB-ENG" (optional)
}
```

When `Subdivision` is set, nationwide holidays always apply, and subdivision-specific holidays are filtered to only those matching the given subdivision. When `Subdivision` is empty, all holidays for the country are included.

### Functions

#### `CountBusinessDays(start, end time.Time) int`

Returns the number of business days (Monday–Friday) between two dates, inclusive. Weekends are excluded; public holidays are **not** considered. The order of dates does not matter.

```go
// 6 business days (Mon Feb 9 – Mon Feb 16, inclusive)
days := businessdayscalculator.CountBusinessDays(
    time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC),
    time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC),
)
```

#### `CountBusinessDaysWithHolidays(start, end time.Time, opts HolidayOptions) int`

Same as `CountBusinessDays` but also excludes public holidays for the given country/subdivision.

```go
opts := businessdayscalculator.HolidayOptions{CountryCode: "US"}

// Jan 2026 has New Year's Day (Jan 1, Thu) and MLK Day (Jan 19, Mon)
count := businessdayscalculator.CountBusinessDaysWithHolidays(
    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
    time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
    opts,
)
```

#### `AddBusinessDaysWithHolidays(start time.Time, n int, opts HolidayOptions) time.Time`

Returns the date that is `n` business days after `start`, skipping weekends and public holidays. Negative values of `n` move backwards in time.

```go
opts := businessdayscalculator.HolidayOptions{CountryCode: "US"}

// Dec 24 (Thu) + 1 business day → Dec 28 (Mon), because Dec 25 is Christmas
result := businessdayscalculator.AddBusinessDaysWithHolidays(
    time.Date(2026, 12, 24, 0, 0, 0, 0, time.UTC),
    1,
    opts,
)
```

#### `IsBusinessDayWithHolidays(date time.Time, opts HolidayOptions) bool`

Reports whether `date` is a business day, accounting for weekends and public holidays.

```go
opts := businessdayscalculator.HolidayOptions{CountryCode: "US"}

// MLK Day is a Monday but a public holiday → false
open := businessdayscalculator.IsBusinessDayWithHolidays(
    time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC),
    opts,
)
```

### Supported country codes (examples)

| Code | Country        | Subdivision support                    |
| ---- | -------------- | -------------------------------------- |
| `US` | United States  | `US-NY`, `US-CA`, `US-TX`, …           |
| `CA` | Canada         | `CA-ON`, `CA-QC`, `CA-BC`, …           |
| `GB` | United Kingdom | `GB-ENG`, `GB-SCT`, `GB-WLS`, `GB-NIR` |
| `AU` | Australia      | `AU-NSW`, `AU-VIC`, `AU-QLD`, …        |
| `DE` | Germany        | `DE-BY`, `DE-NW`, `DE-BE`, …           |

Holiday data is provided by [github.com/bobadilla-tech/holidays-per-country](https://github.com/bobadilla-tech/holidays-per-country).

## How It Works

### `CountBusinessDays`

1. Normalise both dates to midnight and swap if needed
2. Count full weeks × 5, then iterate the remaining days checking the weekday

### `CountBusinessDaysWithHolidays` / `IsBusinessDayWithHolidays` / `AddBusinessDaysWithHolidays`

1. Fetch all holidays for the date range in a single call to `GetHolidaysInRange`
2. Store them in a `map[string]struct{}` keyed by `"YYYY-MM-DD"` for O(1) lookup
3. Iterate day-by-day, skipping weekends and any date in the holiday set
4. `AddBusinessDaysWithHolidays` lazily expands the holiday set if the walk goes beyond the pre-fetched range

## Testing

```bash
go test -v
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

- Inspired by the need for simple, efficient business day calculation in Go
- Holiday data powered by [holidays-per-country](https://github.com/bobadilla-tech/holidays-per-country)

## Related Projects

- [holidays-per-country](https://github.com/bobadilla-tech/holidays-per-country) - Public holiday data for 60+ countries (Go)
- [dateutil](https://github.com/dateutil/dateutil) - Powerful extensions to the standard datetime module (Python)
