package httpserver

import (
	"context"
	"net/http"
	"time"
)

func GetLink(store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		url, err := store.Get(ctx, code)

		if err != nil {
			writeStorageError(w, err)
			return
		}

		var res struct {
			Code  string `json:"code"`
			URL   string `json:"url"`
			Short string `json:"short"`
		}

		res.Code = code
		res.URL = url
		res.Short = "http://localhost:8080/" + code

		writeJSON(w, http.StatusOK, res)
	}
}
