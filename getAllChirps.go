package main

import (
	"net/http"
	"sort"
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

	optionalUserID := r.URL.Query().Get("author_id")
	optonalSort := r.URL.Query().Get("sort")

	if optionalUserID == "" {

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

		if optonalSort == "desc" {
			sort.Slice(chirpList, func(i, j int) bool {
				return chirpList[i].CreatedAt.After(chirpList[j].CreatedAt)
			})
		}

		respondWithJSON(w, http.StatusOK, chirpList)
		return

	}

	queryParams, err := uuid.Parse(optionalUserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbChirps, err := cfg.db.GetChirpsByUserID(r.Context(), queryParams)

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
