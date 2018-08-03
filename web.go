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

    var t *template.Template
    t, _ = template.ParseFiles("tmpl/layout.html", "tmpl/view.html")
    err = t.ExecuteTemplate(w, "layout", dates)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
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

    var t *template.Template
    t, _ = template.ParseFiles("tmpl/layout.html", "tmpl/stats.html")
    err = t.ExecuteTemplate(w, "layout", controllers)
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
