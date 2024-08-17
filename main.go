package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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

	return Date {
		Year: year,
		Month: month,
		Day: day,
	}, nil

}

// data

type Habit struct {
	Name	string
	Days	[]Date
}

type DayNugget struct {
	Date	Date
	Toggle	bool
}

type HabitPageData struct {
	// TODO: Maybe think about storing by month?
	Days [][]DayNugget 
}


// display

func main() {
	dataFilePaths, _ := filepath.Glob("./data/*.txt")

	for _, filePath := range dataFilePaths {
		exp, _ := regexp.Compile("data\\/(\\w+)\\.txt")
		name := exp.FindStringSubmatch(filePath)[1]
		_ = name

		content, _ := os.ReadFile(filePath)

		lines := strings.Split(string(content), "\n")
		fmt.Println(lines)
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
