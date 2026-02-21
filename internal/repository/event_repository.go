package repository

import (
	"database/sql"
	"event-registration-system/internal/models"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) CreateEvent(e *models.Event) (int64, error) {
	query := `INSERT INTO events (title, description, date, capacity, available_spots, organizer_id) 
              VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.db.Exec(query, e.Title, e.Description, e.Date, e.Capacity, e.Capacity, e.OrganizerID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) GetEvents() ([]models.Event, error) {
	query := `SELECT id, title, description, date, capacity, available_spots, organizer_id FROM events`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Date, &e.Capacity, &e.AvailableSpots, &e.OrganizerID); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
