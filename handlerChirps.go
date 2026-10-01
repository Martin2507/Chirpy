package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/martin/Chirpy/internal/auth"
	"github.com/martin/Chirpy/internal/database"
)

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	chirp := Payload{}
	err := decoder.Decode(&chirp)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	profanitiesArr := []string{"kerfuffle", "sharbert", "fornax"}
	splitBody := strings.Split(chirp.Body, " ")

	for i := range splitBody {
		for j := range profanitiesArr {
			if strings.ToLower(splitBody[i]) == profanitiesArr[j] {
				splitBody[i] = "****"
			}
		}
	}

	cleaned_body := strings.Join(splitBody, " ")

	returnPayload := Payload{Body: cleaned_body, UserID: userID}

	params := database.CreateNewChirpParams{
		Body:   returnPayload.Body,
		UserID: userID,
	}

	data, err := cfg.DB.CreateNewChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	returnData := Chirp{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserID:    data.UserID,
	}

	respondWithJSON(w, http.StatusCreated, returnData)

}

func (cfg *apiConfig) handleGetChirpsByOrder(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.SortChirpByCreated_At(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	result := []Chirp{}

	for _, chirp := range chirps {
		result = append(result, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, result)

}

func (cfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {

	chirpToGet := r.PathValue("chirpID")

	uuidChirp, err := uuid.Parse(chirpToGet)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirp, err := cfg.DB.GetChirpByID(r.Context(), uuidChirp)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	returnData := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, returnData)
}

func (cfg *apiConfig) handleChirpDeleteByID(w http.ResponseWriter, r *http.Request) {

	chirpToDelete := r.PathValue("chirpID")

	uuidChirp, err := uuid.Parse(chirpToDelete)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirp, err := cfg.DB.DeleteChirpByID(r.Context(), uuidChirp)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	returnData := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, returnData)

}
