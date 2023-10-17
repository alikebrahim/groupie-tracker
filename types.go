package main

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
}
type Locations struct {
	Index []LocationsIndex `json:"index"`
}
type Dates struct {
	Index []DatesIndex `json:"index"`
}
type Relations struct {
	Index []RelationsIndex `json:"index"`
}
type LocationsIndex struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}
type DatesIndex struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type RelationsIndex struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}