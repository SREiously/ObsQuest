package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
)

type App struct {
	Name         string  `json:"name"`
	Availability float64 `json:"availability"`
	Errors       int64   `json:"errors"`
	Type         string  `json:"type"`
	Trend        float64 `json:"trend"`
}

func loadApps(filename string) []App {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var apps []App
	if err := json.Unmarshal(data, &apps); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	return apps
}

var apps = []App{}

type PageData struct {
	Heroes   []App
	Villains []App
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].Availability > apps[j].Availability
	})

	var heroes, villains []App
	for i := len(apps) - 1; i >= 0; i-- {
		if apps[i].Type == "Hero" && len(heroes) < 10 {
			heroes = append(heroes, apps[i])
		} else if apps[i].Type == "Villain" && len(villains) < 10 {
			villains = append(villains, apps[i])
		}
	}

	tmpl, err := template.ParseFiles("templates/report.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Heroes:   heroes,
		Villains: villains,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	apps = loadApps("canned.json")
	http.HandleFunc("/", mainHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
