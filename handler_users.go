package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type userBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	database.GetUserByEmailRow
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := userBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding request: %v", err))
		return
	}
	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error hashing password: %v", err))
		return
	}
	res, err := cfg.queries.CreateUser(
		r.Context(),
		database.CreateUserParams{Email: body.Email, HashedPassword: hashedPassword},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, res)
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := userBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding request: %v", err))
		return
	}
	hashedPassword, err := cfg.queries.GetUserHashedPasswordByEmail(r.Context(), body.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error getting user: %v", err))
		return
	}
	err = auth.CheckPasswordHash(body.Password, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}
	user, err := cfg.queries.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting user: %v", err))
		return
	}
	expirationTime := time.Hour
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.secretToken,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating access token: %v", err))
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating refresh token: %v", err))
		return
	}
	_, err = cfg.queries.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{Token: refreshToken, UserID: user.ID},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error creating refresh token in db: %v", err),
		)
		return
	}

	respondWithJSON(w, http.StatusOK, tokenResponse{
		GetUserByEmailRow: user,
		Token:             accessToken,
		RefreshToken:      refreshToken,
	})
}

func (cfg *apiConfig) refreshLoginToken(w http.ResponseWriter, r *http.Request) {
	type refreshTokenBody struct {
		Token string `json:"token"`
	}
	w.Header().Set("Content-Type", "application/json")
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error getting bearer token: %v", err))
		return
	}
	res, err := cfg.queries.GetRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error getting refresh token: %v", err))
		return
	}
	if time.Now().UTC().After(res.ExpiresAt) || res.RevokedAt.Time != (time.Time{}) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired or revoked")
		return
	}
	expirationTime := time.Hour
	accessToken, err := auth.MakeJWT(
		res.UserID,
		cfg.secretToken,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating access token: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, refreshTokenBody{Token: accessToken})
}

func (cfg *apiConfig) revokeLoginToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error getting bearer token: %v", err))
		return
	}
	err = cfg.queries.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error revoking refresh token: %v", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := userBody{}
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error getting bearer token: %v", err))
		return
	}
	userID, err := auth.ValidateJWT(bearerToken, cfg.secretToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error validating JWT: %v", err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding request: %v", err))
		return
	}
	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error hashing password: %v", err))
		return
	}
	res, err := cfg.queries.UpdateUser(
		r.Context(),
		database.UpdateUserParams{Email: body.Email, HashedPassword: hashedPassword, ID: userID},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *apiConfig) upgradeUserWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type webhookBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	w.Header().Set("Content-Type", "application/json")
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error getting API key: %v", err))
		return
	}
	body := webhookBody{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding request: %v", err))
		return
	}
	if body.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(body.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing user ID: %v", err))
		return
	}
	user, err := cfg.queries.UpgradeUser(r.Context(), database.UpgradeUserParams{ID: userID, IsChirpyRed: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upgrading user: %v", err))
		return
	}
	if user != (database.UpgradeUserRow{}) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		respondWithError(w, http.StatusNotFound, "User not found")
	}
}
