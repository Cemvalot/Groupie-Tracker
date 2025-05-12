package handlers

import (
	"fmt"
	"groopie_local/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

const templateDir = "templates"

// TemplateData holds the data passed to templates.
// It includes information such as the page title, list of artists,
// details for a specific artist, search parameters, and any custom messages.
type TemplateData struct {
	Title       string
	Artists     []models.ArtistFull
	Artist      models.ArtistFull
	SearchQuery string
	SearchType  string
	SortBy      string
	Message     string
}

// renderTemplate renders a specified HTML template along with an additional filter modal template.
// It constructs the template file paths, parses them, and executes the template with the provided data.
// If parsing or execution fails, it logs the error and sends an HTTP 500 Internal Server Error response.
func renderTemplate(w http.ResponseWriter, name string, data TemplateData) {
	// Construct full file paths for the main template and the filter modal.
	tmplPath := filepath.Join(templateDir, name+".html")
	filterPath := filepath.Join(templateDir, "filterModal.html")

	// Parse the main template file and the filter modal.
	tmpl, err := template.ParseFiles(tmplPath, filterPath)
	if err != nil {
		log.Printf("Error parsing template (%s): %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the parsed template using the provided data.
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template (%s): %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// filterByType filters the provided slice of artists based on the search query and type.
// Depending on the searchType, it filters by artist name, location, date, or a general search across multiple fields.
func FilterByType(artists []models.ArtistFull, query, searchType string) []models.ArtistFull {
	query = strings.ToLower(query)
	var filtered []models.ArtistFull

	for _, artist := range artists {
		nameLower := strings.ToLower(artist.Artist.Name)
		firstAlbumLower := strings.ToLower(artist.Artist.FirstAlbum)
		creationYear := strconv.Itoa(artist.Artist.CreationDate)

		switch searchType {
		case "name":
			// Match artist name or any of the band members.
			if strings.Contains(nameLower, query) ||
				containsAnyLower(artist.Artist.Members, query) {
				filtered = append(filtered, artist)
			}
		case "location":
			// Match any of the artist's event locations.
			if containsAnyLower(artist.Location.Locations, query) {
				filtered = append(filtered, artist)
			}
		case "date":
			// Match any of the artist's concert dates.
			if containsAnyLower(artist.Dates.Dates, query) {
				filtered = append(filtered, artist)
			}
		case "general":
			// General search matches against name, band members, first album,
			// creation year, concert dates, or locations.
			if strings.Contains(nameLower, query) ||
				containsAnyLower(artist.Artist.Members, query) ||
				strings.Contains(firstAlbumLower, query) ||
				strings.Contains(creationYear, query) ||
				containsAnyLower(artist.Dates.Dates, query) ||
				containsAnyLower(artist.Location.Locations, query) {
				filtered = append(filtered, artist)
			}
		default:
			// Default filtering by artist name.
			if strings.Contains(nameLower, query) {
				filtered = append(filtered, artist)
			}
		}
	}

	return filtered
}

// containsAnyLower checks if any string in the given slice contains the query (case insensitive).
// Returns true if a match is found; otherwise, false.
func containsAnyLower(items []string, query string) bool {
	for _, item := range items {
		if strings.Contains(strings.ToLower(item), query) {
			return true
		}
	}
	return false
}

// filterByBandMembers filters the artists based on the number of band members.
// The bandMembers slice contains the desired number(s) of members as strings.
// If no filter is applied (empty slice), all artists are returned.
func FilterByBandMembers(artists []models.ArtistFull, bandMembers []string) []models.ArtistFull {
	var filtered []models.ArtistFull
	for _, artist := range artists {
		// Determine the number of band members for the artist.
		numMembers := len(artist.Artist.Members)
		// Include the artist if no filter is set or if the member count matches one of the filter values.
		if len(bandMembers) == 0 || contains(bandMembers, strconv.Itoa(numMembers)) {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

// contains checks if a slice contains the specified string value.
// Returns true if found, false otherwise.
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// filterByCreationDateFull filters the artists based on a range of creation years.
// Only artists whose creation year falls between minYear and maxYear (inclusive) are returned.
func FilterByCreationDateFull(artists []models.ArtistFull, minYear, maxYear int) []models.ArtistFull {
	var filtered []models.ArtistFull
	for _, artistFull := range artists {
		creationDate := artistFull.Artist.CreationDate
		if creationDate >= minYear && creationDate <= maxYear {
			filtered = append(filtered, artistFull)
		}
	}
	return filtered
}

// filterByFirstAlbumDateFull filters the artists based on the year of their first album release.
// It extracts the year from the album's date string and returns only those artists whose release year is within the specified range.
func FilterByFirstAlbumDateFull(artists []models.ArtistFull, minYear, maxYear int) []models.ArtistFull {
	var filtered []models.ArtistFull
	for _, artistFull := range artists {
		year, err := extractYear(artistFull.Artist.FirstAlbum)
		if err != nil {
			log.Printf("Error extracting year from first album date for artist %s: %v", artistFull.Artist.Name, err)
			continue
		}
		if year >= minYear && year <= maxYear {
			filtered = append(filtered, artistFull)
		}
	}
	return filtered
}

// extractYear extracts the year from a date string formatted as "DD-MM-YYYY".
// Returns the year as an integer or an error if the format is invalid.
func extractYear(dateStr string) (int, error) {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid date format: %s", dateStr)
	}
	year := parts[2]
	return strconv.Atoi(year)
}

// filterByPerformanceLocations filters artists based on provided performance locations.
// For each artist, it checks if all specified locations match at least one of the artist's performance locations.
// If no locations are provided, it returns all artists.
func FilterByPerformanceLocations(artists []models.ArtistFull, locations []string) []models.ArtistFull {
	if len(locations) == 0 {
		return artists // Return all artists if no filter locations are provided.
	}

	var filtered []models.ArtistFull

	// Iterate through each artist.
	for _, artist := range artists {
		matchesAll := true // Assume artist matches all filter criteria initially.

		// Check each filter location.
		for _, filterLoc := range locations {
			normalizedFilterLoc := strings.ToLower(strings.TrimSpace(filterLoc))
			locationFound := false

			// Iterate through the artist's performance locations.
			for artistLoc := range artist.Relations.DatesLocations {
				normalizedArtistLoc := strings.ToLower(strings.TrimSpace(artistLoc))

				// If filter location contains a comma, expect format "City/State, Country".
				if strings.Contains(normalizedFilterLoc, ",") {
					parts := strings.Split(normalizedFilterLoc, ",")
					if len(parts) != 2 {
						continue
					}

					cityState := strings.TrimSpace(parts[0])
					country := strings.TrimSpace(parts[1])

					// Check if both city/state and country are present.
					if strings.Contains(normalizedArtistLoc, cityState) &&
						strings.Contains(normalizedArtistLoc, country) {
						locationFound = true
						break
					}
				} else {
					// Simple case: check if the artist's location contains the filter location.
					if strings.Contains(normalizedArtistLoc, normalizedFilterLoc) {
						locationFound = true
						break
					}
				}
			}

			// If a filter location does not match any of the artist's locations, exclude the artist.
			if !locationFound {
				matchesAll = false
				break
			}
		}

		if matchesAll {
			filtered = append(filtered, artist)
		}
	}

	return filtered
}
