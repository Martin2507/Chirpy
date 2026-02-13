package main

import (
	"database/sql"
	"fmt"
	"log"
	"main/internal/database"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	secret := os.Getenv("SECRET")

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	polka_Key := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		fmt.Println(err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      secret,
		polkaKey:       polka_Key,
	}

	mux := http.NewServeMux()

	// 1. Static Assets & Fileserver
	mux.Handle("/app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))

	// 2. Diagnostics & Metrics (Infrastructure)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	// 3. Admin Actions (Restricted)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// 4. Authentication & User Management
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateNewUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogIn)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	// 5. Chirps (Domain Resource)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	// 6. Webhooks
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerWebhooks)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
