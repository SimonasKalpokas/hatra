package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// data

type DayNugget struct {
	Date time.Time
	Toggle bool
}


// display

func main() {
	data := []DayNugget{
		{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2020, 1, 7, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2020, 1, 8, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC), false},
	}

	templ counts(global, sessiion int) {
		<form action="/clicked" method="post">
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.New("index").Parse(`
			<!DOCTYPE html>
			<html>
			<body>
			<h1>Days of the Week</h1>
			<ul>
			<ol>
			<li><a href="/clicked">Monday</a></li>
			<li><a href="/clicked">Tuesday</a></li>
			<li><a href="/clicked">Wednesday</a></li>
			<li><a href="/clicked">Thursday</a></li>
			<li><a href="/clicked">Friday</a></li>
			<li><a href="/clicked">Saturday</a></li>
			<li><a href="/clicked">Sunday</a></li>
			</ol>
			</ul>
			</body>
			</html>
		`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = templ.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Button Clicked!")
	})

	http.ListenAndServe(":8080", nil)
}	
