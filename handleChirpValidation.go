package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handleChirpValidation(w http.ResponseWriter, r *http.Request) {

	params := JSONBody{}
	returnResponseBody := JSONResponsBody{}

	profanitiesList := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	err := decodeJson(r.Body, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(params.Body) >= 140 {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	splitMessage := strings.Split(params.Body, " ")

	for _, prof := range profanitiesList {

		for j, word := range splitMessage {

			if strings.ToLower(word) == prof {
				splitMessage[j] = "****"
			}
		}
	}

	returnResponseBody.Cleaned_body = strings.Join(splitMessage, " ")

	respondWithJSON(w, returnResponseBody)

}
