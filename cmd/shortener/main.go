package main

import (
	"log"
	"net/http"

	httpserver "github.com/StepanDubrovin/shortener/internal/http-server"
	"github.com/StepanDubrovin/shortener/internal/storage"
)

func main() {
	mux := http.NewServeMux()
	store := storage.New()

	mux.HandleFunc("POST /shorten", httpserver.Shorten(store))

	mux.HandleFunc("GET /{code}", httpserver.Redirect(store))

	log.Println("сервер: http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
