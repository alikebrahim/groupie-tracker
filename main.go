package main

import (
	"fmt"
	"html/template"
	"net/http"
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
		artistHanlder(w, r)
	case "POST":
		filterHandler(w, r)
	}
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: rewrite the code to check for which filters have been set and then filter gradually (e.g using creation date first, then filter the list into another with number members)
	// TODO: fix bug where previous filter results are added to new filter when reloading the website
	var filteredArtists []ArtistRender
	creationYearFrom, creationYearTo, firstAlbumYearFrom, firstAlbumYearTo, location, membersNum := FormParse(r)
	fmt.Println("***** Start - filterHandler: FormParse ******")
	fmt.Println("creation date from: ", creationYearFrom)
	fmt.Println("creation date to: ", creationYearTo)
	fmt.Println("first album from: ", firstAlbumYearFrom)
	fmt.Println("first album to: ", firstAlbumYearTo)
	fmt.Println("location: ", location)
	fmt.Println("members: ", membersNum)
	fmt.Println("***** End - filterHandler: FormParse ******")

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

	// Other filters

	// for id, artist := range Artists {
	// 	location = strings.Split(location, ",")[0]
	// 	artistFirstAlbumYear, loc := GetInfo(id, strings.ToLower(location), artist)
	// 	if (artist.CreationDate >= creationYearFrom && artist.CreationDate <= creationYearTo) || (artistFirstAlbumYear >= firstAlbumYearFrom && artistFirstAlbumYear <= firstAlbumYearTo) || loc || len(artist.Members) == membersNum {
	// 		// fmt.Printf("Artist:%v Creation Date:%v, First Album:%v \n", artist.Name, artist.CreationDate, artist.FirstAlbum)
	// 		artistMatch := ArtistRender{
	// 			ID:             Artists[id].ID,
	// 			Image:          Artists[id].Image,
	// 			Name:           Artists[id].Name,
	// 			Members:        Artists[id].Members,
	// 			CreationDate:   Artists[id].CreationDate,
	// 			FirstAlbum:     Artists[id].FirstAlbum,
	// 			DatesLocations: relations.Index[id].DatesLocations,
	// 		}
	// 		filteredArtists = append(filteredArtists, artistMatch)
	// 	}
	// }
	tmpl := template.Must(template.ParseFiles("assets/templates/filter.html"))
	if err := tmpl.Execute(w, filteredArtists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func artistHanlder(w http.ResponseWriter, r *http.Request) {
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
	// fmt.Println(id)
	artist := ArtistRender{
		ID:             Artists[id].ID,
		Image:          Artists[id].Image,
		Name:           Artists[id].Name,
		Members:        Artists[id].Members,
		CreationDate:   Artists[id].CreationDate,
		FirstAlbum:     Artists[id].FirstAlbum,
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
