package handlers

import (
	"groopie_local/models"
	"net/http"
	"strconv"
)

func ParseFilters(r *http.Request) models.Filters {
	q := r.URL.Query()

	return models.Filters{
		SearchQuery:          q.Get("search"),
		SearchType:           q.Get("searchType"),
		CreationMin:          parseInt(q.Get("creationMin"), DefaultCreationYearMin),
		CreationMax:          parseInt(q.Get("creationMax"), DefaultCreationYearMax),
		AlbumMin:             parseInt(q.Get("albumMin"), DefaultAlbumYearMin),
		AlbumMax:             parseInt(q.Get("albumMax"), DefaultAlbumYearMax),
		BandMembers:          q["bandMembers"],
		PerformanceLocations: q["locations"],
	}
}

func parseInt(val string, fallback int) int {
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	return fallback
}

func FilterArtists(artists []models.ArtistFull, filters models.Filters) []models.ArtistFull {
	filtered := artists

	if filters.SearchQuery != "" {
		filtered = FilterByType(filtered, filters.SearchQuery, filters.SearchType)
	}
	if len(filters.BandMembers) > 0 {
		filtered = FilterByBandMembers(filtered, filters.BandMembers)
	}
	filtered = FilterByCreationDateFull(filtered, filters.CreationMin, filters.CreationMax)
	filtered = FilterByFirstAlbumDateFull(filtered, filters.AlbumMin, filters.AlbumMax)

	if len(filters.PerformanceLocations) > 0 {
		filtered = FilterByPerformanceLocations(filtered, filters.PerformanceLocations)
	}

	return filtered
}
