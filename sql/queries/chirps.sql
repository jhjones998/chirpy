-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
        gen_random_uuid(),
        now() at time zone 'utc',
        now() at time zone 'utc',
        $1,
        $2
)
RETURNING *;

-- name: DeleteChirps :exec
DELETE FROM chirps;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at;

-- name: GetChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at DESC;

-- name: GetChirpsByUserId :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at;

-- name: GetChirpsByUserIdDesc :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1
AND user_id = $2;