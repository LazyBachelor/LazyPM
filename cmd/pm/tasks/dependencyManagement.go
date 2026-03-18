package tasks

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const dependencyManagementDescription = `You are tasked with managing issue dependencies.

Several issues in your project have dependencies on other issues. You need to:

1. Find the 4 issues that mention dependencies in their detail description. For example: "Depends on Issue '123'". Set their status to "blocked".
2. Find the 2 foundational issues that are mentioned by the other issues.
3. Set priority of the 2 foundational issues to 3 (high).
4. Set status of the 2 foundational issues to in-progress.
5. Assign the 2 foundational issues to yourself as "Me".

Resolving dependencies in the right order is critical for efficient team workflow.`

type DependencyManagementTask struct {
	done       bool
	app        *App
	setupIssue *Issue
	depIssues  []*Issue
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
			WithTitle("Create home page for the website").
			WithDescription("Create a page for the website.").
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
		NewIssueBuilder().
			WithTitle("Create user profile page").
			WithDescription("Frontend user profile page. Depends on 'Create home page for the website' issue.").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Create about page").
			WithDescription("Frontend about page. Depends on 'Create home page for the website' issue.").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, t.depIssues, ""); err != nil {
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

	taskIssue := t.setupIssue
	issues, err := FetchIssues(ctx, t.app, t.setupIssue)
	if err != nil {
		return expect.Fatal("Could not fetch issues")
	}

	for _, issue := range issues {
		for _, depIssue := range t.depIssues[2:] {
			if issue.Title == depIssue.Title {
				expect.Equal(issue.Status, models.StatusBlocked,
					fmt.Sprintf("%s status", issue.Title))
			}
		}

		for _, foundationalIssue := range t.depIssues[:2] {
			if issue.Title == foundationalIssue.Title {
				expect.Equal(issue.Priority, 3,
					fmt.Sprintf("%s priority", issue.Title))
				expect.Equal(issue.Assignee, "Me",
					fmt.Sprintf("%s assignee", issue.Title))
				expect.Equal(issue.Status, models.StatusInProgress,
					fmt.Sprintf("%s status", issue.Title))
			}
		}

	}

	if !expect.Valid() {
		return expect.ValidationFeedback
	}

	expect.Equal(taskIssue.Status, models.StatusClosed,
		fmt.Sprintf("%s", taskIssue.Title))

	return expect.Complete()
}
