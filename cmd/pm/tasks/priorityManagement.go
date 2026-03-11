package tasks

import (
	"context"
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
	"github.com/charmbracelet/huh"
)

const priorityManagementDescription = `You are tasked with managing issue priorities.

A critical production issue has been reported.

The database is not working properly and users are not able to connect and access their data.

You need to rebalance the current sprint priorities:

1. Assign the task Issue you are currently reading to yourself as "Me" and set status to "In Progress".
2. A new issue has appeared in the list that needs urgent attention. Change the database related issue's priority to 4 (critical).
3. Set the priority of the feature and chore issues in the list to 1 (low).
4. Change this issue status to "Closed".

`

type PriorityManagementTask struct {
	done       bool
	app        *App
	setupIssue *Issue
}

func NewPriorityManagementTask(app *App) *PriorityManagementTask {
	return &PriorityManagementTask{app: app, done: false}
}

func (t *PriorityManagementTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/priority-task-stats.json")
}

func (t *PriorityManagementTask) Details() TaskDetails {
	return BaseDetails().
		WithTitle("Priority Management Task").
		WithDescription(priorityManagementDescription).
		WithTimeToComplete("8m").
		WithDifficulty("Easy")
}

func (t *PriorityManagementTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).
		With(
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
		NewIssueBuilder().
			WithTitle("Database connection failures").
			WithDescription("PRODUCTION CRITICAL: Intermittent DB connection failures affecting all users. Needs immediate attention.").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("UI theme updates").
			WithDescription("Update color scheme per new brand guidelines. Currently in progress but can wait.").
			WithPriority(2).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeFeature).
			Build(),
		NewIssueBuilder().
			WithTitle("Feature: Dark mode").
			WithDescription("Add dark mode toggle to settings. Nice to have, can be deferred.").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeFeature).
			Build(),
		NewIssueBuilder().
			WithTitle("API rate limiting").
			WithDescription("Add rate limiting to public API endpoints. Security enhancement.").
			WithPriority(2).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeFeature).
			Build(),
		NewIssueBuilder().
			WithTitle("Documentation updates").
			WithDescription("Update API documentation for v2 endpoints. Can be deferred.").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeChore).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, priorityIssues, ""); err != nil {
		return err
	}

	t.setupIssue = NewIssueBuilder().
		WithTitle("Priority Rebalancing").
		WithDescription(priorityManagementDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *PriorityManagementTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app, t.setupIssue)
	if err != nil {
		return expect.ValidationFeedback
	}

	taskIssue := t.setupIssue

	if taskIssue.Assignee != "Me" {
		expect.Assert(taskIssue.Assignee == "Me",
			fmt.Sprintf("'%s' not assigned to you, but to '%s'", taskIssue.Title, taskIssue.Assignee))
		return expect.ValidationFeedback
	}

	if taskIssue.Status == models.StatusOpen {
		expect.Assert(taskIssue.Status == models.StatusInProgress,
			fmt.Sprintf("'%s' status should be 'In Progress', but was '%s'", taskIssue.Title, taskIssue.Status))
		return expect.ValidationFeedback
	}

	for ii := range issues {

		if issues[ii].Title == "Database connection failures" {
			if issues[ii].Priority != 4 {
				expect.Assert(issues[ii].Priority == 4,
					fmt.Sprintf("'%s' priority should be 4 (critical), but was '%d'", issues[ii].Title, issues[ii].Priority))
				return expect.ValidationFeedback
			}
		} else {
			if issues[ii].Priority != 1 {
				expect.Assert(issues[ii].Priority == 1,
					fmt.Sprintf("'%s' priority should be 1 (low), but was '%d'", issues[ii].Title, issues[ii].Priority))
				return expect.ValidationFeedback
			}
		}
	}
	expect.Assert(taskIssue.Status == models.StatusClosed,
		fmt.Sprintf("'%s' is not set to closed", taskIssue.Title))



	return expect.Complete()
}
