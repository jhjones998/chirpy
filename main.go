package main

import (
	"chirpy/internal/database"
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)
import _ "github.com/lib/pq"

func main() {
	const filepathRoot = "."
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		queries:        database.New(db),
		platform:       os.Getenv("PLATFORM"),
		secretToken:    os.Getenv("SECRET_TOKEN"),
		polkaKey:       os.Getenv("POLKA_KEY"),
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))),
	)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	mux.HandleFunc("POST /api/chirps", cfg.createChirp)
	mux.HandleFunc("GET /api/chirps", cfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirp)

	mux.HandleFunc("POST /api/login", cfg.login)
	mux.HandleFunc("POST /api/refresh", cfg.refreshLoginToken)
	mux.HandleFunc("POST /api/revoke", cfg.revokeLoginToken)
	mux.HandleFunc("POST /api/users", cfg.createUser)
	mux.HandleFunc("PUT /api/users", cfg.updateUser)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.upgradeUserWebhookHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s at: %s:%s\n", filepathRoot, "http://localhost", port)
	log.Fatal(srv.ListenAndServe())
}
