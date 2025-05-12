package handlers

import (
	"encoding/json"
	"groopie_local/models"
	"groopie_local/services"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Suggestion represents a single autocomplete suggestion, including its display value and the original input.
type Suggestion struct {
	Value string `json:"value"` // The formatted suggestion string (e.g., "Artist Name - artist").
	Input string `json:"input"` // The original input string that triggered this suggestion.
}

// SearchHandler handles search requests for autocomplete functionality.
// It reads query parameters, fetches cached artist data, generates suggestions based on the search type,
// and responds with a JSON array of suggestions.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the search query (parameter "q") and the search type (parameter "searchType").
	query := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("q")))
	searchType := r.URL.Query().Get("searchType")

	// If no query is provided, redirect the user to the home page.
	if query == "" {
		http.Redirect(w, r, "/home", http.StatusSeeOther) // 303 See Other redirect.
		return
	}

	// Default the search type to "general" if not provided.
	if searchType == "" {
		searchType = "general"
	}

	// Normalize the search query to lowercase for case-insensitive matching.
	query = strings.ToLower(query)

	// Fetch the cached artist data.
	artistsFull, err := services.GetCachedData()
	if err != nil {
		log.Printf("Error fetching cached data: %v", err)
		http.Error(w, "Unable to fetch data", http.StatusInternalServerError)
		return
	}

	// Generate autocomplete suggestions based on the search query and type.
	suggestions := generateSuggestions(artistsFull, query, searchType)

	// Set the response content type to JSON.
	w.Header().Set("Content-Type", "application/json")
	// Encode the suggestions slice into JSON and write it to the response.
	if err := json.NewEncoder(w).Encode(suggestions); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// generateSuggestions creates a list of autocomplete suggestions based on the provided artist data,
// the search query, and the search type (e.g., "name", "location", "date", or "general").
// It returns a slice of Suggestion objects.
func generateSuggestions(artistsFull []models.ArtistFull, query, searchType string) []Suggestion {
	var suggestions []Suggestion
	// uniqueSuggestions ensures that duplicate suggestions are not added.
	uniqueSuggestions := make(map[string]bool)

	// Normalize the query to lowercase for consistent matching.
	query = strings.ToLower(query)
	// Determine if the query is numeric; this helps in filtering date-related suggestions.
	isNumeric := isNumber(query)

	// Iterate over each artist in the full list.
	for _, artist := range artistsFull {
		switch searchType {
		case "name":
			// If the artist's name contains the query, add a suggestion.
			if strings.Contains(strings.ToLower(artist.Artist.Name), query) {
				addSuggestion(&suggestions, uniqueSuggestions, artist.Artist.Name, "artist")
			}
		case "location":
			// Iterate over each location in the artist's performance data.
			for location := range artist.Relations.DatesLocations {
				if strings.Contains(strings.ToLower(location), query) {
					addSuggestion(&suggestions, uniqueSuggestions, location, "venue for "+artist.Artist.Name)
				}
			}
		case "date":
			// For date searches, only proceed if the query is numeric.
			if isNumeric {
				// Iterate over each set of dates in the artist's performance data.
				for _, dates := range artist.Relations.DatesLocations {
					for _, date := range dates {
						if strings.Contains(date, query) {
							addSuggestion(&suggestions, uniqueSuggestions, date, "concert date for "+artist.Artist.Name)
						}
					}
				}
			}
		case "general":
			// General search: check multiple fields.
			// Match artist name.
			if strings.Contains(strings.ToLower(artist.Artist.Name), query) {
				addSuggestion(&suggestions, uniqueSuggestions, artist.Artist.Name, "artist")
			}

			// Match band members.
			for _, member := range artist.Artist.Members {
				if strings.Contains(strings.ToLower(member), query) {
					addSuggestion(&suggestions, uniqueSuggestions, member, "member of "+artist.Artist.Name)
				}
			}

			// Match first album.
			if strings.Contains(strings.ToLower(artist.Artist.FirstAlbum), query) {
				addSuggestion(&suggestions, uniqueSuggestions, artist.Artist.FirstAlbum, "first album by "+artist.Artist.Name)
			}

			// Match creation year.
			creationYear := strconv.Itoa(artist.Artist.CreationDate)
			if strings.Contains(creationYear, query) {
				addSuggestion(&suggestions, uniqueSuggestions, creationYear, "creation year of "+artist.Artist.Name)
			}

			// Match concert dates.
			if isNumeric {
				for _, dates := range artist.Relations.DatesLocations {
					for _, date := range dates {
						if strings.Contains(date, query) {
							addSuggestion(&suggestions, uniqueSuggestions, date, "concert date for "+artist.Artist.Name)
						}
					}
				}
			}

			// Match locations.
			for location := range artist.Relations.DatesLocations {
				if strings.Contains(strings.ToLower(location), query) {
					addSuggestion(&suggestions, uniqueSuggestions, location, "venue for "+artist.Artist.Name)
				}
			}
		}
	}

	return suggestions
}

// addSuggestion constructs a suggestion string by combining the input and a descriptive tag.
// It ensures the suggestion is unique before adding it to the suggestions slice.
func addSuggestion(suggestions *[]Suggestion, uniqueSuggestions map[string]bool, input string, tag string) {
	// Format the suggestion as "input - tag".
	suggestion := input + " - " + tag
	// Add the suggestion only if it hasn't already been added.
	if !uniqueSuggestions[suggestion] {
		*suggestions = append(*suggestions, Suggestion{Value: suggestion, Input: input})
		uniqueSuggestions[suggestion] = true
	}
}

// isNumber checks if the given string consists solely of numeric characters.
// It returns true if the string is a number; otherwise, false.
func isNumber(input string) bool {
	for _, char := range input {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
