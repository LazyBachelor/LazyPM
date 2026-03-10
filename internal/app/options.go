package app

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

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
	lifecycle := NewLifecycle()
	logger := defaultLogger(config, lifecycle)

	return &AppBuilder{
		ctx:         ctx,
		config:      config,
		lifecycle:   lifecycle,
		initializer: &InteractiveInitializer{},
		logger:      logger,
	}
}

func defaultLogger(config Config, lifecycle *Lifecycle) *slog.Logger {
	statsDir := filepath.Dir(config.StatisticsStoragePath)
	if statsDir == "" {
		statsDir = "."
	}

	if err := os.MkdirAll(statsDir, 0o755); err != nil {
		return slog.New(slog.NewJSONHandler(io.Discard, nil))
	}

	logPath := filepath.Join(statsDir, "app-logs.jsonl")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return slog.New(slog.NewJSONHandler(io.Discard, nil))
	}

	if lifecycle != nil {
		lifecycle.Add(func() {
			_ = logFile.Close()
		})
	}

	return slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{}))
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
