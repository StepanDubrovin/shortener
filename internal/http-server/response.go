package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/StepanDubrovin/shortener/internal/storage"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("failed to write JSON response: ", err)
	}
}

func writeStorageError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		http.Error(w, "storage timeout", http.StatusGatewayTimeout)
		return
	}

	if errors.Is(err, storage.ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	log.Print("storage operation failed: ", err)

	http.Error(w, "storage operation failed", http.StatusInternalServerError)
}
