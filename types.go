package main

import "html/template"

var Tmpl *template.Template

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
	// DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistRender struct {
	ID               int                 `json:"id"`
	Image            string              `json:"image"`
	Name             string              `json:"name"`
	Members          []string            `json:"members"`
	CreationDate     int                 `json:"creationDate"`
	FirstAlbum       string              `json:"firstAlbum"`
	DatesLocations   map[string][]string `json:"datesLocations"`
	DatesLocations_F map[string][]string `json:"datesLocations_f"`
	MapDetails       MapLocations        `json:"map_details"`
}

type MapLocations struct {
	MapURL    string        `json:"map_url"`
	Locations []GeoLocation `json:"locations"`
}

type GeoLocation struct {
	Location string  `json:"location"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

// locations API
type Locations struct {
	Index []LocationsIndex `json:"index"`
}

type LocationsIndex struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

// Relations API
type RelationsIndex struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Relations struct {
	Index []RelationsIndex `json:"index"`
}

// I forgot API
type DatesIndex struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Dates struct {
	Index []DatesIndex `json:"index"`
}

// Geocoding API
type GeocodingResponse struct {
	Results []Result `json:"results"`
	Status  string   `json:"status"`
}

type Result struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          Geometry           `json:"geometry"`
	PlaceID           string             `json:"place_id"`
	Types             []string           `json:"types"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Geometry struct {
	Bounds       Bounds     `json:"bounds"`
	Location     Coordinate `json:"location"`
	LocationType string     `json:"location_type"`
	Viewport     Bounds     `json:"viewport"`
}

type Bounds struct {
	Northeast Coordinate `json:"northeast"`
	Southwest Coordinate `json:"southwest"`
}

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
