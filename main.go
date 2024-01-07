package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// *************************************************************************//
// Note 1: activate the rest of the APIs before audit as they are required //
// Note 2: Ensure http status errors handling                             //
// **********************************************************************//

var (
	Artists   []Artist
	relations Relations

// locations Locations
// dates Dates
)

func main() {
	http.HandleFunc("/", indexRouter)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/artist/", artistHandler)
	http.HandleFunc("/filter-artist/", filterHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets"))))
	fmt.Println("Server starting at port 3000")
	http.ListenAndServe(":3000", nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var filteredArtists []Artist
	for _, artist := range Artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(query)) ||
			containsMember(artist.Members, query) || strings.Contains(strconv.Itoa(artist.CreationDate), strings.ToLower(query)) {
			if !checkFilteredArtists(filteredArtists, artist) {
				filteredArtists = append(filteredArtists, artist)
			}

		}
		var finalLocs []string
		locationDates := relations.Index[(artist.ID)-1].DatesLocations
		for loc := range locationDates {
			city := strings.Split(loc, "-")[0]
			country := strings.Split(loc, "-")[1]
			cityformatted := strings.Replace(city, "_", "", -1)
			countryformatted := strings.Replace(country, "_", "", -1)
			finalLocs = append(finalLocs, countryformatted)
			finalLocs = append(finalLocs, cityformatted)
		}

		for i := 0; i < len(finalLocs); i++ {
			if strings.Contains(strings.ToLower(finalLocs[i]), strings.ToLower(query)) {
				if !checkFilteredArtists(filteredArtists, artist) {
					filteredArtists = append(filteredArtists, artist)
				}
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("assets/templates/index.html"))
	if err := tmpl.Execute(w, filteredArtists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	tmpl := template.Must(template.ParseFiles("assets/templates/index.html"))
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
	var artists []ArtistRender
	for i := 0; i < len(Artists); i++ {
		artist := MakeArtistRender(i)
		artists = append(artists, artist)
	}

	// FilterParamsCheck find the min and max values for each filter value
	// FilterParamsCheck(Artists)

	tmpl := template.Must(template.ParseFiles("assets/templates/index.html"))
	if err := tmpl.Execute(w, artists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(time.Since(start))
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := GetID(r.URL.Query().Get("id")) - 1

	artist := MakeArtistRender(id)

	artist.MapDetails.Locations = Geocoding(artist)

	artist.MapDetails.MapURL = CreatMap(artist)

	tmpl := template.Must(template.ParseFiles("assets/templates/artist.html"))
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
		size string = "500x400"
		// markers string
	)

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
