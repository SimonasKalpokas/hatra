package main

import (
	"flag"
	"fmt"
	"hatra/date"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

type Habit struct {
	Name string
	Days []date.Date
}

func GetHabits(habitsToExclude *[]string) []Habit {
	habits := make([]Habit, 0)

	dataFilePaths, _ := filepath.Glob("./data/*.txt")
	for _, filePath := range dataFilePaths {
		exp, _ := regexp.Compile("data\\/(\\w+)\\.txt")
		name := exp.FindStringSubmatch(filePath)[1]

		if slices.Contains(*habitsToExclude, name) {
			continue
		}

		content, _ := os.ReadFile(filePath)
		lines := strings.Split(string(content), "\n")

		dates := make([]date.Date, 0)
		for _, line := range lines {
			if line == "" {
				continue
			}
			date, err := date.ParseDate(line[:10])
			if err != nil {
				fmt.Println(line)
				panic("parsing failed")
			}
			dates = append(dates, date)
		}
		slices.SortFunc(dates, date.Compare)
		habits = append(habits, Habit{
			Name: name,
			Days: dates,
		})
	}

	return habits
}

type DayNugget struct {
	Date   date.Date
	Toggle bool
}

type HabitPageData struct {
	// TODO: Maybe think about storing by month?
	Days [][]DayNugget
}

// display

func DisplayHabitsByMonthHorizontal(habits []Habit) {
	for _, habit := range habits {
		DisplayHabitByMonthHorizontal(habit)
		fmt.Println()
	}
}

func DisplayHabitByMonthHorizontal(habit Habit) {
	today := date.Today()
	fmt.Println(habit.Name)

	monthNames := [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	monthDayCounts := [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	dateIndex := 0
	firstMonth := habit.Days[0].Month()
	lastMonth := today.Month()
	for month := firstMonth; month <= lastMonth; month = month + 1 {
		fmt.Printf("%s ", monthNames[month-1])
		for day := 1; day <= monthDayCounts[month-1]; day = day + 1 {
			if today.Month() == month && today.Day() < day {
				break
			}
			if len(habit.Days) <= dateIndex {
				fmt.Print("\033[91m□\033[39m")
				continue
			}
			date := habit.Days[dateIndex]
			if date.Month() == month && date.Day() == day {
				dateIndex = dateIndex + 1
				fmt.Print("\033[92m■\033[39m")
			} else {
				fmt.Print("\033[91m□\033[39m")
			}
		}
		fmt.Println()
	}
}

func DisplayHabitsByWeekHorizontal(habits []Habit) {
	for _, habit := range habits {
		if len(habit.Days) == 0 {
			continue
		}
		DisplayHabitByWeekHorizontal(habit)
		fmt.Println()
	}
}

func DisplayHabitByWeekHorizontal(habit Habit) {
	today := date.Today()

	fmt.Println(habit.Name)

	habitStartDay := habit.Days[0]
	date := habitStartDay.AddDays(7 - int(habitStartDay.WeekDay()))
	monthNames := [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	fmt.Printf("    ")
	for date.Year() < today.Year() ||
		(date.Year() == today.Year() && date.Month() < today.Month()) ||
		(date.Year() == today.Year() && date.Month() == today.Month() && date.Day() <= today.Day()) {
		if date.Day() <= 7 {
			fmt.Printf("%s", monthNames[date.Month()-1])
		} else {
			fmt.Printf("   ")
		}

		date = date.AddDays(7)
	}
	fmt.Println()

	weekDays := [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

	for weekDay := 1; weekDay <= 7; weekDay = weekDay + 1 {
		fmt.Printf("%s ", weekDays[weekDay-1])
		date := habitStartDay.AddDays(weekDay - int(habitStartDay.WeekDay()))
		for date.Year() < today.Year() ||
			(date.Year() == today.Year() && date.Month() < today.Month()) ||
			(date.Year() == today.Year() && date.Month() == today.Month() && date.Day() <= today.Day()) {
			if slices.Contains(habit.Days, date) {
				fmt.Printf(" \033[92m■\033[39m ")
			} else {
				fmt.Printf(" \033[91m□\033[39m ")
			}
			date = date.AddDays(7)
		}
		fmt.Println()
	}
}

func DisplayHabitsByMonthVertical(habits []Habit) {
	today := date.Today()

	for _, habit := range habits {
		monthCount := today.Month() - habit.Days[0].Month()
		fmt.Printf("%*.*s", monthCount*7, monthCount*7, habit.Name)
	}
	fmt.Println()

	monthNames := [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	monthDayCounts := [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	lastMonth := today.Month()
	for _, habit := range habits {
		firstMonth := habit.Days[0].Month()
		for month := firstMonth; month <= lastMonth; month = month + 1 {
			fmt.Printf("%s ", monthNames[month-1])
		}
		fmt.Print("  ")
	}
	fmt.Println()

	for day := 1; day <= 31; day = day + 1 {
		for _, habit := range habits {
			firstMonth := habit.Days[0].Month()
			for month := firstMonth; month <= lastMonth; month = month + 1 {
				if (day > monthDayCounts[month-1]) || (month == lastMonth && day > today.Day()) {
					fmt.Printf("    ")
					continue
				}
				if slices.Contains(habit.Days, date.NewDate(today.Year(), month, day)) {
					fmt.Printf(" \033[92m■\033[39m  ")
				} else {
					fmt.Printf(" \033[91m□\033[39m  ")
				}
			}
			fmt.Print("  ")
		}
		fmt.Println()
	}
}

func parseExcludeFlag(argument *string) []string {
	return strings.Split(*argument, ",")
}

func main() {
	directionFlag := flag.String("direction", "horizontal", "either vertical or horizontal")
	periodFlag := flag.String("period", "month", "either month or week")
	excludeFlag := flag.String("exclude", "", "a list of comma separated habits to exclude")
	flag.Parse()
	habitsToExclude := parseExcludeFlag(excludeFlag)

	habits := GetHabits(&habitsToExclude)
	if *directionFlag == "horizontal" {
		if *periodFlag == "month" {
			DisplayHabitsByMonthHorizontal(habits)
		} else if *periodFlag == "week" {
			DisplayHabitsByWeekHorizontal(habits)
		} else {
			panic("Unrecognized period")
		}
	} else if *directionFlag == "vertical" {
		if *periodFlag == "month" {
			DisplayHabitsByMonthVertical(habits)
		} else if *periodFlag == "week" {
			panic("Sorry this is not yet implemented")
		} else {
			panic("Unrecognized period")
		}
	} else {
		panic("Unrecognized direction")
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
					{Date: date.NewDate(2020, 1, 1), Toggle: false},
					{Date: date.NewDate(2020, 1, 2), Toggle: false},
					{Date: date.NewDate(2020, 1, 3), Toggle: false},
					{Date: date.NewDate(2020, 1, 4), Toggle: false},
				},
				{
					{Date: date.NewDate(2020, 1, 1), Toggle: true},
					{Date: date.NewDate(2020, 1, 2), Toggle: false},
					{Date: date.NewDate(2020, 1, 3), Toggle: false},
					{Date: date.NewDate(2020, 1, 4), Toggle: false},
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
