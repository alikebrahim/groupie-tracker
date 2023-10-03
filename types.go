package main

type ArtistComplete struct {
	ID           int       `json:"id"`
	Image        string    `json:"image"`
	Name         string    `json:"name"`
	Members      []string  `json:"members"`
	CreationDate int       `json:"creationDate"`
	FirstAlbum   string    `json:"firstAlbum"`
	Locations    Locations `json:"locations"`
	ConcertDates Dates     `json:"concertDates"`
	Relations    string    `json:"relations"`
}
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
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}
type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type Relations struct {
	ID             int      `json:"id"`
	DatesLocations []string `json:"datesLocations"`
}
