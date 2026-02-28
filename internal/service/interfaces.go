package service

import (
	"context"
	"log/slog"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/steveyegge/beads"
)

type App struct {
	Config Config
	Issues IssueService
	Stats  StatsService

	Logger *slog.Logger

	CurrentFeedback *ValidationFeedback
}

type IssueService interface {
	beads.Storage
	AllIssues(ctx context.Context) ([]models.Issue, error)
	DeleteIssues() error
}

type StatsService interface {
	Load(ctx context.Context) error
	Save(ctx context.Context) error
	GetStatistics() (models.Statistics, error)
}

type ValidationFeedback struct {
	Success bool
	Message string
}
