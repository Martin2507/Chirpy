package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	type ChirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	dbChirps, err := cfg.db.GetAllChirpsByDate(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirpList := []ChirpResponse{}

	for _, record := range dbChirps {

		tempRes := ChirpResponse{
			ID:        record.ID,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
			Body:      record.Body,
			UserID:    record.UserID,
		}

		chirpList = append(chirpList, tempRes)
	}

	respondWithJSON(w, http.StatusOK, chirpList)
}
