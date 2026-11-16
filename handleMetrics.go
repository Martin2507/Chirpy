package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {

	returnString := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(returnString))

}
