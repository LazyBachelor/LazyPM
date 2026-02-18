package handler

import (
	"context"
	"net/http"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/donseba/go-htmx"
)

type contextKey string

const (
	servicesKey contextKey = "services"
	htmxKey     contextKey = "htmx"
)

func ServicesMiddleware(svc *service.Services) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), servicesKey, svc)
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

func Services(r *http.Request) *service.Services {
	return r.Context().Value(servicesKey).(*service.Services)
}

func HTMX(r *http.Request) *htmx.Handler {
	return r.Context().Value(htmxKey).(*htmx.Handler)
}
