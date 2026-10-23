package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/StepanDubrovin/shortener/internal/storage"
)

type URLStore interface {
	Save(url string) (string, error)
	Get(code string) (string, error)
}

func Redirect(store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		url, err := store.Get(code)

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

		code, err := store.Save(req.URL)

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
