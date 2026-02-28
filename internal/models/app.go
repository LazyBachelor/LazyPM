package models

import (
	"context"
	"log/slog"

	"github.com/steveyegge/beads"
)

type App struct {
	Config Config
	Issues IssueService
	Stats  StatsService

	Logger *slog.Logger

	Tasks      map[string]Tasker
	Interfaces map[string]Interface

	CurrentFeedback *ValidationFeedback
}

type IssueService interface {
	beads.Storage // Want to get rid of this dependency, but it provides a lot of useful methods that would be a pain to re-implement right now
	DeleteIssues() error
}

type StatsService interface {
	Load(ctx context.Context) error
	Save(ctx context.Context) error
	GetStatistics() (Statistics, error)
}
