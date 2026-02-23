package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/storage"
	"github.com/charmbracelet/huh"

	"github.com/google/uuid"
	"github.com/steveyegge/beads"
)

func NewServices(ctx context.Context, config Config) (*App, func(), error) {
	var cleanupFuncs []func()

	if !initialized(config.BeadsDBPath) {
		fmt.Println("PM is not initialized")
		os.Exit(0)
	}

	store, err := beads.NewSQLiteStorage(ctx, config.BeadsDBPath)
	if err != nil {
		return nil, nil, err
	}
	cleanupFuncs = append(cleanupFuncs, func() { store.Close() })

	beadsSvc, err := NewBeadsService(ctx, store, config.IssuePrefix)
	if err != nil {
		return nil, nil, err
	}
	cleanupFuncs = append(cleanupFuncs, func() { beadsSvc.Close() })

	statStore := storage.NewJsonStorage(config.StatisticsStoragePath, &models.Statistics{
		ID:        uuid.New(),
		StartTime: time.Now(),
	})

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	statSvc, err := NewStatisticsService(statStore)
	if err != nil {
		return nil, nil, err
	}

	return &App{
		Issues: beadsSvc,
		Stats:  statSvc,
		Config: config,
		Logger: logger,
	}, func() { runCleanup(cleanupFuncs) }, nil

}

func runCleanup(funcs []func()) {
	for _, fn := range funcs {
		fn()
	}
}

func initialized(beadsPath string) bool {
	_, err := os.Stat(beadsPath)

	if os.IsNotExist(err) {
		var initialize bool

		huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Title("PM is not initialized in this directory!").
					Description("Do you want to initialize it here?").
					Value(&initialize),
			),
		).WithTheme(huh.ThemeBase16()).WithAccessible(true).Run()

		if !initialize {
			return false
		}
	}
	return true
}
