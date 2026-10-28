package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	mux.HandleFunc("GET /api/links/{code}", httpserver.GetLink(store))

	mux.HandleFunc("DELETE /api/links/{code}", httpserver.DeleteLink(store))

	stopCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serverErr := make(chan error, 1)

	log.Println("сервер: http://localhost:8080")

	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case <-stopCtx.Done():
		log.Println("shutdown requested")
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("server error: ", err)
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(shutdownCtx)

	if err != nil {
		log.Println("shutdown error: ", err)
		server.Close()
		return
	}

	log.Println("server stopped")
}
