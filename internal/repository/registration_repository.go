package repository

import (
	"context"
	"database/sql"
	"errors"
	"event-registration-system/internal/models"
)

type RegistrationRepository struct {
	db *sql.DB
}

func NewRegistrationRepository(db *sql.DB) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

// RegisterForEvent uses an atomic update and database transaction to handle concurrent bookings safely.
func (r *RegistrationRepository) RegisterForEvent(ctx context.Context, reg *models.Registration) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Atomic update to decrement available spots if > 0
	// This prevents race conditions where multiple users try to book the last spot simultaneously.
	// Only those updates that actually decrease the value from > 0 will succeed.
	queryUpdate := `UPDATE events SET available_spots = available_spots - 1 WHERE id = ? AND available_spots > 0`
	res, err := tx.ExecContext(ctx, queryUpdate, reg.EventID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("event is full or does not exist")
	}

	// Insert registration
	queryInsert := `INSERT INTO registrations (event_id, user_id) VALUES (?, ?)`
	_, err = tx.ExecContext(ctx, queryInsert, reg.EventID, reg.UserID)
	if err != nil {
		// E.g., unique constraint violation if user already registered for this event
		return err
	}

	return tx.Commit()
}
