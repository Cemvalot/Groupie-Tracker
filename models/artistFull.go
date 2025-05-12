package models

type ArtistFull struct {
	Artist    Artist    `json:"artist"`
	Location  Location  `json:"location"`
	Relations Relations `json:"relations"`
	Dates     Date      `json:"date"`
}
