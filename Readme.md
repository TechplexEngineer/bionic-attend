Bionic Attendance
=================

Web application to track student attendance for FRC robotics team meetings.

## How it works
Users can be created by entering a userid that does not already exist in the system.
Designed to be used with a barcode scanner which appends a "return" to the end of the scan.
Does not track time, only if the student attended a meeting.

## How we use it
We have a raspberry pi with a keyboard, mouse, and barcode scanner near the entrance to our robotics shop.
Students scan their student IDs when they enter the shop.

## Technical
Data is stored in a SQLite database in two tables, users and attendance.
The database structure is specified in [schema.sql](src/db/schema.sql).
THe queries are specified in [queries.sql](src/db/queries.sql).
This project uses [sqlc](https://github.com/kyleconroy/sqlc) to generate go code which converts the query results into
go datatypes. 

### Dependencies
- modernc.org/sqlite - pure go implementation of sqlite
- github.com/gorilla/mux - for http routing
- github.com/gorilla/sessions - for cross browser cookie creation and deletion
- github.com/matryer/is - lightweight unit test helpers