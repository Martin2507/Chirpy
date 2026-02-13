package main

import (
	"main/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	urlIDValue, err := uuid.Parse(r.PathValue("chirpID"))
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

	chirp, err := cfg.db.GetChirpByID(r.Context(), urlIDValue)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "No chirps were found")
		return
	}

	dbError := cfg.db.DeleteChirpByID(r.Context(), urlIDValue)
	if dbError != nil {
		respondWithError(w, http.StatusInternalServerError, dbError.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
