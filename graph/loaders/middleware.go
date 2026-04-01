package loaders

import (
	"context"
	"database/sql"
	"net/http"
)

func Middleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ldr := NewLoaders(db)
		ctx := context.WithValue(r.Context(), loadersKey, ldr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
