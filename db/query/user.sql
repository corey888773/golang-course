-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    email,
    full_name
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE username = $1;


-- name: UpdateUser :one
UPDATE users SET 
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password), 
    password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at), 
    email = COALESCE(sqlc.narg(email), email), 
    full_name = COALESCE(sqlc.narg(full_name), full_name)
WHERE 
    username = sqlc.arg(username)
RETURNING *;