package models

import "time"

type Registration struct {
	ID               int       `json:"id"`
	EventID          int       `json:"event_id"`
	UserID           int       `json:"user_id"`
	RegistrationDate time.Time `json:"registration_date"`
}

type RegistrationResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}
