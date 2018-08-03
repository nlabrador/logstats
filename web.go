package main

import (
    "encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
	"regexp"
	"os"
)

var templates = template.Must(template.ParseFiles("tmpl/view.html", "tmpl/stats.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var search = regexp.MustCompile("\\[([a-zA-Z0-9]+)\\]") 

func rootHandler(w http.ResponseWriter, r *http.Request) {
    files, err := ioutil.ReadDir("./data")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    var dates []string
    for _, f := range files {
        dates = append(dates, f.Name())
    }

    renderRoot(w, "view", dates)
}

func controllerVisitsHandler(w http.ResponseWriter, r *http.Request) {
    date := r.URL.Path[len("/controller/visits/"):]
    jsonFile, err := os.Open("./data/" + date + "/controllervisits.json")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    defer jsonFile.Close()
    byteValue, _ := ioutil.ReadAll(jsonFile)
    controllers := make(map[string]interface{})
    json.Unmarshal(byteValue, &controllers)

    renderStats(w, "stats", controllers)
}

func renderRoot(w http.ResponseWriter, tmpl string, dates []string) {
    err := templates.ExecuteTemplate(w, tmpl+".html", dates)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func renderStats(w http.ResponseWriter, tmpl string, controllers map[string]interface{}) {
    err := templates.ExecuteTemplate(w, tmpl+".html", controllers)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

type Controller struct {
    Name string
    Count int
}

func main() {
    http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/controller/visits/", controllerVisitsHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
