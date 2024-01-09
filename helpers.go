package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func GetID(id string) int {
	artist, err := strconv.Atoi(id)
	if err != nil {
		return 0
	}
	return artist
}

func MakeArtistRender(i int) ArtistRender {
	return ArtistRender{
		ID:               Artists[i].ID,
		Image:            Artists[i].Image,
		Name:             Artists[i].Name,
		Members:          Artists[i].Members,
		CreationDate:     Artists[i].CreationDate,
		FirstAlbum:       Artists[i].FirstAlbum,
		DatesLocations:   relations.Index[i].DatesLocations,
		DatesLocations_F: formatLocations(relations.Index[i].DatesLocations),
	}
}

// GetInfo returns the formatted first album year && a bool if location exists in the artist's locations
func GetInfo(id int, loc string, a Artist) (int, bool, error) {
	var locationExist bool
	var locationsList []string
	dateParsed, err := time.Parse("02-01-2006", a.FirstAlbum)
	if err != nil {
		return 0, false, err
	}
	firstAlbumYear := dateParsed.Year()

	locationDates := relations.Index[id].DatesLocations

	for k := range locationDates {
		city := strings.Split(k, "-")[0]
		cityformatted := strings.Replace(city, "_", " ", -1)
		locationsList = append(locationsList, cityformatted)
	}
	for _, l := range locationsList {
		if l == loc {
			locationExist = true
			break
		}
	}
	return firstAlbumYear, locationExist, nil
}

func GetInfoFiltered(id int, loc string, a ArtistRender) (int, bool, error) {
	var locationExist bool
	var locationsList []string
	dateParsed, err := time.Parse("02-01-2006", a.FirstAlbum)
	if err != nil {
		return 0, false, err
	}
	firstAlbumYear := dateParsed.Year()

	locationDates := relations.Index[id].DatesLocations
	for k := range locationDates {
		city := strings.Split(k, "-")[0]
		cityformatted := strings.Replace(city, "_", " ", -1)
		locationsList = append(locationsList, cityformatted)
	}
	for _, l := range locationsList {
		if l == loc {
			locationExist = true
			break
		}
	}
	return firstAlbumYear, locationExist, nil
}

func formatLocations(dl map[string][]string) map[string][]string {
	newDL := make(map[string][]string)

	for k, v := range dl {
		input := strings.ToTitle(k)

		input = strings.ReplaceAll(input, "_", " ")

		input = strings.ReplaceAll(input, "-", ", ")

		newDL[input] = v
	}

	return newDL
}

func FormParse(r *http.Request) (int, int, int, int, string, []int, error) {
	creationDateFrom := r.PostFormValue("yearRangeFrom")
	creationYearFrom, err := strconv.Atoi(creationDateFrom)
	if err != nil {
		return 0, 0, 0, 0, "", []int{}, err
	}

	creationDateTo := r.PostFormValue("yearRangeTo")
	creationYearTo, err := strconv.Atoi(creationDateTo)
	if err != nil {
		return 0, 0, 0, 0, "", []int{}, err
	}

	firstAlbumDateFrom := r.PostFormValue("firstAlbumFrom")
	firstAlbumYearFrom, err := strconv.Atoi(firstAlbumDateFrom)
	if err != nil {
		return 0, 0, 0, 0, "", []int{}, err
	}

	firstAlbumDateTo := r.PostFormValue("firstAlbumTo")
	firstAlbumYearTo, err := strconv.Atoi(firstAlbumDateTo)
	if err != nil {
		return 0, 0, 0, 0, "", []int{}, err
	}

	location := r.PostFormValue("location")

	r.ParseForm()
	membersNumStr := r.Form["membersNum[]"]
	var memmbersNum []int
	for _, str := range membersNumStr {

		int, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println(err)
			continue
		}
		memmbersNum = append(memmbersNum, int)
	}

	return creationYearFrom, creationYearTo, firstAlbumYearFrom, firstAlbumYearTo, location, memmbersNum, nil
}

func containsMember(members []string, query string) bool {
	for _, m := range members {
		if strings.Contains(strings.ToLower(m), strings.ToLower(query)) {
			return true
		}
	}
	return false
}

func checkFilteredArtists(arr []Artist, artist Artist) bool {
	if len(arr) == 0 {
		return false
	}
	id := artist.ID
	for i := 0; i < len(arr); i++ {
		if arr[i].ID == id {
			return true
		}

	}
	return false

}

func Geocoding(artist ArtistRender) ([]GeoLocation, error) {
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
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &respData)
		if err != nil {
			return nil, err
		}

		locationAddress.Location = respData.Results[0].FormattedAddress
		locationAddress.Lat, locationAddress.Lng = respData.Results[0].Geometry.Location.Lat, respData.Results[0].Geometry.Location.Lng

		GeoLocations = append(GeoLocations, locationAddress)
	}

	return GeoLocations, nil
}

func CreatMap(artist ArtistRender) string {
	var (
		baseURL string = "https://maps.googleapis.com/maps/api/staticmap?"
		APIKey  string = "AIzaSyD0AeIdWqfSMZujmXzaOHAQx1deLoYFnFE"

		size string = "500x400"
	)

	params := url.Values{}
	params.Add("size", size)
	for _, loc := range artist.MapDetails.Locations {
		marker := fmt.Sprintf("%v,%v", loc.Lat, loc.Lng)
		params.Add("markers", marker)
	}
	params.Add("key", APIKey)

	return baseURL + params.Encode()
}

func RenderPage(w http.ResponseWriter, status int, page string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	err := Tmpl.ExecuteTemplate(w, page, data)
	if err != nil {
		if page != "404.html" {
			RenderPage(w, http.StatusNotFound, "404.html", data)
		}
	}
}
