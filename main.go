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

func init() {
	Tmpl = template.Must(template.ParseGlob("assets/templates/*.html"))
}

func main() {
	http.HandleFunc("/", indexRouter)
	http.HandleFunc("/search", searchHandler)
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		RenderPage(w, http.StatusNotFound, "404.html", struct {
			error string
		}{
			error: "404 Not Found",
		})
		return
	}
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

	RenderPage(w, http.StatusOK, "index.html", artists)
	fmt.Println(time.Since(start))
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := GetID(r.URL.Query().Get("id")) - 1

	artist := MakeArtistRender(id)

	data, err := Geocoding(artist)
	if err != nil {
		RenderPage(w, http.StatusInternalServerError, "500.html", struct {
			error string
		}{
			error: "500 error",
		})
		return
	}
	artist.MapDetails.Locations = data

	artist.MapDetails.MapURL = CreatMap(artist)

	RenderPage(w, http.StatusOK, "artist.html", artist)
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

	RenderPage(w, http.StatusOK, "index.html", filteredArtists)
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	var filteredArtists []ArtistRender
	creationYearFrom, creationYearTo, firstAlbumYearFrom, firstAlbumYearTo, location, membersNum, err := FormParse(r)
	if err != nil {
		RenderPage(w, http.StatusInternalServerError, "500.html", struct {
			error string
		}{
			error: "500 error",
		})
		return
	}

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
				artistFirstAlbumYear, _, err := GetInfo(i, strings.ToLower(location), Artists[filteredArtists[i].ID-1])
				if err != nil {
					RenderPage(w, http.StatusInternalServerError, "500.html", struct {
						error string
					}{
						error: "500 error",
					})
					return
				}
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
				_, loc, err := GetInfoFiltered(filteredArtists[i].ID-1, strings.ToLower(location), filteredArtists[i])
				if err != nil {
					RenderPage(w, http.StatusInternalServerError, "500.html", struct {
						error string
					}{
						error: "500 error",
					})
					return
				}
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
			artistFirstAlbumYear, _, err := GetInfo(id, strings.ToLower(location), artist)
			if err != nil {
				RenderPage(w, http.StatusInternalServerError, "500.html", struct {
					error string
				}{
					error: "500 error",
				})
				return
			}
			if artistFirstAlbumYear >= firstAlbumYearFrom && artistFirstAlbumYear <= firstAlbumYearTo {
				artistMatch := MakeArtistRender(id)
				filteredArtists = append(filteredArtists, artistMatch)
			}
		}
		if location != "" {
			var newfilteredArtists []ArtistRender
			for i := range filteredArtists {
				location = strings.Split(location, ",")[0]
				_, loc, err := GetInfoFiltered(filteredArtists[i].ID-1, strings.ToLower(location), filteredArtists[i])
				if err != nil {
					RenderPage(w, http.StatusInternalServerError, "500.html", struct {
						error string
					}{
						error: "500 error",
					})
					return
				}
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
			_, loc, err := GetInfo(id, strings.ToLower(location), artist)
			if err != nil {
				RenderPage(w, http.StatusInternalServerError, "500.html", struct {
					error string
				}{
					error: "500 error",
				})
				return
			}
			if loc {
				artistMatch := MakeArtistRender(id)
				filteredArtists = append(filteredArtists, artistMatch)
			}
		}
		if firstAlbumYearFrom >= 1963 {
			var newfilteredArtists []ArtistRender
			for i := range filteredArtists {
				artistFirstAlbumYear, _, err := GetInfo(i, strings.ToLower(location), Artists[filteredArtists[i].ID-1])
				if err != nil {
					RenderPage(w, http.StatusInternalServerError, "500.html", struct {
						error string
					}{
						error: "500 error",
					})
					return
				}
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
				_, loc, err := GetInfoFiltered(filteredArtists[i].ID-1, strings.ToLower(location), filteredArtists[i])
				if err != nil {
					RenderPage(w, http.StatusInternalServerError, "500.html", struct {
						error string
					}{
						error: "500 error",
					})
					return
				}
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

	RenderPage(w, http.StatusOK, "index.html", filteredArtists)
}
