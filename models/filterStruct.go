package models

type Filters struct {
	SearchQuery              string
	SearchType               string
	CreationMin, CreationMax int
	AlbumMin, AlbumMax       int
	BandMembers              []string
	PerformanceLocations     []string
}
