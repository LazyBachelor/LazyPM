package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/internal/utils"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	"github.com/charmbracelet/huh"
)

const backlogRefinementDescription = `You are tasked with backlog refinement.

The product backlog has become cluttered with old and unclear issues. You need to groom the backlog:

1. Review all issues in the backlog
2. Identify stale or obsolete issues (older items that are no longer relevant)
3. Update issue descriptions for clarity where needed
4. Close issues that are duplicates or no longer applicable
5. Reprioritize issues based on current business value
6. Ensure remaining issues are well-defined and actionable

Focus on making the backlog a reliable source of upcoming work.`

type BacklogRefinementTask struct {
	done       bool
	app        *service.App
	setupIssue *models.Issue
}

func NewBacklogRefinementTask(app *service.App) *BacklogRefinementTask {
	return &BacklogRefinementTask{app: app, done: false}
}

func (t *BacklogRefinementTask) Config() task.Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/refinement-task-stats.json")
}

func (t *BacklogRefinementTask) Details() taskui.TaskDetails {
	return BaseDetails().
		WithTitle("Backlog Refinement Task").
		WithDescription(backlogRefinementDescription).
		WithTimeToComplete("12m").
		WithDifficulty("Medium")
}

func (t *BacklogRefinementTask) Questions(interfaceType task.InterfaceType) taskui.Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("How many issues did you close or update during refinement?").
				Options(
					huh.NewOption("1-2", 1),
					huh.NewOption("3-4", 2),
					huh.NewOption("5+", 3),
				),
		),
	)
}

func (t *BacklogRefinementTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	refinementIssues := []*models.Issue{
		models.NewIssueBuilder().
			WithTitle("Old feature request: Fax integration").
			WithDescription("Allow sending reports via fax. DEPRECATED - nobody uses fax anymore").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("User profile page").
			WithDescription("Create page for users to view profile. DUPLICATE of user-management epic").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Mobile app redesign").
			WithDescription("Redesign mobile interface with modern UI patterns. Still relevant, needs clarity").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Legacy data export tool").
			WithDescription("Tool for exporting data in old format. OBSOLETE - format no longer supported").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("API v1 documentation").
			WithDescription("Document old API version. DEPRECATED - migrating to v2").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Customer feedback system").
			WithDescription("Build system for collecting user feedback. HIGH VALUE - prioritize").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, refinementIssues, ""); err != nil {
		return err
	}

	t.setupIssue = models.NewBaseIssue().
		WithTitle("Backlog Refinement Session").
		WithDescription(backlogRefinementDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *BacklogRefinementTask) Validate(ctx context.Context) ValidationFeedback {
	expect := utils.NewExpector()

	return expect.Complete()
}
