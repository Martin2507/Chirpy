package main

import (
	"log"
	"net/http"
)

func main() {

	// Const Variables
	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{}

	// Create a new server using the HTTP package
	mux := http.NewServeMux()

	// Serve content of the index file to the root of the project
	mux.Handle("/app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// Serve content of the assests folder to the client
	mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// Admin Handlers
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleMetricsReset)

	// Api Handlers
	mux.HandleFunc("GET /api/healthz", handleHealthz)

	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handleChirpValidation)

	// Assigne variables to the server struct of the HTTP package
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Print server details to the console
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	// Listen and server files to the client and incase of failier log a fatal error to the console
	log.Fatal(server.ListenAndServe())

}
