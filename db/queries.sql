-- name: GetUser :one
SELECT * FROM users
    WHERE userid = ? LIMIT 1;

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
    WHERE date = ? AND userid = ?
