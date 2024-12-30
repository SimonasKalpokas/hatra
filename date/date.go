package date

import (
	"cmp"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	year  int
	month int
	day   int
}

func NewDate(year, month, day int) Date {
	return Date{year, month, day}
}

const (
	daysPer400Years = 365*400 + 97
	daysPer100Years = 365*100 + 24
	daysPer4Years   = 365*4 + 1
)

func daysSinceEpochForYear(year int) int {
	y := year - 2000

	// Add in days from 400-year cycles.
	n := y / 400
	y -= 400 * n
	d := daysPer400Years * n

	// Add in 100-year cycles.
	n = y / 100
	y -= 100 * n
	d += daysPer100Years * n

	// Add in 4-year cycles.
	n = y / 4
	y -= 4 * n
	d += daysPer4Years * n

	// Add in non-leap years.
	n = y
	d += 365 * n

	return d
}

// TODO: explain values and operation with comments
// add tests and refactor
func daysSinceEpoch(date Date) int {
	days := daysSinceEpochForYear(date.year)

	monthDayCounts := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	for month := 0; month < date.month-1; month++ {
		days += monthDayCounts[month]
	}

	if isLeapYear(date.year) && date.month > 2 {
		days++
	}

	days += date.day - 2

	return days
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func (date Date) WeekDay() int {
	days := daysSinceEpoch(date)
	return (days+6)%7 + 1
}

// Parses Date that is in format of YYYY-MM-DD
func ParseDate(input string) (Date, error) {
	parts := strings.Split(input, "-")

	if len(parts) != 3 ||
		len(parts[0]) != 4 ||
		len(parts[1]) != 2 ||
		len(parts[2]) != 2 {
		return Date{}, errors.New(fmt.Sprintf("bad string format: %s", input))
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil || year < 0 {
		return Date{}, errors.New(fmt.Sprintf("bad string format: %s", input))
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil || month > 12 || month < 1 {
		return Date{}, errors.New(fmt.Sprintf("bad string format: %s", input))
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil || month > 12 || month < 1 {
		return Date{}, errors.New(fmt.Sprintf("bad string format: %s", input))
	}

	return NewDate(year, month, day), nil
}

func Compare(a, b Date) int {
	compareYear := cmp.Compare(a.year, b.year)
	if compareYear != 0 {
		return compareYear
	}
	compareMonth := cmp.Compare(a.month, b.month)
	if compareMonth != 0 {
		return compareMonth
	}
	compareDay := cmp.Compare(a.day, b.day)
	return compareDay
}

func Today() Date {
	now := time.Now()
	return Date{
		year:  now.Year(),
		month: int(now.Month()),
		day:   now.Day(),
	}
}

func (date Date) AddDays(days int) Date {
	// TODO: handle leap years
	monthDayCounts := [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	day := date.day + days
	month := date.month
	year := date.year
	for day > monthDayCounts[month-1] {
		day = day - monthDayCounts[month-1]
		month = month + 1
		for month > 12 {
			month -= 12
			year = year + 1
		}
	}
	return Date{year, month, day}
}

func (date Date) Year() int {
	return date.year
}

func (date Date) Month() int {
	return date.month
}

func (date Date) Day() int {
	return date.day
}
