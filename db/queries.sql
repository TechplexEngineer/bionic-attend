-- name: GetUser :one
SELECT * FROM users
    WHERE userid = ? LIMIT 1;

-- name: UserIDExists :one
SELECT count(*) FROM users
    WHERE userid = ? AND hidden = FALSE LIMIT 1;

-- name: GetUserByName :one
SELECT count(*) FROM users
    WHERE first_name = ? AND last_name = ? AND hidden = FALSE LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
		WHERE  hidden = FALSE
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

-- name: SoftDeleteUser :exec
UPDATE users
		SET hidden = TRUE
    WHERE userid = ?;

-- name: UnHideUser :exec
UPDATE users
		SET hidden = FALSE
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
SELECT DISTINCT date, count(*) FROM attendance GROUP BY date;

-- name: GetAttendance :many
SELECT * FROM attendance JOIN users ON users.userid=attendance.userid WHERE users.hidden = FALSE;
-- SELECT *, count(*) AS total FROM attendance JOIN users ON users.userid=attendance.userid GROUP BY date;
