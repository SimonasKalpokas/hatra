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
	monthDayCounts := [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	day := date.day + days
	month := date.month
	for day > monthDayCounts[month] {
		day = day - monthDayCounts[month]
		month = month + 1
	}
	return Date{date.year, month, day}
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
