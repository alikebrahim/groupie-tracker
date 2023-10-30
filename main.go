package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var artists []Artist
var locations Locations
var dates Dates
var relations Relations

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/artist/", artistHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets"))))
	fmt.Println("Server starting at port 3000")
	http.ListenAndServe(":3000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	wg := &sync.WaitGroup{}

	wg.Add(4)
	go artistAPI(w, wg)
	go datesAPI(w, wg)
	go locationsAPI(w, wg)
	go relationsAPI(w, wg)

	wg.Wait()

	// jsonMarshal, err := json.Marshal(artists)
	// if err != nil {
	// 	_, err := fmt.Fprintln(w, "error marshalling json", err)
	// 	if err != nil {
	// 		log.Fatal("err printing marshalling error", err)
	// 	}
	// }
	// _, err2 := fmt.Fprintln(w, string(jsonMarshal))
	// if err != nil {
	// 	log.Fatal(
	// 		"error printing json to webpage",
	// 		err2,
	// 	)
	// }
	tmpl := template.Must(template.ParseFiles("assets/templates/index.html"))
	if err := tmpl.Execute(w, artists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(time.Since(start))
}

type art struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := getID(r) - 1
	art := art{
		ID:             artists[id].ID,
		Image:          artists[id].Image,
		Name:           artists[id].Name,
		Members:        artists[id].Members,
		CreationDate:   artists[id].CreationDate,
		FirstAlbum:     artists[id].FirstAlbum,
		DatesLocations: relations.Index[id].DatesLocations,
	}
	tmpl := template.Must(template.ParseFiles("assets/templates/artist.html"))
	if err := tmpl.Execute(w, art); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func getID(r *http.Request) int {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0
	}
	return id
}
