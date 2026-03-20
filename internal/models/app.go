package models

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type App struct {
	Config Config
	Issues IssueService
	Stats  StatsService

	Logger *slog.Logger

	Tasks      map[string]Tasker
	Interfaces map[string]Interface

	CurrentFeedback *ValidationFeedback
	ActionLogger    func(string)

	SubmitChan chan<- struct{}
}

func (a *App) LogAction(action string) {
	if a == nil || action == "" {
		return
	}
	if a.Logger != nil {
		a.Logger.Info("action event", "action", action)
	}
	if a.ActionLogger != nil {
		a.ActionLogger(action)
	}
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

	AddDependency(ctx context.Context, dep *Dependency, actor string) error
	RemoveDependency(ctx context.Context, issueID, dependsOnID string, actor string) error
	GetDependencies(ctx context.Context, issueID string) ([]*Issue, error)
	GetDependents(ctx context.Context, issueID string) ([]*Issue, error)
}

type StatsService interface {
	Load(ctx context.Context) error
	Save(ctx context.Context) error
	GetStatistics() (Statistics, error)
	GetParticipantID() bson.ObjectID
	RecordTaskRun(ctx context.Context, run TaskRunMetrics) error
	RecordIntroQuestionnaireAnswers(answers map[string]any) error
}
