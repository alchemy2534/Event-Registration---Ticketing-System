package models

import "time"

type Event struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Date           time.Time `json:"date"`
	Capacity       int       `json:"capacity"`
	AvailableSpots int       `json:"available_spots"`
	OrganizerID    int       `json:"organizer_id"`
}
