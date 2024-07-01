package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	filepath := filepath.Join("admin", "metrics.html")

	htmlContent, err := os.ReadFile(filepath)

	formattedHtml := fmt.Sprintf(string(htmlContent), cfg.fileserverHits)

	if err != nil {
		log.Fatal("Something went wrong while loading the html file")
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(formattedHtml))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

type User struct {
	FirstName string `json:"first_name"`
	BirthYear int    `json:"birth_year"`
	Email     string
}

func testing(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Start testing...")
	data, error := json.Marshal(User{
		FirstName: "Moe",
		BirthYear: 1996,
		Email:     "mahmoud.othman.ajam@gmail.com",
	})

	if error != nil {
		log.Fatal("Something went wrong...")
	}

	w.Write(data)

	fmt.Println(string(data))
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{fileserverHits: 0}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer)))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)
	mux.HandleFunc("/api/reset", apiCfg.reset)
	mux.HandleFunc("/api/test", testing)
	mux.HandleFunc("POST /api/chirps", createChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
