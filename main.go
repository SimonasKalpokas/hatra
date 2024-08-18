package main

import (
	"cmp"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// time

type Date struct {
	Year  int
	Month int
	Day   int
}

// Parses Date that is in format of YYYY-MM-DD
//
// TODO: move to separate module?
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

	return Date{
		Year:  year,
		Month: month,
		Day:   day,
	}, nil
}

// data

type Habit struct {
	Name string
	Days []Date
}

func GetHabits() []Habit {
	habits := make([]Habit, 0)

	dataFilePaths, _ := filepath.Glob("./data/*.txt")
	for _, filePath := range dataFilePaths {
		exp, _ := regexp.Compile("data\\/(\\w+)\\.txt")
		name := exp.FindStringSubmatch(filePath)[1]

		content, _ := os.ReadFile(filePath)
		lines := strings.Split(string(content), "\n")

		dates := make([]Date, 0)
		for _, line := range lines {
			if line == "" {
				continue
			}
			date, err := ParseDate(line)
			if err != nil {
				fmt.Println(line)
				panic("parsing failed")
			}
			dates = append(dates, date)
		}
		slices.SortFunc(dates, func(a, b Date) int {
			compareYear := cmp.Compare(a.Year, b.Year)
			if compareYear != 0 {
				return compareYear
			}
			compareMonth := cmp.Compare(a.Month, b.Month)
			if compareMonth != 0 {
				return compareMonth
			}
			compareDay := cmp.Compare(a.Day, b.Day)
			return compareDay
		})
		habits = append(habits, Habit{
			Name: name,
			Days: dates,
		})
	}

	return habits
}

type DayNugget struct {
	Date   Date
	Toggle bool
}

type HabitPageData struct {
	// TODO: Maybe think about storing by month?
	Days [][]DayNugget
}

func Now() Date {
	time := time.Now()
	return Date{
		Year:  time.Year(),
		Month: int(time.Month()),
		Day:   time.Day(),
	}
}

// display

func DisplayHabit(habit Habit) {
	today := Now()
	fmt.Println(habit.Name)

	monthNames := [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	monthDayCounts := [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	dateIndex := 0
	firstMonth := habit.Days[0].Month
	lastMonth := today.Month
	for month := firstMonth; month <= lastMonth; month = month + 1 {
		fmt.Printf("%s ", monthNames[month-1])
		for day := 1; day <= monthDayCounts[month-1]; day = day + 1 {
			if today.Month == month && today.Day < day {
				break
			}
			if len(habit.Days) <= dateIndex {
				fmt.Print(".")
				continue
			}
			date := habit.Days[dateIndex]
			if date.Month == month && date.Day == day {
				dateIndex = dateIndex + 1
				fmt.Print("+")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func main() {
	habits := GetHabits()
	for _, habit := range habits {
		DisplayHabit(habit)
		fmt.Println()
	}
}

func main1() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("index.html"))

		data := HabitPageData{
			Days: [][]DayNugget{
				{
					{Date: Date{Year: 2020, Month: 1, Day: 1}, Toggle: false},
					{Date: Date{Year: 2020, Month: 1, Day: 2}, Toggle: false},
					{Date: Date{Year: 2020, Month: 1, Day: 3}, Toggle: false},
					{Date: Date{Year: 2020, Month: 1, Day: 4}, Toggle: false},
				},
				{
					{Date: Date{Year: 2020, Month: 1, Day: 1}, Toggle: true},
					{Date: Date{Year: 2020, Month: 1, Day: 2}, Toggle: false},
					{Date: Date{Year: 2020, Month: 1, Day: 3}, Toggle: false},
					{Date: Date{Year: 2020, Month: 1, Day: 4}, Toggle: false},
				},
			},
		}
		templ.Execute(w, data)
	})

	http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Button Clicked!")
	})

	http.ListenAndServe(":8080", nil)
}
