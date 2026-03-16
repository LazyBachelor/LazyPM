package tasks

import (
	"context"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
	"github.com/charmbracelet/huh"
)

const backlogRefinementDescription = `You are tasked with backlog refinement.

The product backlog has become cluttered with old and unclear issues. You need to groom the backlog:

1. Go to the backlog.
2. Find two issues that got the same name or describe the same problem.
3. Open one of these issues.
4. Select "Close issue"
5. Choose "Duplicate issue" as closing reason.
6. Save/close issue.

Focus on making the backlog a reliable source of upcoming work.`

type BacklogRefinementTask struct {
	done       bool
	app        *App
	setupIssue *Issue
}

func NewBacklogRefinementTask(app *App) *BacklogRefinementTask {
	return &BacklogRefinementTask{app: app, done: false}
}

func (t *BacklogRefinementTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/refinement-task-stats.json")
}

func (t *BacklogRefinementTask) Details(interfaceType InterfaceType) TaskDetails {
	return BaseDetails(interfaceType).
		WithTitle("Backlog Refinement Task").
		WithDescription(backlogRefinementDescription).
		WithTimeToComplete("2m").
		WithDifficulty("Easy")
}

func (t *BacklogRefinementTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Key("duplicate-issues-closed").
				Title("How many issues did you close during refinement?").
				Description("There was only one issue required for the task; if you closed more, why?").
				Options(
					huh.NewOption("1", 1),
					huh.NewOption("1+ I saw there were more to be closed", 2),
					huh.NewOption("1+ I closed more by mistake", 3),
					huh.NewOption("0 I didn't close any", 4),
				),
		),
	)
}

func (t *BacklogRefinementTask) QuestionnaireKeys(_ InterfaceType) []string {
	return BaseKeys().With("duplicate-issues-closed")
}

func (t *BacklogRefinementTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	refinementIssues := []*models.Issue{
		NewIssueBuilder().
			WithTitle("User profile page").
			WithDescription("Create page for users to view and edit their profile").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("User profile page").
			WithDescription("Allow users to view their profile information").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Fix login timeout").
			WithDescription("Login sometimes times out after 30 seconds").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeBug).
			Build(),
		NewIssueBuilder().
			WithTitle("Fix login timeout").
			WithDescription("Users report login requests timing out").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeBug).
			Build(),
		NewIssueBuilder().
			WithTitle("Mobile app redesign").
			WithDescription("Redesign mobile interface with modern UI patterns").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Customer feedback system").
			WithDescription("Build system for collecting user feedback").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, refinementIssues, ""); err != nil {
		return err
	}

	t.setupIssue = NewIssueBuilder().
		WithTitle("Backlog Refinement Session").
		WithDescription(backlogRefinementDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *BacklogRefinementTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app, t.setupIssue)
	if err != nil {
		return expect.ValidationFeedback
	}

	var closedDuplicate *models.Issue
	for _, issue := range issues {
		if issue.Status == models.StatusClosed &&
			strings.Contains(strings.ToLower(issue.CloseReason), "duplicate") {
			closedDuplicate = issue
			break
		}
	}

	expect.Assert(closedDuplicate != nil,
		"Expected one duplicate issue to be closed with 'Duplicate issue' as closing reason")

	return expect.Complete()
}
