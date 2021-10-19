-- name: GetUser :one
SELECT * FROM users
    WHERE userid = ? LIMIT 1;

-- name: UserIDExists :one
SELECT count(*) FROM users
    WHERE userid = ? LIMIT 1;

-- name: GetUserByName :one
SELECT count(*) FROM users
    WHERE first_name = ? AND last_name = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
    ORDER BY last_name;

-- name: CreateUser :exec
INSERT INTO users (
  userid, first_name, last_name, data
) VALUES (
  ?, ?, ?, ?
);

-- name: UpdateUser :exec
UPDATE users SET first_name = ?, last_name = ?, data = ?
    WHERE userid = ?;

-- name: UpdateUserIDinAttendance :exec
UPDATE attendance
SET userid = ?
WHERE userid = ?;

-- name: UpdateUserIDinUsers :exec
UPDATE users
SET userid = ?
WHERE userid = ?;

-- name: DeleteUser :exec
DELETE FROM users
    WHERE userid = ?;

-- name: CheckinUser :exec
INSERT INTO attendance (
  userid, date
) VALUES (
  ?, ?
);

-- name: IsUserCheckedIn :one
SELECT count(*) FROM attendance
    WHERE date = ? AND userid = ?;

-- name: GetMeetings :many
SELECT DISTINCT date FROM attendance;

-- name: GetAttendance :many
-- SELECT * FROM attendance;
SELECT * FROM attendance JOIN users ON users.userid=attendance.userid;
-- SELECT *, count(*) AS total FROM attendance JOIN users ON users.userid=attendance.userid GROUP BY date;
