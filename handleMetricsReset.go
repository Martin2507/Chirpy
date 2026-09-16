package main

import "net/http"

func (cfg *apiConfig) handleMetricsReset(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)

}
