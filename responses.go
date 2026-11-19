package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, message string) {

	returnError := JSONError{
		Error: message,
	}

	data, err := encodeJson(returnError)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)

}

func respondWithJSON(w http.ResponseWriter, payload any) {

	data, err := encodeJson(payload)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func encodeJson(payloda any) ([]byte, error) {

	data, err := json.Marshal(payloda)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return nil, err
	}

	return data, nil
}

func decodeJson(body io.Reader, shape any) error {

	decoder := json.NewDecoder(body)
	err := decoder.Decode(shape)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		return err
	}

	return nil

}
