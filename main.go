package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// *************************************************************************//
// Note 1: activate the rest of the APIs before audit as they are required //
// Note 2: Ensure http status errors handling                             //
// **********************************************************************//

var Artists []Artist

// var locations Locations
// var dates Dates
var relations Relations

func main() {
	http.HandleFunc("/", indexRouter)
	http.HandleFunc("/artist/", artistHandler)
	http.HandleFunc("/filter-artist/", filterHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets"))))
	fmt.Println("Server starting at port 3000")
	http.ListenAndServe(":3000", nil)
}

func indexRouter(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		indexHandler(w, r)
	case "POST":
		filterHandler(w, r)
	}
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	var filteredArtists []ArtistRender
	creationYearFrom, creationYearTo, firstAlbumYearFrom, firstAlbumYearTo, location, membersNum := FormParse(r)

	// filter loop
	if creationYearFrom >= 1958 {
		for i, artist := range Artists {
			location = strings.Split(location, ",")[0]
			if artist.CreationDate >= creationYearFrom && artist.CreationDate <= creationYearTo {
				artistMatch := MakeArtistRender(i)
				filteredArtists = append(filteredArtists, artistMatch)
			}
		}
		if firstAlbumYearFrom >= 1963 {
			var newfilteredArtists []ArtistRender
			for i := range filteredArtists {
				artistFirstAlbumYear, _ := GetInfo(i, strings.ToLower(location), Artists[filteredArtists[i].ID-1])
				if artistFirstAlbumYear >= firstAlbumYearFrom && artistFirstAlbumYear <= firstAlbumYearTo {
					newfilteredArtists = append(newfilteredArtists, filteredArtists[i])
				}
			}
			filteredArtists = newfilteredArtists
		}
		if location != "" {
			var newfilteredArtists []ArtistRender
			for i := range filteredArtists {
				location = strings.Split(location, ",")[0]
				_, loc := GetInfoFiltered(filteredArtists[i].ID-1, strings.ToLower(location), filteredArtists[i])
				if loc {
					newfilteredArtists = append(newfilteredArtists, filteredArtists[i])
				}
			}
			filteredArtists = newfilteredArtists
		}
		if len(membersNum) > 0 {
			var newfilteredArtists []ArtistRender
			for _, n := range membersNum {
				for i := range filteredArtists {
					if len(filteredArtists[i].Members) == n {
						newfilteredArtists = append(newfilteredArtists, filteredArtists[i])
					}
				}
			}
			filteredArtists = newfilteredArtists
		}
	} else if firstAlbumYearFrom >= 1963 {
		for id, artist := range Artists {
			artistFirstAlbumYear, _ := GetInfo(id, strings.ToLower(location), artist)
			if artistFirstAlbumYear >= firstAlbumYearFrom && artistFirstAlbumYear <= firstAlbumYearTo {
				artistMatch := MakeArtistRender(id)
				filteredArtists = append(filteredArtists, artistMatch)
			}
		}
		if location != "" {
			var newfilteredArtists []ArtistRender
			for i := range filteredArtists {
				location = strings.Split(location, ",")[0]
				_, loc := GetInfoFiltered(filteredArtists[i].ID-1, strings.ToLower(location), filteredArtists[i])
				if loc {
					newfilteredArtists = append(newfilteredArtists, filteredArtists[i])
				}
			}
			filteredArtists = newfilteredArtists
		}
		if len(membersNum) > 0 {
			var newfilteredArtists []ArtistRender
			for _, n := range membersNum {
				for i := range filteredArtists {
					if len(filteredArtists[i].Members) == n {
						newfilteredArtists = append(newfilteredArtists, filteredArtists[i])
					}
				}
			}
			filteredArtists = newfilteredArtists
		}
	} else if location != "" {
		for id, artist := range Artists {
			location = strings.Split(location, ",")[0]
			_, loc := GetInfo(id, strings.ToLower(location), artist)
			if loc {
				artistMatch := MakeArtistRender(id)
				filteredArtists = append(filteredArtists, artistMatch)
			}
		}
		if firstAlbumYearFrom >= 1963 {
			var newfilteredArtists []ArtistRender
			for i := range filteredArtists {
				artistFirstAlbumYear, _ := GetInfo(i, strings.ToLower(location), Artists[filteredArtists[i].ID-1])
				if artistFirstAlbumYear >= firstAlbumYearFrom && artistFirstAlbumYear <= firstAlbumYearTo {
					newfilteredArtists = append(newfilteredArtists, filteredArtists[i])
				}
			}
			filteredArtists = newfilteredArtists
		}
		if location != "" {
			var newfilteredArtists []ArtistRender
			for i := range filteredArtists {
				location = strings.Split(location, ",")[0]
				_, loc := GetInfoFiltered(filteredArtists[i].ID-1, strings.ToLower(location), filteredArtists[i])
				if loc {
					newfilteredArtists = append(newfilteredArtists, filteredArtists[i])
				}
			}
			filteredArtists = newfilteredArtists
		}
		if len(membersNum) > 0 {
			var newfilteredArtists []ArtistRender
			for _, n := range membersNum {
				for i := range filteredArtists {
					if len(filteredArtists[i].Members) == n {
						newfilteredArtists = append(newfilteredArtists, filteredArtists[i])
					}
				}
			}
			filteredArtists = newfilteredArtists
		}
	} else if len(membersNum) > 0 {
		for _, n := range membersNum {
			for id, artist := range Artists {
				if len(artist.Members) == n {
					artistMatch := MakeArtistRender(id)
					filteredArtists = append(filteredArtists, artistMatch)
				}
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("assets/templates/filter.html"))
	if err := tmpl.Execute(w, filteredArtists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go artistAPI(w, wg)
	// ***********************************//
	// Note: To be activated before audit//
	// *********************************//
	// go datesAPI(w, wg)
	// go locationsAPI(w, wg)
	go relationsAPI(w, wg)

	wg.Wait()

	// FilterParamsCheck find the min and max values for each filter value
	// FilterParamsCheck(Artists)

	tmpl := template.Must(template.ParseFiles("assets/templates/index.html"))
	if err := tmpl.Execute(w, Artists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(time.Since(start))
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := GetID(r.URL.Query().Get("id")) - 1
	artist := ArtistRender{
		ID:             Artists[id].ID,
		Image:          Artists[id].Image,
		Name:           Artists[id].Name,
		Members:        Artists[id].Members,
		CreationDate:   Artists[id].CreationDate,
		FirstAlbum:     Artists[id].FirstAlbum,
		DatesLocations: relations.Index[id].DatesLocations,
	}

	artist.MapDetails.Locations = Geocoding(artist)

	artist.MapDetails.MapURL = CreatMap(artist)
	// fmt.Println("From artistHanderl (fullUrl)\n", artist.MapDetails.mapURL)
	// for _, loc := range artist.MapDetails.Locations {
	// 	fmt.Println(loc.Location + ":")
	// 	fmt.Printf("Lat: (%v) Lng: (%v)\n", loc.Lat, loc.Lng)
	// }

	// Geo := GeoLocation{}
	// _, Geo.Lat, Geo.Lng = Geocoding(artist)
	// artist.Location = append(artist.Location, Geo)
	// fmt.Printf("Latitude: (%f) Longtitude: (%f)\n", artist.Location[0].Lat, artist.Location[0].Lng)

	tmpl := template.Must(template.New("artist.html").Funcs(template.FuncMap{
		"FormatText": FormatText,
	}).ParseFiles("assets/templates/artist.html"))
	if err := tmpl.Execute(w, artist); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func Geocoding(artist ArtistRender) []GeoLocation {
	var (
		GeoLocations []GeoLocation
		APIKey       = "AIzaSyD0AeIdWqfSMZujmXzaOHAQx1deLoYFnFE"
		baseURL      = "https://maps.googleapis.com/maps/api/geocode/json?"
	)

	for loc := range artist.DatesLocations {
		var (
			respData        GeocodingResponse
			locationAddress GeoLocation
			address         = strings.ReplaceAll(loc, "-", ",")
		)

		params := url.Values{}
		params.Add("address", address)
		params.Add("key", APIKey)

		resp, err := http.Get(baseURL + params.Encode())
		if err != nil {
			fmt.Println("Error making request")
			return nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response", err)
			return nil
		}

		err = json.Unmarshal(body, &respData)
		if err != nil {
			fmt.Println("failed to unmarshall data", err)
		}

		// fmt.Println("Address", address)
		// fmt.Println("Formatted Address", respData.Results[0].FormattedAddress)
		// fmt.Printf("Lat: (%v) Lng: (%v)\n", respData.Results[0].Geometry.Location.Lat, respData.Results[0].Geometry.Location.Lng)

		locationAddress.Location = respData.Results[0].FormattedAddress
		locationAddress.Lat, locationAddress.Lng = respData.Results[0].Geometry.Location.Lat, respData.Results[0].Geometry.Location.Lng

		GeoLocations = append(GeoLocations, locationAddress)
	}

	return GeoLocations
}

func CreatMap(artist ArtistRender) string {
	var (
		baseURL string = "https://maps.googleapis.com/maps/api/staticmap?"
		APIKey  string = "AIzaSyD0AeIdWqfSMZujmXzaOHAQx1deLoYFnFE"

		// mapCenter string = "41.0082,28.9784"
		// zoom    string = "2"
		size    string = "500x400"
		markers string
	)
	fmt.Println("Markers:\n", markers)

	params := url.Values{}
	// params.Add("center", mapCenter)
	params.Add("size", size)
	// params.Add("zoom", zoom)
	for _, loc := range artist.MapDetails.Locations {
		marker := fmt.Sprintf("%v,%v", loc.Lat, loc.Lng)
		params.Add("markers", marker)
	}
	params.Add("key", APIKey)

	return baseURL + params.Encode()
}
