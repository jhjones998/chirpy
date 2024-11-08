-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
        gen_random_uuid(),
       now() at time zone 'utc',
    now() at time zone 'utc',
        $1,
        $2
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: GetUserHashedPasswordByEmail :one
SELECT hashed_password FROM users
WHERE email = $1;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, is_chirpy_red FROM users
WHERE email = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: UpdateUser :one
UPDATE users
SET email = $1, updated_at = now() at time zone 'utc', hashed_password = $2
WHERE id = $3
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpgradeUser :one
UPDATE users
SET is_chirpy_red = $1
WHERE id = $2
RETURNING id, created_at, updated_at, email, is_chirpy_red;