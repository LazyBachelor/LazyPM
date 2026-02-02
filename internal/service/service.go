package service

import (
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/storage"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/steveyegge/beads"
)

type Config struct {
	WebAddress            string
	BeadsDBPath           string
	IssuePrefix           string
	StatisticsStoragePath string
}

type Services struct {
	Config     Config
	Beads      *BeadsService
	Statistics *StatisticsService
}

func NewServices(ctx context.Context, config Config) (*Services, func(), error) {
	var cleanupFuncs []func()

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

	statSvc, err := NewStatisticsService(statStore)
	if err != nil {
		return nil, nil, err
	}

	return &Services{
		Beads:      beadsSvc,
		Statistics: statSvc,
		Config:     config,
	}, func() { runCleanup(cleanupFuncs) }, nil

}

func runCleanup(funcs []func()) {
	for _, fn := range funcs {
		fn()
	}
}
