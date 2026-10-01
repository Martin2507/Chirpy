package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/martin/Chirpy/internal/database"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file was found or an error has occured while loading it")
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")

	const filePathRoot = "."
	const port = "8080"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
		Platform:       platform,
		Secret:         secret,
	}

	// Something
	mux.Handle("GET /app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))
	mux.Handle("GET /app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir(filePathRoot+"/assets")))))

	// Admin Functions
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

	// MISC
	mux.HandleFunc("GET /api/healthz", handleHealthz)

	// Chirp Functions
	mux.HandleFunc("POST /api/chirps", apiCfg.handleChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirpsByOrder)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetChirpByID)
	// mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleChirpDeleteByID)

	// User Function
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateNewUser)
	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)

	err = http.ListenAndServe(server.Addr, server.Handler)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}
