package main

import (
	"beadstest/cmd/web/routes"
	"beadstest/cmd/web/server"
	"beadstest/internal/models"
	"beadstest/internal/service"
	"beadstest/internal/storage"
	"context"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/steveyegge/beads"
)

//go:embed assets/*
var assets embed.FS

func main() {
	ctx := context.Background()

	beadStore, err := beads.NewSQLiteStorage(ctx, "./.pm/db.db")
	handleError(err, "Error initializing Beads storage")
	defer beadStore.Close()

	beadSvc, err := service.NewService(ctx, beadStore, "pm")
	handleError(err, "Error initializing Beads service")
	defer beadSvc.Close()

	statStore := storage.NewJsonStorage("./.pm/stats.json", &models.Statistics{
		ID:        uuid.New(),
		StartTime: time.Now(),
	})

	statSvc, err := service.NewStatisticsService(statStore)
	handleError(err, "Error initializing Statistics service")

	server := server.NewServer(server.Server{
		Port:              8080,
		Assets:            assets,
		Service:           beadSvc,
		StatisticsService: statSvc,
		Routes: []server.Route{
			{Pattern: "/", Component: routes.Index()},
		},
	})

	fmt.Printf("Starting web server on port %s...\n", server.Addr)
	handleError(server.ListenAndServe(), "Server closed")
}

func handleError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
		os.Exit(1)
	}
}
