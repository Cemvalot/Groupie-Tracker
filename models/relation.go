package models

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// TotalLocations returns the total number of unique locations in the relations.
func (r *Relations) TotalLocations() int {
	return len(r.DatesLocations)
}
