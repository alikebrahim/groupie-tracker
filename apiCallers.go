package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const (
	ArtistAPI = "https://groupietrackers.herokuapp.com/api/artists"
	// LocationsAPI = "https://groupietrackers.herokuapp.com/api/locations"
	// DatesAPI     = "https://groupietrackers.herokuapp.com/api/dates"
	RelationsAPI = "https://groupietrackers.herokuapp.com/api/relation"
)

func artistAPI(w http.ResponseWriter, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := http.Get(ArtistAPI)
	if err != nil {
		_, err := fmt.Fprintln(w, "http.Get", "Error getting api:", err)
		if err != nil {
			log.Fatal("Error writing response") // error 500
		}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		_, err := fmt.Fprintln(w, "http.Get", "Error body reading api:", err)
		if err != nil {
			log.Fatal("error writing response") // 500 error
		}
	}
	err = json.Unmarshal(body, &Artists)
	if err != nil {
		_, err := fmt.Fprintln(w, "error ummarshalling JSON", err)
		if err != nil {
			log.Fatal("error writing response") // 500 error
		}
	}
}

// func locationsAPI(w http.ResponseWriter, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	res, err := http.Get(LocationsAPI)
// 	if err != nil {
// 		_, err := fmt.Fprintln(w, "http.Get", "Error getting api:", err)
// 		if err != nil {
// 			log.Fatal("Error writing response") // error 500
// 		}
// 	}
// 	defer res.Body.Close()
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		_, err := fmt.Fprintln(w, "http.Get", "Error body reading api:", err)
// 		if err != nil {
// 			log.Fatal("error writing response") // 500 error
// 		}
// 	}
// 	err = json.Unmarshal(body, &locations)
// 	if err != nil {
// 		_, err := fmt.Fprintln(w, "error ummarshalling JSON", err)
// 		if err != nil {
// 			log.Fatal("error writing response") // 500 error
// 		}
// 	}
// }
//
// func datesAPI(w http.ResponseWriter, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	res, err := http.Get(DatesAPI)
// 	if err != nil {
// 		_, err := fmt.Fprintln(w, "http.Get", "Error getting api:", err)
// 		if err != nil {
// 			log.Fatal("Error writing response") // error 500
// 		}
// 	}
// 	defer res.Body.Close()
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		_, err := fmt.Fprintln(w, "http.Get", "Error body reading api:", err)
// 		if err != nil {
// 			log.Fatal("error writing response") // 500 error
// 		}
// 	}
// 	err = json.Unmarshal(body, &dates)
// 	if err != nil {
// 		_, err := fmt.Fprintln(w, "error ummarshalling JSON", err)
// 		if err != nil {
// 			log.Fatal("error writing response") // 500 error
// 		}
// 	}
// }

func relationsAPI(w http.ResponseWriter, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := http.Get(RelationsAPI)
	if err != nil {
		_, err := fmt.Fprintln(w, "http.Get", "Error getting api:", err)
		if err != nil {
			log.Fatal("Error writing response") // error 500
		}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		_, err := fmt.Fprintln(w, "http.Get", "Error body reading api:", err)
		if err != nil {
			log.Fatal("error writing response") // 500 error
		}
	}
	err = json.Unmarshal(body, &relations)
	if err != nil {
		_, err := fmt.Fprintln(w, "error ummarshalling JSON", err)
		if err != nil {
			log.Fatal("error writing response") // 500 error
		}
	}
}
