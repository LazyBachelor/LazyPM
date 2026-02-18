package server

import (
	"embed"
	"net/http"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/service"
)

type Server struct {
	Address  string
	Assets   embed.FS
	Services *service.Services
}

// NewServer creates and configures a new HTTP server instance.
func NewServer(props Server) *http.Server {
	if props.Address == "" {
		props.Address = ":8080"
	}

	handler := props.RegisterRoutes(props.Assets)

	return &http.Server{
		Addr:         props.Address,
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
}
