package app

import (
	"context"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/storage"
	"github.com/steveyegge/beads"
)

type App = models.App
type Config = models.Config

func New(ctx context.Context, config Config, opts ...Option) (*App, func(), error) {
	b := defaultBuilder(ctx, config)

	for _, opt := range opts {
		if err := opt(b); err != nil {
			return nil, nil, err
		}
	}

	if !config.AutoInit {
		if err := b.initializer.Init(config.BeadsDBPath); err != nil {
			return nil, nil, err
		}
	}

	if b.issueService == nil {
		sqliteStore, err := beads.NewSQLiteStorage(b.ctx, config.BeadsDBPath)
		if err != nil {
			return nil, nil, err
		}
		b.lifecycle.Add(func() { sqliteStore.Close() })

		b.issueService, err = storage.NewBeadsIssueStorage(b.ctx, sqliteStore, config.IssuePrefix)
		if err != nil {
			return nil, nil, err
		}
	}

	if b.statsService == nil {
		statStore := storage.NewJsonStorage(config.StatisticsStoragePath, &models.Statistics{
			ID:        0,
			StartTime: time.Now(),
		})

		jsonStatService, err := NewStatisticsService(statStore)
		if err != nil {
			return nil, nil, err
		}

		b.statsService = jsonStatService
	}

	app := &App{
		Config: config,
		Logger: b.logger,

		Issues: b.issueService,
		Stats:  b.statsService,
	}

	return app, b.lifecycle.Close, nil
}
