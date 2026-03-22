package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const issueReviewCleanupDescription = `You are responsible for reviewing and maintaining the current project issues.

Complete the following steps:
1. Add a comment to two issues
2. Delete the issue titled "Delete this issue"`

type IssueReviewCleanupTask struct {
	done         bool
	app          *App
	reviewIssues []*models.Issue
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
		WithTimeToComplete("3m").
		WithDifficulty("Medium")
}

func (t *IssueReviewCleanupTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType)
}

func (t *IssueReviewCleanupTask) QuestionnaireKeys(_ InterfaceType) []string {
	return BaseKeys()
}

func (t *IssueReviewCleanupTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	t.reviewIssues = []*models.Issue{
		models.NewIssueBuilder().
			WithTitle("Delete this issue").
			WithDescription("").
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

	return t.app.Issues.CreateIssues(ctx, t.reviewIssues, "")
}

func (t *IssueReviewCleanupTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app)
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
	expect.Equal(len(issues), len(t.reviewIssues)-1, "Remaining issues after cleanup")

	issue, err := t.app.Issues.GetIssue(ctx, t.reviewIssues[0].ID)
	if err != nil {
		return expect.Fatal("Deleted issue should not be found")
	}

	if issue != nil {
		expect.Fail("The issue titled 'Delete this issue' should have been deleted.")
	} else {
		expect.Pass("Issue deleted successfully")
	}

	return expect.Complete()
}
