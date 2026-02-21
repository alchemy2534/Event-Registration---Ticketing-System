# Event Registration & Ticketing System

A REST API for event registration and ticketing built with Go and SQLite, simulating an Eventbrite-like experience. The core focus of this project is to safely and deterministically handle high-concurrency race conditions — specifically when hundreds of users simultaneously attempt to register for an event's final few available spots.

## Key Features
- **RESTful API**: Browse events, register for an event, and organize new events.
- **SQLite Database**: Lightweight data store utilizing schema constraints.
- **Concurrency Saftey**: Strict mechanisms are implemented to prevent overbooking via an `Atomic Updates` + `Database Transactions` strategy.

## Tech Stack
- **Go**: 1.21+
- **Database**: SQLite (built-in standard library with `github.com/mattn/go-sqlite3`)
- **Routing**: `net/http` standard mux

## API Endpoints

### Organizers
- `POST /api/events/create`
  - Body: `{ "title": "Tech Conference", "description": "Annual tech meet", "date": "2026-10-10T10:00:00Z", "capacity": 100, "organizer_id": 1 }`
  - **Returns**: Event object details

### Users
- `GET /api/events`
  - **Returns**: A JSON array of all existing events, displaying their up-to-date capacity vs available spots.
- `POST /api/register`
  - Body: `{ "event_id": 1, "user_id": 2 }`
  - **Returns**: `{"message": "Registered successfully", "status": true}` on success, or a `409 Conflict`/`400 Bad Request` if the event is full or user already signed up. 

## The Concurrency Challenge Strategy

When multiple users try to register for the very last spot at the exact same millisecond, typical `SELECT capacity THEN UPDATE if > 0` patterns result in a race condition. Both concurrent requests could read that `capacity > 0` before either one finishes updating the row, leading to *both* users successfully booking the exact same physical spot, thereby exceeding maximum capacity.

### Approach: Database Constraints & Atomic Actions

The application resolves race conditions using the following principles:

1. **Atomic Update within Transaction**: Rather than selecting data into Go memory and doing conditional checks before writing back, the condition gets moved into the very `UPDATE` row lock itself inside a single atomic database statement using:
   ```sql
   UPDATE events SET available_spots = available_spots - 1 WHERE id = ? AND available_spots > 0
   ```
   *If the spots are already 0, the database returns `0 rows affected`. Go then knows it failed.*

2. **Transaction Isolation**: 
   All this happens enveloped within a transaction wrapper (`sql.LevelSerializable`).

3. **Database Constraints**: The SQLite layout acts as a secondary failsafe:
   ```sql
   available_spots INTEGER NOT NULL CHECK (available_spots >= 0)
   UNIQUE(event_id, user_id)
   ```
   Applying the explicit `CHECK` prevents numerical capacities dropping into the negatives regardless. Similarly, the unique-composite key restricts users duplicate-buying.

## Run Concurrency Test
An automated unit test is included to violently parallel-bomb the database utilizing Go's `sync.WaitGroup` to spawn 100 simultaneous goroutines attempting to secure a 10-seat event exactly at identical moments.

```bash
# Navigate to the correct directory
cd internal/repository
# Run test and simulate race conditions
go test -v
```

## How to Run the Server

```bash
# Download dependencies
go mod tidy 

# Start API
go run cmd/main.go
# Server listens on :8080 by default
```
