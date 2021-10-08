-- name: GetUser :one
SELECT * FROM users
    WHERE userid = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
    ORDER BY name;

-- name: CreateUser :exec
INSERT INTO users (
  userid, name, data
) VALUES (
  ?, ?, ?
);

-- name: UpdateUser :exec
UPDATE users SET name = ?, data = ?
    WHERE userid = ?;

-- name: DeleteUser :exec
DELETE FROM users
    WHERE userid = ?;

-- ----------------------

-- name: CheckinUser :exec
INSERT INTO attendance (
  userid, date
) VALUES (
  ?, ?
);