package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/martin/Chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {

	usersParams := UserParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&usersParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if usersParams.ExpiresInSeconds <= 0 || usersParams.ExpiresInSeconds > 3600 {
		usersParams.ExpiresInSeconds = time.Duration(60*60) * time.Second
	}

	userData, err := cfg.DB.GetUserByEmail(r.Context(), usersParams.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := auth.CheckPasswordHash(usersParams.Password, userData.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "ncorrect email or password")
		return
	}

	token, err := auth.MakeJWT(userData.ID, cfg.Secret, time.Duration(usersParams.ExpiresInSeconds))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Print(usersParams.ExpiresInSeconds)

	returnData := User{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email,
		Token:     token,
	}

	respondWithJSON(w, http.StatusOK, returnData)

}
