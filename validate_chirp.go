package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerValidationChirp(w http.ResponseWriter, r *http.Request) {

	type jsonBody struct {
		Body string `json:"body"`
	}

	type jsonCleanedBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := jsonBody{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned_params := jsonCleanedBody{}
	splitString := strings.Split(params.Body, " ")

	for i := range splitString {
		if strings.ToLower(splitString[i]) == "kerfuffle" || strings.ToLower(splitString[i]) == "sharbert" || strings.ToLower(splitString[i]) == "fornax" {
			splitString[i] = "****"
		}
	}

	cleanedString := strings.Join(splitString, " ")
	cleaned_params.CleanedBody = cleanedString

	respondWithJSON(w, http.StatusOK, cleaned_params)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	w.Header().Set("Content-Type", "application/json")

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(code)
	w.Write(jsonPayload)
}
