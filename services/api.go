package services

import (
	"context"
	"encoding/json"
	"fmt"
	"groopie_local/models"
	"net/http"
	"sync"
	"time"
)

var (
	// cache stores the merged data for artists to avoid repeated API calls.
	cache []models.ArtistFull
	// cacheLock ensures thread-safe access to the cache.
	cacheLock sync.Mutex
	// lastCacheTime tracks the last time the cache was updated.
	lastCacheTime time.Time
)

// client is a custom HTTP client configured with connection pooling and timeouts.
// This improves performance by reusing connections and ensures API requests do not hang.
var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	},
}

// fetchFromAPI performs an HTTP GET request to the specified URL using a custom HTTP client,
// with context-based timeout to avoid long waits.
// It checks if the response status is 200 OK, and decodes the JSON response into the target interface.
// It returns an error if the request fails or if the status code is not OK.
func fetchFromAPI(url string, target interface{}) error {
	// Create a context with timeout to ensure the request is cancelled if it takes too long.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new HTTP request with the given context.
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	// Execute the request using the custom client.
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verify that the HTTP status is 200 OK.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: received status code %d from %s", resp.StatusCode, url)
	}

	// Decode the JSON response into the target interface.
	return json.NewDecoder(resp.Body).Decode(target)
}

// FetchArtists retrieves the list of artists from the external API.
// It returns a slice of models.Artist and an error if the request or decoding fails.
func FetchArtists() ([]models.Artist, error) {
	var artists []models.Artist
	err := fetchFromAPI("https://groupietrackers.herokuapp.com/api/artists", &artists)
	return artists, err
}

// FetchLocations retrieves location data from the external API.
// It returns a slice of models.Location and an error if the request or decoding fails.
func FetchLocations() ([]models.Location, error) {
	var response struct {
		Index []models.Location `json:"index"`
	}
	err := fetchFromAPI("https://groupietrackers.herokuapp.com/api/locations", &response)
	return response.Index, err
}

// FetchRelations retrieves relations data from the external API.
// It returns a slice of models.Relations and an error if the request or decoding fails.
func FetchRelations() ([]models.Relations, error) {
	var response struct {
		Index []models.Relations `json:"index"`
	}
	err := fetchFromAPI("https://groupietrackers.herokuapp.com/api/relation", &response)
	return response.Index, err
}

// FetchDates retrieves dates data from the external API.
// It returns a slice of models.Date and an error if the request or decoding fails.
func FetchDates() ([]models.Date, error) {
	var response struct {
		Index []models.Date `json:"index"`
	}
	err := fetchFromAPI("https://groupietrackers.herokuapp.com/api/dates", &response)
	return response.Index, err
}

// MergeData concurrently fetches artists, locations, relations, and dates data,
// then merges them into a slice of models.ArtistFull by matching their IDs.
// It returns the merged data or an error if any of the API calls fail.
func MergeData() ([]models.ArtistFull, error) {
	var (
		artists   []models.Artist
		locations []models.Location
		relations []models.Relations
		dates     []models.Date
	)

	var wg sync.WaitGroup
	// Buffered channel to capture errors from the four concurrent operations.
	errChan := make(chan error, 4)

	wg.Add(4)
	go func() {
		defer wg.Done()
		a, err := FetchArtists()
		if err != nil {
			errChan <- fmt.Errorf("FetchArtists: %w", err)
			return
		}
		artists = a
	}()
	go func() {
		defer wg.Done()
		l, err := FetchLocations()
		if err != nil {
			errChan <- fmt.Errorf("FetchLocations: %w", err)
			return
		}
		locations = l
	}()
	go func() {
		defer wg.Done()
		r, err := FetchRelations()
		if err != nil {
			errChan <- fmt.Errorf("FetchRelations: %w", err)
			return
		}
		relations = r
	}()
	go func() {
		defer wg.Done()
		d, err := FetchDates()
		if err != nil {
			errChan <- fmt.Errorf("FetchDates: %w", err)
			return
		}
		dates = d
	}()
	wg.Wait()
	close(errChan)

	// Return the first encountered error, if any.
	for e := range errChan {
		return nil, e
	}

	// Create maps for fast lookup of locations, relations, and dates by their IDs.
	locationMap := make(map[int]models.Location)
	for _, loc := range locations {
		locationMap[loc.ID] = loc
	}

	relationsMap := make(map[int]models.Relations)
	for _, rel := range relations {
		relationsMap[rel.ID] = rel
	}

	datesMap := make(map[int]models.Date)
	for _, d := range dates {
		datesMap[d.ID] = d
	}

	// Merge the data into a slice of models.ArtistFull.
	var artistsFull []models.ArtistFull
	for _, artist := range artists {
		artistsFull = append(artistsFull, models.ArtistFull{
			Artist:    artist,
			Location:  locationMap[artist.ID],
			Relations: relationsMap[artist.ID],
			Dates:     datesMap[artist.ID],
		})
	}

	return artistsFull, nil
}

// GetCachedData returns the cached merged artist data.
// If the cache is empty or older than 5 minutes, it refreshes the cache by calling MergeData.
// It uses a mutex to ensure thread safety.
func GetCachedData() ([]models.ArtistFull, error) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if len(cache) == 0 || time.Since(lastCacheTime) > 5*time.Minute {
		data, err := MergeData()
		if err != nil {
			return nil, err
		}
		cache = data
		lastCacheTime = time.Now()
	}
	return cache, nil
}
