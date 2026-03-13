package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const issueReviewCleanupDescription = `You are responsible for reviewing and maintaining the current project issues.

Using the system, complete the following steps:

1. Open three different issues and read their titles and descriptions
2. Add a comment to two issues
3. Delete this cleanup task issue ("Issue Review and Cleanup Task") from the issue list — do not delete the other project issues
4. Confirm that the cleanup task issue no longer appears in the list`

type IssueReviewCleanupTask struct {
	done       bool
	app        *App
	setupIssue *models.Issue
}

func NewIssueReviewCleanupTask(app *App) *IssueReviewCleanupTask {
	return &IssueReviewCleanupTask{app: app, done: false}
}

func (t *IssueReviewCleanupTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/issue-review-cleanup-stats.json")
}

func (t *IssueReviewCleanupTask) Details(interfaceType InterfaceType) TaskDetails {
	return BaseDetails(interfaceType).
		WithTitle("Issue Review and Cleanup Task").
		WithDescription(issueReviewCleanupDescription).
		WithTimeToComplete("8m").
		WithDifficulty("Easy")
}

func (t *IssueReviewCleanupTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType)
}

func (t *IssueReviewCleanupTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	reviewIssues := []*models.Issue{
		models.NewIssueBuilder().
			WithTitle("Fix login page layout").
			WithDescription("The login form is misaligned on smaller screens. Needs responsive CSS adjustments.").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Add unit tests for API endpoints").
			WithDescription("Current coverage is low. Add tests for /users and /projects endpoints.").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Update dependencies").
			WithDescription("Several npm packages have security advisories. Run npm audit and update.").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Improve error messages").
			WithDescription("Form validation errors are generic. Make them more helpful for users.").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Document deployment process").
			WithDescription("Add README section for deploying to production environment.").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, reviewIssues, ""); err != nil {
		return err
	}

	t.setupIssue = models.NewBaseIssue().
		WithTitle("Issue Review and Cleanup Task").
		WithDescription(issueReviewCleanupDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *IssueReviewCleanupTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app, t.setupIssue)
	if err != nil {
		return expect.Fatal("Could not fetch issues")
	}

	// Count the amount of comments on each issue
	commentCount := map[string]int{}
	for _, issue := range issues {
		comments, _ := t.app.Issues.GetIssueComments(ctx, issue.ID)
		commentCount[issue.ID] = len(comments)
	}

	// Add a count if there are comment on a issue
	commentsInIssues := 0
	for _, count := range commentCount {
		if count > 0 {
			commentsInIssues++
		}
	}

	expect.Equal(commentsInIssues, 2, "Comments on issues")

	// Fetching the setup issue as as FetchIssues only updates the setup issue if it exists,
	// we can check if it was deleted by seeing if it can be fetched again
	if t.setupIssue, err = t.app.Issues.GetIssue(ctx, t.setupIssue.ID); t.setupIssue != nil {
		expect.Fail("Setup issue still exists")
	} else {
		expect.Pass("Setup issue deleted")
	}

	return expect.Complete()
}
