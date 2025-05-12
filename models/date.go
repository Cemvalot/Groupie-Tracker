package models

import "strings"

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Add helper methods if needed, e.g., to format or validate dates
func (d *Date) IsValid() bool {
	return d.ID > 0 && len(d.Dates) > 0
}

// Example method to return dates as a single formatted string
func (d *Date) FormattedDates() string {
	return strings.Join(d.Dates, ", ")
}
