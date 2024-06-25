package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hits:", cfg.fileserverHits)
	fmt.Fprintln(w, "Hits:", cfg.fileserverHits)
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{fileserverHits: 0}

	mux := http.NewServeMux()
	// mux.Handle("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))))
	fileServer := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer)))

	mux.HandleFunc("/metrics", apiCfg.metricHandler)
	mux.HandleFunc("/reset", apiCfg.reset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
