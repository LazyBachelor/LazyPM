package server

import (
	"embed"
	"net/http"
	"strings"
	"time"

	"github.com/LazyBachelor/LazyPM/pkg/web/handler"
	"github.com/NYTimes/gziphandler"
	"github.com/rs/cors"
)

func (s *Server) RegisterRoutes(assets embed.FS) http.Handler {
	mux := http.NewServeMux()

	s.handleAssets(mux, assets)

	for _, route := range handler.GetRoutes(s.Services) {
		mux.Handle(route.Pattern, route.Handler)
	}

	handler := cors.Default().Handler(mux)

	return gziphandler.GzipHandler(handler)
}

func (s *Server) handleAssets(mux *http.ServeMux, assets embed.FS) {
	mux.Handle("/assets/",
		http.StripPrefix("/assets/",
			http.FileServer(http.Dir("pkg/web/assets"))))

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeContent(w, r, "robots.txt", time.Now(), strings.NewReader("User-agent: *\nAllow: /"))
	})
}
