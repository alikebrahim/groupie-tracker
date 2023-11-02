package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"
)

var artists []Artist

// var locations Locations
// var dates Dates
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

	wg.Add(2)
	go artistAPI(w, wg)
	// go datesAPI(w, wg)
	// go locationsAPI(w, wg)
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

func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := GetID(r) - 1
	artist := ArtistRender{
		ID:             artists[id].ID,
		Image:          artists[id].Image,
		Name:           artists[id].Name,
		Members:        artists[id].Members,
		CreationDate:   artists[id].CreationDate,
		FirstAlbum:     artists[id].FirstAlbum,
		DatesLocations: relations.Index[id].DatesLocations,
	}

	tmpl := template.Must(template.New("artist.html").Funcs(template.FuncMap{
		"FormatText": FormatText,
	}).ParseFiles("assets/templates/artist.html"))
	// tmpl := template.Must(template.ParseFiles("assets/templates/artist.html"))
	if err := tmpl.Execute(w, artist); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
