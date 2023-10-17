package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var artists []Artist
var locations Locations
var dates Dates
var relations Relations

func main() {
	http.HandleFunc("/", indexHandler)
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

	jsonMarshal, err := json.Marshal(artists)
	if err != nil {
		_, err := fmt.Fprintln(w, "error marshalling json", err)
		if err != nil {
			log.Fatal("err printing marshalling error", err)
		}
	}
	_, err2 := fmt.Fprintln(w, string(jsonMarshal))
	if err != nil {
		log.Fatal(
			"error printing json to webpage",
			err2,
		)
	}
	fmt.Println(time.Since(start))
}