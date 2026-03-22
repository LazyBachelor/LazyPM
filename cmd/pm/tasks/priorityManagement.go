package tasks

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const priorityManagementDescription = `You are tasked with managing issue priorities.

A critical production issue has been reported.
The database is not working properly and users are not able to connect and access their data.

You need to rebalance the current sprint priorities:

1. A new issue has appeared in the list that needs urgent attention:
   - Change the database related issue's priority to 4.
2. Set the priority of the feature and chore issues in the list to 1.`

type PriorityManagementTask struct {
	done           bool
	app            *App
	priorityIssues []*models.Issue
	isInProgress   bool
}

func NewPriorityManagementTask(app *App) *PriorityManagementTask {
	return &PriorityManagementTask{app: app, done: false}
}

func (t *PriorityManagementTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/priority-task-stats.json")
}

func (t *PriorityManagementTask) Details(interfaceType InterfaceType) TaskDetails {
	return BaseDetails(interfaceType).
		WithTitle("Priority Management Task").
		WithDescription(priorityManagementDescription).
		WithTimeToComplete("4m").
		WithDifficulty("Medium")
}

func (t *PriorityManagementTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().Key("management_difficulty").
				Title("How did it feel managing multiple issues?").
				Description("Managing multiple issues would require careful prioritization and coordination.").
				Options(
					huh.NewOption("Very Easy", 1),
					huh.NewOption("Easy", 2),
					huh.NewOption("Moderate", 3),
					huh.NewOption("Difficult", 4),
					huh.NewOption("Very Difficult", 5),
				),
		),
	)
}

func (t *PriorityManagementTask) QuestionnaireKeys(_ InterfaceType) []string {
	return BaseKeys().With("management_difficulty")
}

func (t *PriorityManagementTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	t.priorityIssues = []*models.Issue{
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

	return t.app.Issues.CreateIssues(ctx, t.priorityIssues, "")
}

func (t *PriorityManagementTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app)
	if err != nil {
		return expect.Fatal("Failed to fetch issues for validation")
	}

	for _, issue := range issues {
		if issue.Title == t.priorityIssues[0].Title {
			expect.Equal(issue.Priority, 4,
				fmt.Sprintf("Priority of issue %s", issue.Title))
		} else {
			expect.Equal(issue.Priority, 1,
				fmt.Sprintf("Priority of issue %s", issue.Title))
		}
	}

	return expect.Complete()
}
