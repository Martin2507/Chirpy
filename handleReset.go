package main

import (
	"net/http"
)

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {

	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "")
		return
	}

	cfg.FileserverHits.Store(0)

	_, err := cfg.DB.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
