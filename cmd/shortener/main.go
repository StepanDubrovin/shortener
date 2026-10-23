package main

import (
	"log"
	"net/http"
	"os"

	httpserver "github.com/StepanDubrovin/shortener/internal/http-server"
	"github.com/StepanDubrovin/shortener/internal/storage/postgres"
)

func main() {
	mux := http.NewServeMux()
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	store, err := postgres.New(dsn)

	if err != nil {
		log.Fatal("no connection", err)
	}

	defer store.Close()

	mux.HandleFunc("POST /shorten", httpserver.Shorten(store))

	mux.HandleFunc("GET /{code}", httpserver.Redirect(store))

	log.Println("сервер: http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Print(err)
	}
}
