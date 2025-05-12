package models

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"` // Updated to match API response
}

// IsValid checks if the location has a valid ID and at least one location.
func (l *Location) IsValid() bool {
	return l.ID > 0 && len(l.Locations) > 0
}
