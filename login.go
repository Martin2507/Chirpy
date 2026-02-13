package main

import (
	"database/sql"
	"encoding/json"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogIn(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type User struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	correctPassword, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !correctPassword {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tempParams := database.CreateNewRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{},
	}

	dbToken, err := cfg.db.CreateNewRefreshToken(r.Context(), tempParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tokenString, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        tokenString,
		RefreshToken: dbToken.Token,
		IsChirpyRed:  dbUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, resp)

}
