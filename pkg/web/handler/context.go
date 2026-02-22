package handler

import (
	"context"
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/donseba/go-htmx"
)

type contextKey string

const (
	appKey  contextKey = "app"
	htmxKey contextKey = "htmx"
)

func AppMiddleware(app *service.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), appKey, app)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HTMXMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxInstance := htmx.New()
		handler := htmxInstance.NewHandler(w, r)
		ctx := context.WithValue(r.Context(), htmxKey, handler)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func App(r *http.Request) *service.App {
	return r.Context().Value(appKey).(*service.App)
}

func HTMX(r *http.Request) *htmx.Handler {
	return r.Context().Value(htmxKey).(*htmx.Handler)
}
