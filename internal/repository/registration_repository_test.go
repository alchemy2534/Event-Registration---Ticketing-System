package repository

import (
	"context"
	"database/sql"
	"event-registration-system/internal/models"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	_ "modernc.org/sqlite"
)

// TestConcurrentRegistration simulates multiple users trying to book the last spots.
func TestConcurrentRegistration(t *testing.T) {
	// Create an db file for testing
	os.Remove("test_db.sqlite") // Clean up old ones just in case
	os.Remove("test_db.sqlite-wal")
	os.Remove("test_db.sqlite-shm")
	db, err := sql.Open("sqlite", "test_db.sqlite?_pragma=busy_timeout=10000&_pragma=journal_mode=WAL")
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove("test_db.sqlite")
	}()
	defer db.Close()

	// 1. Initialize DB Scheme
	schema1 := `CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, email TEXT UNIQUE NOT NULL);`
	schema2 := `CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL, description TEXT, date DATETIME NOT NULL, capacity INTEGER NOT NULL, available_spots INTEGER NOT NULL CHECK (available_spots >= 0), organizer_id INTEGER NOT NULL);`
	schema3 := `CREATE TABLE IF NOT EXISTS registrations (id INTEGER PRIMARY KEY AUTOINCREMENT, event_id INTEGER NOT NULL, user_id INTEGER NOT NULL, UNIQUE(event_id, user_id));`

	for _, s := range []string{schema1, schema2, schema3} {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("Failed to initialize test schema: %v", err)
		}
	}

	// 2. Setup mock data: 1 event with capacity 10, and 100 users
	if _, err := db.Exec("INSERT INTO users (name, email) VALUES ('Organizer', 'org@admin.com')"); err != nil {
		t.Fatalf("Failed to insert mock user: %v", err)
	}
	if _, err := db.Exec("INSERT INTO events (title, date, capacity, available_spots, organizer_id) VALUES ('Exclusive Tech Meetup', '2026-10-10', 10, 10, 1)"); err != nil {
		t.Fatalf("Failed to insert mock event: %v", err)
	}

	for i := 1; i <= 100; i++ {
		email := fmt.Sprintf("user%d@test.com", i) // safely unique email for tests
		if _, err := db.Exec("INSERT INTO users (name, email) VALUES ('TestUser', ?)", email); err != nil {
			t.Fatalf("Failed to insert mock user loop %d: %v", i, err)
		}
	}

	// 3. Prepare Registration Repository
	repo := NewRegistrationRepository(db)

	var successfulRegistrations int32
	var failedRegistrations int32

	var wg sync.WaitGroup

	// 4. Simulate 100 concurrent requests for only 10 spots
	workers := 100
	for i := 1; i <= workers; i++ {
		wg.Add(1)

		go func(userID int) {
			defer wg.Done()

			err := repo.RegisterForEvent(context.Background(), &models.Registration{
				EventID: 1,          // 'Exclusive Tech Meetup'
				UserID:  userID + 1, // Offset organizer ID
			})

			if err == nil {
				atomic.AddInt32(&successfulRegistrations, 1)
			} else {
				atomic.AddInt32(&failedRegistrations, 1)
			}
		}(i)
	}

	wg.Wait()

	// 5. Assertions

	if successfulRegistrations != 10 {
		t.Errorf("Expected exactly 10 successful registrations, got %d", successfulRegistrations)
	}

	if failedRegistrations != 90 {
		t.Errorf("Expected exactly 90 failed registrations, got %d", failedRegistrations)
	}

	// Verify database state matches
	var finalSpots int
	err = db.QueryRow("SELECT available_spots FROM events WHERE id = 1").Scan(&finalSpots)
	if err != nil {
		t.Fatalf("Failed to check final spots: %v", err)
	}

	if finalSpots != 0 {
		t.Errorf("Expected 0 available spots, got %d", finalSpots)
	}

	var countRegs int
	err = db.QueryRow("SELECT COUNT(*) FROM registrations WHERE event_id = 1").Scan(&countRegs)
	if countRegs != 10 {
		t.Errorf("Expected exactly 10 registration rows, got %d", countRegs)
	}
}
