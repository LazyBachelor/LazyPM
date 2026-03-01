package models

import (
	"context"
	"log/slog"
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

// IssueService defines the interface for managing issues in the application.
// This a siplification of the Beads issue storage interface, tailored to our needs.
type IssueService interface {
	CreateIssue(ctx context.Context, issue *Issue, actor string) error
	CreateIssues(ctx context.Context, issues []*Issue, actor string) error

	UpdateIssue(ctx context.Context, id string, updates map[string]any, actor string) error

	CloseIssue(ctx context.Context, id string, reason string, actor string, session string) error
	DeleteIssue(ctx context.Context, id string) error
	DeleteIssues() error

	GetIssue(ctx context.Context, id string) (*Issue, error)
	SearchIssues(ctx context.Context, query string, filter IssueFilter) ([]*Issue, error)

	AddIssueComment(ctx context.Context, issueID, author, text string) (*Comment, error)
	GetIssueComments(ctx context.Context, issueID string) ([]*Comment, error)
	GetCommentCounts(ctx context.Context, issueIDs []string) (map[string]int, error)
}

type StatsService interface {
	Load(ctx context.Context) error
	Save(ctx context.Context) error
	GetStatistics() (Statistics, error)
}
