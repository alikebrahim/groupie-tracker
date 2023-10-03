package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var artists []Artist
var artistsComplete []ArtistComplete

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":3000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	artistAPI := "https://groupietrackers.herokuapp.com/api/artists"
	res, err := http.Get(artistAPI)
	if err != nil {
		fmt.Fprintln(w, "http.Get", "Error getting api:", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(w, "http.Get", "Error body reading api:", err)
	}
	//fmt.Fprintln(w, string(body))
	err = json.Unmarshal(body, &artists)
	if err != nil {
		fmt.Fprintln(w, "error ummarshalling JSON", err)
	}
	//fmt.Fprintln(w, artists)
	//for _, item := range artists {
	//	fmt.Fprintln(w, item)
	//}
	for _, item := range artists {
		artist := ArtistComplete{
			ID:           item.ID,
			Image:        item.Image,
			Name:         item.Name,
			Members:      item.Members,
			CreationDate: item.CreationDate,
			FirstAlbum:   item.FirstAlbum,
			Locations:    locationsAPI(item.Locations),
			ConcertDates: datesAPI(item.ConcertDates),
			//Relations:    relationsAPI(item.Relations),
			Relations: item.Relations,
		}
		artistsComplete = append(artistsComplete, artist)

	}
	jsonMarshal, err := json.Marshal(artistsComplete)
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
}

func locationsAPI(loc string) Locations {
	locations := Locations{}
	API := loc
	res, err := http.Get(API)
	if err != nil {
		fmt.Println("error", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error", err)
	}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println("error unmarshalling to locations", err)
	}

	return locations
}

func datesAPI(dat string) Dates {
	dates := Dates{}
	API := dat
	res, err := http.Get(API)
	if err != nil {
		fmt.Println("error", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error", err)
	}
	err = json.Unmarshal(body, &dates)
	if err != nil {
		fmt.Println("error unmarshalling to dates", err)
	}

	return dates
}

func relationsAPI(rel string) Relations {
	relations := Relations{}
	API := rel
	res, err := http.Get(API)
	if err != nil {
		fmt.Println("error", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error", err)
	}
	err = json.Unmarshal(body, &relations)
	if err != nil {
		fmt.Println("error unmarshalling to relations", err)
	}

	return relations

}
