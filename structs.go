package main

import "sync/atomic"

type apiConfig struct {
	fileserverHits atomic.Int32
}

type JSONBody struct {
	Body string `json:"body"`
}

type JSONError struct {
	Error string `json:"error"`
}

type JSONResponsBody struct {
	Cleaned_body string `json:"cleaned_body"`
}
