package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
	"github.com/charmbracelet/huh"
)

const dependencyManagementDescription = `You are tasked with managing issue dependencies.

Several issues in your project have dependencies on other issues. You need to:

1. Review the dependency chain described in issue descriptions
2. Identify issues that are blocked by others
3. Prioritize work on foundational issues (those that unblock others)
4. Update issue statuses to reflect dependency resolution
5. Ensure no circular dependencies exist

Resolving dependencies in the right order is critical for efficient team workflow.`

type DependencyManagementTask struct {
	done       bool
	app        *App
	setupIssue *Issue
}

func NewDependencyManagementTask(app *App) *DependencyManagementTask {
	return &DependencyManagementTask{app: app, done: false}
}

func (t *DependencyManagementTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/dependency-task-stats.json")
}

func (t *DependencyManagementTask) Details() TaskDetails {
	return BaseDetails().
		WithTitle("Dependency Management Task").
		WithDescription(dependencyManagementDescription).
		WithTimeToComplete("15m").
		WithDifficulty("Hard")
}

func (t *DependencyManagementTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("How many foundational issues did you identify?").
				Options(
					huh.NewOption("1", 1),
					huh.NewOption("2", 2),
					huh.NewOption("3+", 3),
				),
		),
	)
}

func (t *DependencyManagementTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	depIssues := []*Issue{
		NewIssueBuilder().
			WithTitle("Setup database connection").
			WithDescription("Configure database connection pool.").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Define API contract").
			WithDescription("Create OpenAPI spec.").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Implement user repository").
			WithDescription("Implement data access layer for user management.").
			WithPriority(2).
			WithStatus(models.StatusBlocked).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Create user endpoints").
			WithDescription("REST API for users.").
			WithPriority(2).
			WithStatus(models.StatusBlocked).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Build user profile UI").
			WithDescription("Frontend user profile page.").
			WithPriority(3).
			WithStatus(models.StatusBlocked).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, depIssues, ""); err != nil {
		return err
	}

	t.setupIssue = NewIssueBuilder().
		WithTitle("Dependency Management").
		WithDescription(dependencyManagementDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *DependencyManagementTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	return expect.Complete()
}
