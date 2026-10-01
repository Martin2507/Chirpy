package main

import (
	"encoding/json"
	"net/http"

	"github.com/martin/Chirpy/internal/auth"
	"github.com/martin/Chirpy/internal/database"
)

func (cfg *apiConfig) handleCreateNewUser(w http.ResponseWriter, r *http.Request) {

	userParams := UserParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(userParams.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	params := database.CreateUserParams{
		Email:          userParams.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.DB.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	parsedUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, parsedUser)

}
