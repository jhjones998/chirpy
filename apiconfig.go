package main

import (
	"chirpy/internal/database"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	platform       string
	secretToken    string
	polkaKey       string
}
