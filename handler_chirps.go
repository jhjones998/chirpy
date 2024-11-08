package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func cleanChirp(chirp string) string {
	badWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	cleanedChirp := []string{}
	chirpWords := strings.Split(chirp, " ")
	for _, word := range chirpWords {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			cleanedChirp = append(cleanedChirp, "****")
		} else {
			cleanedChirp = append(cleanedChirp, word)
		}
	}
	return strings.Join(cleanedChirp, " ")
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type chirpBody struct {
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error getting token: %v", err))
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't validate user: %v", err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirpBody{}
	err = decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding request: %v", err))
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleanedChirp := cleanChirp(chirp.Body)
	res, err := cfg.queries.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanedChirp, UserID: userID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating chirp: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, res)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authorId := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "asc"
	}
	if authorId != "" {
		authorID, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing author ID: %v", err))
			return
		}
		if sort == "asc" {
			res, err := cfg.queries.GetChirpsByUserId(r.Context(), authorID)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting chirps: %v", err))
				return
			}
			respondWithJSON(w, http.StatusOK, res)
		} else {
			res, err := cfg.queries.GetChirpsByUserIdDesc(r.Context(), authorID)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting chirps: %v", err))
				return
			}
			respondWithJSON(w, http.StatusOK, res)
		}
		return
	}
	if sort == "asc" {
		res, err := cfg.queries.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting chirps: %v", err))
			return
		}
		respondWithJSON(w, http.StatusOK, res)
	} else {
		res, err := cfg.queries.GetChirpsDesc(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting chirps: %v", err))
			return
		}
		respondWithJSON(w, http.StatusOK, res)
	}
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	chirpID := r.PathValue("chirpID")
	uuidChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing chirp ID: %v", err))
		return
	}
	res, err := cfg.queries.GetChirp(r.Context(), uuidChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error getting chirp: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error getting token: %v", err))
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't validate user: %v", err))
		return
	}

	chirpID := r.PathValue("chirpID")
	uuidChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing chirp ID: %v", err))
		return
	}
	chirp, err := cfg.queries.GetChirp(r.Context(), uuidChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error getting chirp: %v", err))
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "User does not own chirp")
		return
	}
	err = cfg.queries.DeleteChirp(r.Context(), database.DeleteChirpParams{ID: uuidChirpID, UserID: userID})
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error deleting chirp: %v", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
