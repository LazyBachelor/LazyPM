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
1. Find the 2 issues that depend on 'Setup database connection'.
2. Add the 'Setup database connection' issue as a dependency to the 2 issues.
3. Set the 2 issues' status to "Blocked".`

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
			WithPriority(3).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Implement Authentication System").
			WithDescription("Add login/logout functionality.").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Add user management operations").
			WithDescription("Add operations for user management.").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, t.depIssues, ""); err != nil {
		return err
	}



	return nil
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
					fmt.Sprintf("%s status", issue.Title,))
			}
		}
	}

	for _, depIssue := range t.depIssues[1:] {

		deps, _ := t.app.Issues.GetDependencies(ctx, depIssue.ID)
		expect.Equal(len(deps), 1,
			fmt.Sprintf("%s - %d dependencies", depIssue.Title, len(deps)))
		if len(deps) == 1 {
			expect.Equal(deps[0].ID, t.depIssues[0].ID,
				fmt.Sprintf("%s dependency", depIssue.Title))
		}
	
	}




	if !expect.Valid() {
		return expect.ValidationFeedback
	}

	return expect.Complete()
}
