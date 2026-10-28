package httpserver

import (
	"context"
	"net/http"
	"time"
)

func Redirect(store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		url, err := store.Get(ctx, code)

		if err != nil {
			writeStorageError(w, err)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
