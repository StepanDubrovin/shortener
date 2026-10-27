package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/StepanDubrovin/shortener/internal/storage"
)

type URLStore interface {
	Save(ctx context.Context, url string) (string, error)
	Get(ctx context.Context, code string) (string, error)
}

func Redirect(store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		url, err := store.Get(ctx, code)

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "storage timeout", http.StatusGatewayTimeout)
			return
		}

		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "failed to get url", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}

func Shorten(store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			URL string `json:"url"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, "url is required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		code, err := store.Save(ctx, req.URL)

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "storage timeout", http.StatusGatewayTimeout)
			return
		}

		if err != nil {
			http.Error(w, "failed to save url", http.StatusInternalServerError)
			return
		}

		var res struct {
			Code  string `json:"code"`
			Short string `json:"short"`
		}

		res.Code = code
		res.Short = "http://localhost:8080/" + code

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	}
}
