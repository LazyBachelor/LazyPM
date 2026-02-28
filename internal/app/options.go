package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/LazyBachelor/LazyPM/internal/models"
)

type Option func(*AppBuilder) error

type AppBuilder struct {
	config Config
	ctx    context.Context

	logger      *slog.Logger
	lifecycle   *Lifecycle
	initializer Initializer

	issueService models.IssueService
	statsService models.StatsService
}

func defaultBuilder(ctx context.Context, config Config) *AppBuilder {
	return &AppBuilder{
		ctx:         ctx,
		config:      config,
		lifecycle:   NewLifecycle(),
		initializer: &InteractiveInitializer{},
		logger:      slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(b *AppBuilder) error {
		b.logger = logger
		return nil
	}
}

func WithInitializer(initializer Initializer) Option {
	return func(b *AppBuilder) error {
		b.initializer = initializer
		return nil
	}
}

func WithLifecycle(l *Lifecycle) Option {
	return func(b *AppBuilder) error {
		b.lifecycle = l
		return nil
	}
}

func WithIssueService(svc models.IssueService) Option {
	return func(b *AppBuilder) error {
		b.issueService = svc
		return nil
	}
}

func WithStatsService(svc models.StatsService) Option {
	return func(b *AppBuilder) error {
		b.statsService = svc
		return nil
	}
}
