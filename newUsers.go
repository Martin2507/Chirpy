package main

import (
	"encoding/json"
	"main/internal/database"
	"net/http"
	"time"

	"main/internal/auth"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateNewUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
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
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	hPassword, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	tempParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hPassword,
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), tempParams)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	respUser := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, respUser)

}
