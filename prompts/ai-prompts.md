# AI Prompts Used

This file contains the prompts used during the development of the Event Registration & Ticketing System, fulfilling the project requirements for AI transparency.

## Prompt 1: Project Scaffolding
**Prompt**: 
> "make this structure" (attached image of the directory tree: event-registration-system containing cmd, internal/handlers, internal/services, internal/repository, internal/models, internal/middleware, pkg/database, migrations, docs, prompts)

**Action**: 
The AI was used to execute shell commands creating the exact directory architecture presented in the image, and initializing the `go.mod` file.

## Prompt 2: Core Logic and Race Condition Handling
**Prompt**: 
> "Event Registration & Ticketing System
Build a REST API for event registration and ticketing (like Eventbrite). Users can browse events, register
for events with limited capacity, and organizers can create events and manage registrations. The critical
challenge is handling concurrent registrations to prevent overbooking when multiple users try to register
for the last few spots simultaneously.
Deliverables
Complete REST API
Database schema with proper constraints
Concurrent booking test (simulate multiple users booking last spot)
README with concurrency strategy explanation
Document explaining your approach to preventing race conditions , This all have to done in go language pls make it"

**Action**: 
The AI was used to generate:
- The SQLite tables enforcing `CHECK (available_spots >= 0)`.
- The atomic `UPDATE events SET available_spots = available_spots - 1 WHERE id = ? AND available_spots > 0` logic wrapped inside the `Serializable` transaction within `registration_repository.go`.
- The testing file utilizing `sync.WaitGroup` to fire 100 concurrent mock requests evaluating the lock.
- The `README.md` and `docs/design.md` analyzing the architecture.

## Prompt 3: CGO/SQLite Driver Fixes
**Prompt**: 
> "i have indtalled / same error / how to check and where to check server"

**Action**: 
The AI was used to diagnose that `gcc` was missing from the Windows PATH, making native `go-sqlite3` fail. It rewrote the database logic to universally bind to `modernc.org/sqlite`, a 100% pure-Go driver that circumvents the need for Windows CGO compilers, and updated the PRAGMA flags to natively support `WAL` logging mode to alleviate locking issues during the test.
