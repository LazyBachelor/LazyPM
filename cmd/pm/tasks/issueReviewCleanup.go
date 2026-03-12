package tasks

import (
	"context"
	"fmt"

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

func (t *IssueReviewCleanupTask) Details() TaskDetails {
	return BaseDetails().
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
		return expect.Fail("Could not fetch issues: " + err.Error()).Complete()
	}

	// Step 1 cannot be verified (opening/viewing leaves no trace) — not shown in results

	// Step 2: Add comment to 2 issues
	var issueIDs []string
	for _, issue := range issues {
		issueIDs = append(issueIDs, issue.ID)
	}
	issuesWithComments := 0
	if len(issueIDs) > 0 {
		counts, err := t.app.Issues.GetCommentCounts(ctx, issueIDs)
		if err != nil {
			expect.Fail("Step 2: Could not verify comments: " + err.Error())
		} else {
			for _, count := range counts {
				if count >= 1 {
					issuesWithComments++
				}
			}
			if issuesWithComments >= 2 {
				expect.Pass(fmt.Sprintf("Step 2: Added comments to at least 2 issues (%d issues have comments)", issuesWithComments))
			} else {
				expect.Fail(fmt.Sprintf("Step 2: Add comments to at least 2 issues (found %d with comments)", issuesWithComments))
			}
		}
	} else {
		expect.Fail("Step 2: Add comments to at least 2 issues (no issues to check)")
	}

	// Step 3: Delete the cleanup task issue — we verify it's gone (Step 4 is same outcome; we can't verify they "looked")
	setupStillExists, _ := t.app.Issues.GetIssue(ctx, t.setupIssue.ID)
	deletedCleanupIssue := setupStillExists == nil

	if len(issues) == 5 && deletedCleanupIssue {
		expect.Pass("Step 3: Deleted the cleanup task issue — it no longer appears in the list")
	} else if len(issues) < 5 {
		expect.Fail(fmt.Sprintf("Step 3: Delete only the cleanup task issue — you deleted a project issue (expected 5, got %d)", len(issues)))
	} else if !deletedCleanupIssue {
		expect.Fail("Step 3: Delete the cleanup task issue (\"Issue Review and Cleanup Task\") from the list")
	}

	return expect.Complete()
}
