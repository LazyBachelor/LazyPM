package server

import (
	"embed"
	"net/http"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/a-h/templ"
	"github.com/rs/cors"
)

type Route struct {
	Pattern   string
	Component templ.Component
}

func (s *Server) RegisterRoutes(assets embed.FS, routes []Route) http.Handler {
	mux := http.NewServeMux()

	s.handleAssets(mux, assets)

	for _, route := range routes {
		mux.Handle(route.Pattern, templ.Handler(route.Component))
	}

	handler := cors.Default().Handler(mux)

	return gziphandler.GzipHandler(handler)
}

func (s *Server) handleAssets(mux *http.ServeMux, assets embed.FS) {
	fileServer := http.FileServer(http.FS(assets))
	mux.Handle("/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Optionally set long-term caching headers for static assets
		//w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		fileServer.ServeHTTP(w, r)
	}))

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeContent(w, r, "robots.txt", time.Now(), strings.NewReader("User-agent: *\nAllow: /"))
	})
}
