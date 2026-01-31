package server

import (
	"beadstest/internal/service"
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/steveyegge/beads"
)

type Server struct {
	Port              int
	Assets            embed.FS
	Routes            []Route
	Service           beads.Storage
	StatisticsService *service.StatisticsService
}

// NewServer creates and configures a new HTTP server instance.
func NewServer(props Server) *http.Server {
	if props.Port == 0 {
		props.Port = 8080
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", props.Port),
		Handler:      props.RegisterRoutes(props.Assets, props.Routes),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
}
