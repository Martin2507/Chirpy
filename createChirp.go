package main

import (
	"encoding/json"
	"fmt"
	"main/internal/auth"
	"main/internal/database"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type jsonBody struct {
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type ChirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	cleaned_body, err := validationChirp(r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	params := database.CreateChirpParams{
		Body:   cleaned_body.Body,
		UserID: userID,
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := ChirpResponse{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, resp)

}

func validationChirp(r *http.Request) (jsonBody, error) {

	// Decode JSON Requiest
	decoder := json.NewDecoder(r.Body)
	params := jsonBody{}
	err := decoder.Decode(&params)

	if err != nil {
		return jsonBody{}, err
	}

	// Check if body of the requiest matches the length reuirements
	if len(params.Body) > 140 {
		return jsonBody{}, fmt.Errorf("Chirp is too long")
	}

	// Split the reuqiest body and clear any profanities found inside
	cleaned_params := jsonBody{}
	splitString := strings.Split(params.Body, " ")

	for i := range splitString {
		if strings.ToLower(splitString[i]) == "kerfuffle" || strings.ToLower(splitString[i]) == "sharbert" || strings.ToLower(splitString[i]) == "fornax" {
			splitString[i] = "****"
		}
	}

	cleanedString := strings.Join(splitString, " ")

	cleaned_params.Body = cleanedString

	// Return Clean params to the handlerCreateChirp
	return cleaned_params, nil
}
