package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

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

		if err != nil {
			writeStorageError(w, err)
			return
		}

		var res struct {
			Code  string `json:"code"`
			Short string `json:"short"`
		}

		res.Code = code
		res.Short = "http://localhost:8080/" + code

		writeJSON(w, http.StatusCreated, res)
	}
}
