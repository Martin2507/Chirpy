package main

import (
	"encoding/json"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type User struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	params := RequestBody{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tokenBearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(tokenBearer, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tempParams := database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashed_password,
	}

	dbUsers, err := cfg.db.UpdateUser(r.Context(), tempParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := User{
		ID:          dbUsers.ID,
		CreatedAt:   dbUsers.CreatedAt,
		UpdatedAt:   dbUsers.UpdatedAt,
		Email:       dbUsers.Email,
		IsChirpyRed: dbUsers.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, resp)

}
