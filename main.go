package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	apiCfg := apiConfig{
		atomic.Int32{},
	}

	mux.Handle("GET /app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))
	mux.Handle("GET /app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir(filePathRoot+"/assets")))))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleMetricsReset)

	mux.HandleFunc("GET /api/healthz", handleHealthz)

	err := http.ListenAndServe(server.Addr, server.Handler)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}
