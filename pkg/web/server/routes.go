package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
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
	r.Get("/status/modal", handler.HandleTaskStatusModal)

	r.Route("/issues", func(r chi.Router) {
		r.Get("/", handler.ListIssues)
		r.Post("/", handler.CreateIssue)
		r.Get("/create", handler.CreateIssueFormModal)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(handler.IssueCtx)
			r.Get("/", handler.GetIssue)
			r.Patch("/", handler.UpdateIssue)
			r.Get("/edit", handler.EditIssueFormModal)
			r.Get("/assignee", handler.AssigneeFormModal)
			r.Patch("/assignee", handler.UpdateAssignee)
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
	var fileServer http.Handler

	if os.Getenv("DEV") == "True" {
		log.Println("Running in development mode: serving assets from disk")
		fileServer = http.FileServer(http.Dir("pkg/web/assets"))
	} else {
		subFS, err := fs.Sub(assets, "assets")
		if err != nil {
			log.Fatalf("failed to create sub filesystem: %v", err)
		}
		fileServer = http.FileServer(http.FS(subFS))
	}

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
