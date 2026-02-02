package server

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"embed"
	"net/http"
	"time"
)

type Server struct {
	Address  string
	Assets   embed.FS
	Routes   []Route
	Services *service.Services
}

// NewServer creates and configures a new HTTP server instance.
func NewServer(props Server) *http.Server {
	if props.Address == "" {
		props.Address = "localhost:8080"
	}

	return &http.Server{
		Addr:         props.Address,
		Handler:      props.RegisterRoutes(props.Assets, props.Routes),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
}
