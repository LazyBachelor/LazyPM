package server

import (
	"embed"
	"net/http"
	"strings"
	"time"

	"github.com/LazyBachelor/LazyPM/pkg/web/handler"
	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func (s *Server) RegisterRoutes(assets embed.FS) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.CleanPath)

	r.Use(handler.HTMXMiddleware)
	r.Use(handler.AppMiddleware(s.App))

	s.handleAssets(r, assets)

	r.Get("/", handler.DashboardHandler)
	r.Get("/status", handler.HandleTaskStatus)

	r.Get("/create-issue", handler.CreateIssueFormModal)

	r.Route("/issues", func(r chi.Router) {
		r.Get("/", handler.ListIssues)
		r.Post("/", handler.CreateIssue)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(handler.IssueCtx)
			r.Get("/", handler.GetIssue)
			r.Patch("/", handler.UpdateIssue)
			r.Delete("/", handler.DeleteIssue)

			r.Route("/comments", func(r chi.Router) {
				r.Get("/", handler.ListComments)
				r.Post("/", handler.CreateComment)
			})
		})
	})
	return gziphandler.GzipHandler(r)
}

func (s *Server) handleAssets(r chi.Router, assets embed.FS) {
	fileServer := http.FileServer(http.Dir("pkg/web/assets"))

	r.Handle("/assets/*",
		http.StripPrefix("/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			fileServer.ServeHTTP(w, r)
		})))

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeContent(w, r, "robots.txt", time.Now(), strings.NewReader("User-agent: *\nAllow: /"))
	})
}
