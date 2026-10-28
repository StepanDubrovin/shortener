package httpserver

import (
	"context"
	"net/http"
	"time"
)

func DeleteLink(store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		err := store.Delete(ctx, code)

		if err != nil {
			writeStorageError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
