package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// time

type Date struct {
	Year  int
	Month int
	Day   int
}

// data

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
