package tasks

import (
	"context"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	"github.com/charmbracelet/huh"
)

const priorityManagementDescription = `You are tasked with managing issue priorities.

A critical production issue has been reported. You need to rebalance the current sprint priorities:

1. Review all current issues and their priorities
2. Identify the most urgent production issue
3. Reprioritize existing work to accommodate the urgent fix
4. Defer lower priority items if necessary
5. Update the team on priority changes via comments
6. Ensure the critical path is clear for the urgent fix

The production database is experiencing intermittent connection failures affecting all users.`

type PriorityManagementTask struct {
	done       bool
	app        *service.App
	setupIssue *models.Issue
}

func NewPriorityManagementTask(app *service.App) *PriorityManagementTask {
	return &PriorityManagementTask{app: app, done: false}
}

func (t *PriorityManagementTask) Config() task.Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/priority-task-stats.json")
}

func (t *PriorityManagementTask) Details() taskui.TaskDetails {
	return BaseDetails().
		WithTitle("Priority Management Task").
		WithDescription(priorityManagementDescription).
		WithTimeToComplete("8m").
		WithDifficulty("Easy")
}

func (t *PriorityManagementTask) Questions(interfaceType task.InterfaceType) taskui.Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("How many issues did you reprioritize or comment on?").
				Options(
					huh.NewOption("1-2", 1),
					huh.NewOption("3-4", 2),
					huh.NewOption("5+", 3),
				),
		),
	)
}

func (t *PriorityManagementTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	priorityIssues := []*models.Issue{
		models.NewIssueBuilder().
			WithTitle("Database connection failures").
			WithDescription("PRODUCTION CRITICAL: Intermittent DB connection failures affecting all users. Needs immediate attention.").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("UI theme updates").
			WithDescription("Update color scheme per new brand guidelines. Currently in progress but can wait.").
			WithPriority(1).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Feature: Dark mode").
			WithDescription("Add dark mode toggle to settings. Nice to have, can be deferred.").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("API rate limiting").
			WithDescription("Add rate limiting to public API endpoints. Security enhancement.").
			WithPriority(2).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Documentation updates").
			WithDescription("Update API documentation for v2 endpoints. Can be deferred.").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, priorityIssues, ""); err != nil {
		return err
	}

	t.setupIssue = models.NewBaseIssue().
		WithTitle("Priority Rebalancing").
		WithDescription(priorityManagementDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *PriorityManagementTask) Validate(ctx context.Context) (bool, error) {
	return EndTaskWithTimeout(&t.done, "Priority management task completed!", 5*time.Second)
}
