package tasks

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const backlogRefinementDescription = `You are tasked with backlog refinement.

The product backlog has become cluttered with old and unclear issues.

You need to groom the backlog:
1. Go to the backlog.
2. Find two issues that have the same name.
3. Open one of these issues.
4. Select "Close issue"
5. Set the closing reason to "Duplicate issue".
6. Save/close issue.

Focus on making the backlog a reliable source of upcoming work.`

type BacklogRefinementTask struct {
	done bool
	app  *App
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

	return t.app.Issues.CreateIssues(ctx, refinementIssues, "")
}

func (t *BacklogRefinementTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app)
	if err != nil {
		return expect.Fatal("Could not fetch issues")
	}

	fmt.Printf("DEBUG: Found %d issues\n", len(issues))

	var closedDuplicate *models.Issue
	for _, issue := range issues {
		fmt.Printf("DEBUG: Issue '%s' has status: %v\n", issue.Title, issue.Status)
		if issue.Status == models.StatusClosed {
			closedDuplicate = issue
			break
		}
	}

	fmt.Printf("DEBUG: closedDuplicate is nil: %v\n", closedDuplicate == nil)
	if closedDuplicate == nil {
		expect.Fail("Closed duplicate issue is wrong")
	} else {
		expect.Pass("Closed duplicate issue is correct")
		title := closedDuplicate.Title
		isDuplicate := title == "User profile page" || title == "Fix login timeout"
		expect.Assert(isDuplicate, "Closed issue was one of the duplicates")
	}

	feedback := expect.Complete()
	fmt.Printf("DEBUG: Validation Success=%v, Checks=%d\n", feedback.Success, len(feedback.Checks))
	for i, check := range feedback.Checks {
		fmt.Printf("DEBUG: Check[%d]: Valid=%v, Message=%s\n", i, check.Valid, check.Message)
	}
	return feedback
}
