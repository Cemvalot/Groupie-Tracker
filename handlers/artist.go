package handlers

import (
	"groopie_local/models"
	"groopie_local/services"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// ArtistHandler handles HTTP requests for a specific artist.
// It extracts the artist ID from the URL, fetches the corresponding artist data from the cache,
// and renders the artist page or an error page if the artist is not found.
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the "/artist/" prefix to extract the artist ID as a string.
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")

	// If the extracted string is empty or just "/", log the issue and serve an error page.
	if idStr == "" || idStr == "/" {
		log.Println("No artist ID provided")
		http.ServeFile(w, r, "templates/error.html") // Serve a generic error page.
		return
	}

	// Convert the artist ID from string to integer.
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 || id > 52 {
		log.Printf("Invalid artist ID: %s", idStr)
		http.ServeFile(w, r, "templates/error.html") // Serve error page for invalid ID.
		return
	}

	// Retrieve the merged and cached artist data.
	artistsFull, err := services.GetCachedData()
	if err != nil {
		log.Printf("Error fetching cached data: %v", err)
		// Render an error template using a helper function with a custom error title.
		renderTemplate(w, "error", TemplateData{Title: "Error"})
		return
	}

	// Build a map to quickly look up an artist by its ID.
	artistMap := make(map[int]models.ArtistFull)
	for _, artist := range artistsFull {
		artistMap[artist.Artist.ID] = artist
	}

	// Find the artist corresponding to the provided ID.
	artistFull, found := artistMap[id]
	if !found {
		log.Printf("Artist with ID %d not found", id)
		// If not found, respond with a 404 Not Found.
		http.NotFound(w, r)
		return
	}

	// Prepare data for the template, including a dynamic title.
	data := TemplateData{
		Title:  artistFull.Artist.Name + " - Groopie Tracker",
		Artist: artistFull,
	}

	// Render the artist template with the retrieved data.
	renderTemplate(w, "artist", data)
}
