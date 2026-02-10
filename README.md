# Business Days Calculator

A Go package to calculate the number of business days (Monday to Friday) between two dates.

[![Go Reference](https://pkg.go.dev/badge/github.com/bobadilla-tech/business-days-calculator.svg)](https://pkg.go.dev/github.com/bobadilla-tech/business-days-calculator)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobadilla-tech/business-days-calculator)](https://goreportcard.com/report/github.com/bobadilla-tech/business-days-calculator)
[![License](https://img.shields.io/github/license/bobadilla-tech/business-days-calculator)](LICENSE)

## Features

- 📅 **Business day calculation** - Counts only Monday to Friday
- 🚀 **Zero dependencies** - Pure Go implementation
- ⚡ **Fast and efficient** - Simple, reliable algorithm
- 🧪 **Well tested** - Comprehensive test coverage
- 💻 **Simple API** - Easy to integrate

## Installation

```bash
go get github.com/bobadilla-tech/business-days-calculator
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "time"
    "github.com/bobadilla-tech/business-days-calculator"
)

func main() {
    start := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC) // Monday
    end := time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC)   // Next Monday
    days := businessdayscalculator.CalculateBusinessDays(start, end)
    fmt.Printf("Business days: %d\n", days)
}
```

### API Reference

#### `CalculateBusinessDays(start, end time.Time) int`

Returns the number of business days (Monday to Friday) between two dates, inclusive. The order of dates does not matter.

```go
businessdayscalculator.CalculateBusinessDays(start, end) // returns business day count
```

- Weekends (Saturday, Sunday) are excluded
- Time component is ignored (only the date matters)
- If start and end are the same business day, returns 1
- If both are on a weekend, returns 0

## How It Works

1. **Date Normalization**: Both dates are normalized to midnight (00:00:00)
2. **Order Handling**: If start is after end, the dates are swapped
3. **Full Weeks**: Calculates the number of full weeks and multiplies by 5 (business days per week)
4. **Partial Week**: Counts business days in the remaining days
5. **Excludes weekends**: Only Monday-Friday are counted

## Testing

Run the tests:

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

## Related Projects

- [dateutil](https://github.com/dateutil/dateutil) - Powerful extensions to the standard datetime module (Python)
