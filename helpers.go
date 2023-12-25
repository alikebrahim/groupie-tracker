package main

import (
	"fmt"
	"net/http"
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
		ID:             Artists[i].ID,
		Image:          Artists[i].Image,
		Name:           Artists[i].Name,
		Members:        Artists[i].Members,
		CreationDate:   Artists[i].CreationDate,
		FirstAlbum:     Artists[i].FirstAlbum,
		DatesLocations: relations.Index[i].DatesLocations,
	}
}

// GetInfo returns the formatted first album year && a bool if location exists in the artist's locations
func GetInfo(id int, loc string, a Artist) (int, bool) {
	var locationExist bool
	var locationsList []string
	dateParsed, err := time.Parse("02-01-2006", a.FirstAlbum)
	if err != nil {
		fmt.Println(err) // server side error #500
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
	return firstAlbumYear, locationExist
}

func GetInfoFiltered(id int, loc string, a ArtistRender) (int, bool) {
	fmt.Println("Start for: ", a.Name)
	var locationExist bool
	var locationsList []string
	dateParsed, err := time.Parse("02-01-2006", a.FirstAlbum)
	if err != nil {
		fmt.Println(err) // server side error #500
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
	fmt.Println("End for: ", a.Name)
	return firstAlbumYear, locationExist
}

func FormatText(input string) string {
	input = strings.ToTitle(input)

	input = strings.ReplaceAll(input, "_", " ")

	input = strings.ReplaceAll(input, "-", ", ")

	return input
}

func FormParse(r *http.Request) (int, int, int, int, string, []int) {
	creationDateFrom := r.PostFormValue("yearRangeFrom")
	creationYearFrom, err := strconv.Atoi(creationDateFrom)
	if err != nil {
		fmt.Println(err) // server side error #500
	}

	creationDateTo := r.PostFormValue("yearRangeTo")
	creationYearTo, err := strconv.Atoi(creationDateTo)
	if err != nil {
		fmt.Println(err) // server side error #500
	}

	firstAlbumDateFrom := r.PostFormValue("firstAlbumFrom")
	firstAlbumYearFrom, err := strconv.Atoi(firstAlbumDateFrom)
	if err != nil {
		fmt.Println(err) // server side error #500
	}

	firstAlbumDateTo := r.PostFormValue("firstAlbumTo")
	firstAlbumYearTo, err := strconv.Atoi(firstAlbumDateTo)
	if err != nil {
		fmt.Println(err) // server side error #500
	}

	location := r.PostFormValue("location")

	r.ParseForm()
	membersNumStr := r.Form["membersNum[]"]
	var memmbersNum []int
	for _, str := range membersNumStr {

		int, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println(err) // server side error #500
			continue
		}
		memmbersNum = append(memmbersNum, int)
	}

	return creationYearFrom, creationYearTo, firstAlbumYearFrom, firstAlbumYearTo, location, memmbersNum
}

// below code was used to find the values of highest number of members of a group, earliest album and creation dates
// and latest album and creation dates
func FilterParamsCheck(Artists []Artist) {
	var (
		maxMem, earliestAlbum, latestAlbum, earliestCreation, latestCreation int
	)
	for indx, artist := range Artists {
		albumYearParsed, err := time.Parse("02-01-2006", artist.FirstAlbum)
		if err != nil {
			fmt.Println(err) // server side error #500
		}
		firstAlbumYear := albumYearParsed.Year()
		if err != nil {
			fmt.Println(err) // server side error #500
		}
		if indx == 0 {
			earliestAlbum = firstAlbumYear
			earliestCreation = artist.CreationDate
		}
		if len(artist.Members) > maxMem {
			maxMem = len(artist.Members)
		}
		if artist.CreationDate < earliestCreation {
			earliestCreation = artist.CreationDate
		}
		if artist.CreationDate > latestCreation {
			latestCreation = artist.CreationDate
		}
		if firstAlbumYear < earliestAlbum {
			earliestAlbum = firstAlbumYear
		}
		if firstAlbumYear > latestAlbum {
			latestAlbum = firstAlbumYear
		}

	}
	fmt.Printf("max members: %d <<>> earliest album year: %d <<>> latest album year: %d <<>> earliest creation date: %d <<>> latest creation date: %d \n", maxMem, earliestAlbum, latestAlbum, earliestCreation, latestCreation)

}
