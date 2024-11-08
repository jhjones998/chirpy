// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: chirps.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createChirp = `-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
        gen_random_uuid(),
        now() at time zone 'utc',
        now() at time zone 'utc',
        $1,
        $2
)
RETURNING id, created_at, updated_at, body, user_id
`

type CreateChirpParams struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirp, arg.Body, arg.UserID)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const deleteChirp = `-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1
AND user_id = $2
`

type DeleteChirpParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteChirp(ctx context.Context, arg DeleteChirpParams) error {
	_, err := q.db.ExecContext(ctx, deleteChirp, arg.ID, arg.UserID)
	return err
}

const deleteChirps = `-- name: DeleteChirps :exec
DELETE FROM chirps
`

func (q *Queries) DeleteChirps(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteChirps)
	return err
}

const getChirp = `-- name: GetChirp :one
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE id = $1
`

func (q *Queries) GetChirp(ctx context.Context, id uuid.UUID) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getChirp, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const getChirps = `-- name: GetChirps :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
ORDER BY created_at
`

func (q *Queries) GetChirps(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getChirpsByUserId = `-- name: GetChirpsByUserId :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE user_id = $1
ORDER BY created_at
`

func (q *Queries) GetChirpsByUserId(ctx context.Context, userID uuid.UUID) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getChirpsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getChirpsByUserIdDesc = `-- name: GetChirpsByUserIdDesc :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE user_id = $1
ORDER BY created_at DESC
`

func (q *Queries) GetChirpsByUserIdDesc(ctx context.Context, userID uuid.UUID) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getChirpsByUserIdDesc, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getChirpsDesc = `-- name: GetChirpsDesc :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
ORDER BY created_at DESC
`

func (q *Queries) GetChirpsDesc(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getChirpsDesc)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}