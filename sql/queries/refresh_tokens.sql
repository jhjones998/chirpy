-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
           $1,
           now() at time zone 'utc',
           now() at time zone 'utc',
           $2,
           now() at time zone 'utc' + interval '60 days'
       )
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = now() at time zone 'utc', updated_at = now() at time zone 'utc'
WHERE token = $1;