# Design & Architecture Document

### The Race Condition Challenge
Event registrations naturally generate extremely localized and heavy simultaneous throughput on single records. Specifically, tickets for exclusive events sell out precisely during high-load periods on identical resources.

Building a standard CRUD implementation:
1. `GET` current spots for event
2. Check if Spots > 0
3. `POST` new registration
4. `PUT` (Spots = Spots - 1)

This implementation breaks instantly under concurrency since concurrent reads all pull `Spots = 1` and all process successful modifications.

### Our Solution Framework

#### 1. SQL Layer Failsafes 
Standard Relational Database management systems like SQLite handle atomic operations beautifully if designed around them. 
We introduced explicit schemas:
- `CHECK(available_spots >= 0)`: The database engine naturally rejects any mathematical attempts at rendering events over-booked.
- `UNIQUE(user_id, event_id)`: A secondary constraint strictly restricting user replication.

#### 2. Service Layer Atomicity
Inside `internal/repository/registration_repository.go`:
We compress the typical `READ/CHECK/WRITE` workflow exclusively down into a solitary query executing atop a locked row block:
`UPDATE events SET available_spots = available_spots - 1 WHERE id = 1 AND available_spots > 0`

By enforcing the read constraint inside the `WHERE` clause instead of utilizing Go's CPU memory logic, it utilizes native database mechanisms to synchronize rows efficiently, allowing only 1 caller to succeed the decrement operation when evaluating the integer values. Any subsequently evaluated update query will result in `rowsAffected = 0.`

We wrap the update and the registration creation inside of an atomic SQL transaction utilizing `Serializable` locking levels; ensuring the insert completely rolls back if the update criteria was declined. 

#### 3. Why SQLite? 
With simple setups, lightweight SQL deployments demonstrate concurrent atomicity effectively without installing large infrastructural dependencies. To swap into larger deployments (i.e., Postgres), only the `pkg/database/db.go` init payload requires adaptation — the mathematical principles mapping `UPDATE .. AND spots > 0` retain flawless integrity across platform variances.
