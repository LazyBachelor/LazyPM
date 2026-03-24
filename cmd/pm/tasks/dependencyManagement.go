package tasks

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const dependencyManagementDescription = `You are tasked with managing issue dependencies.

Several issues in your project have dependencies on other issues. 

You need to:
1. Find 2 issues that mention dependencies in their description.
   - Set their status to "Blocked".
2. Find the issue that is mentioned by the other issues:
   - Set priority to 3.
   - Set status to in-progress.
   - Assign to yourself as "Me".`

type DependencyManagementTask struct {
	done      bool
	app       *App
	depIssues []*Issue
}

func NewDependencyManagementTask(app *App) *DependencyManagementTask {
	return &DependencyManagementTask{app: app, done: false}
}

func (t *DependencyManagementTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/dependency-task-stats.json")
}

func (t *DependencyManagementTask) Details(interfaceType InterfaceType) TaskDetails {
	return BaseDetails(interfaceType).
		WithTitle("Dependency Management Task").
		WithDescription(dependencyManagementDescription).
		WithTimeToComplete("3m").
		WithDifficulty("Medium")
}

func (t *DependencyManagementTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().Key("discovery_difficulty").
				Title("How difficult was it discovering the dependencies?").
				Description("We are interested in how difficult it was to discover the issues that needed to be addressed.").
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

func (t *DependencyManagementTask) QuestionnaireKeys(_ InterfaceType) []string {
	return BaseKeys().With("discovery_difficulty")
}

func (t *DependencyManagementTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	t.depIssues = []*Issue{
		NewIssueBuilder().
			WithTitle("Setup database connection").
			WithDescription("Configure database connection pool.").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Implement Authentication System").
			WithDescription("Add login/logout functionality. Depends on 'Setup database connection' issue.").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Add user management operations").
			WithDescription("Add operations for user management. Depends on 'Setup database connection' issue.").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
	}

	return t.app.Issues.CreateIssues(ctx, t.depIssues, "")
}

func (t *DependencyManagementTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app)
	if err != nil {
		return expect.Fatal("Could not fetch issues")
	}

	for _, issue := range issues {
		for _, depIssue := range t.depIssues[1:] {
			if issue.Title == depIssue.Title {
				expect.Equal(issue.Status, models.StatusBlocked,
					fmt.Sprintf("%s status", issue.Title))
			}
		}
	}
	if !expect.Valid() {
		return expect.ValidationFeedback
	}
	for _, issue := range issues {
		for _, foundationalIssue := range t.depIssues[:1] {
			if issue.Title == foundationalIssue.Title {
				expect.Equal(issue.Priority, 3,
					fmt.Sprintf("%s priority", issue.Title))
				expect.NotEmptyAndEqual(issue.Assignee, "Me",
					fmt.Sprintf("%s assignee", issue.Title))
				expect.Equal(issue.Status, models.StatusInProgress,
					fmt.Sprintf("%s status", issue.Title))
			}
		}
	}

	return expect.Complete()
}
