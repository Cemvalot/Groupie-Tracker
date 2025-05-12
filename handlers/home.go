package handlers

import (
	"groopie_local/services"
	"log"
	"net/http"
)

// WelcomeHandler serves the welcome page by sending the welcome.html file to the client.
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the welcome page file.
	http.ServeFile(w, r, "templates/welcome.html")
}

// NotFoundHandler serves a generic error page for undefined routes.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Serve an error page for any undefined route.
	http.ServeFile(w, r, "templates/error.html")
}

// HomeHandler processes requests for the home page, applying various filters on artist data.
// It retrieves cached data, applies search and filter parameters from the query string, and then renders the home template.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}

	artists, err := services.GetCachedData()
	if err != nil {
		log.Printf("Error fetching cached data: %v", err)
		renderTemplate(w, "error", TemplateData{
			Title:   "Error",
			Message: "Unable to load data. Please try again later.",
		})
		return
	}

	filters := ParseFilters(r)

	filteredArtists := FilterArtists(artists, filters)

	data := TemplateData{
		Title:       "Home - Groupie Tracker",
		Artists:     filteredArtists,
		SearchQuery: filters.SearchQuery,
		SearchType:  filters.SearchType,
	}

	renderTemplate(w, "home", data)
}

