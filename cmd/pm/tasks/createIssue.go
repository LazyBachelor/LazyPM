package tasks

import (
	"context"
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const description = `You are tasked with creating a new issue in the project management system.

This task will test your ability to use the issue creation workflow effectively.

Your task:
1. Create a new issue with the title "My first Issue"
2. Add this detailed description "I need to do some coding"
3. Assign the issue to yourself as "Me"
4. Mark the issue as In Progress when you are done.

Make sure to fill out all the necessary details to help others understand the work item.`

type CreateIssueTask struct {
	done       bool
	app        *App
	setupIssue *Issue
}

func NewCreateIssueTask(app *App) *CreateIssueTask {
	return &CreateIssueTask{app: app, done: false}
}

func (t *CreateIssueTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/create-issue-stats.json")
}

func (t *CreateIssueTask) Details() TaskDetails {
	return BaseDetails().WithTitle("Create Issue Task").WithDescription(description)
}

func (t *CreateIssueTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType)
}

func (t *CreateIssueTask) QuestionnaireKeys(_ InterfaceType) []string {
	return []string{"task_completed", "task_difficulty"}
}

func (t *CreateIssueTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	t.setupIssue = NewIssueBuilder().
		WithTitle("Create a New Issue").
		WithDescription(description).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *CreateIssueTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app, t.setupIssue)
	if err != nil {
		return expect.ValidationFeedback
	}

	if len(issues) == 0 {
		expect.Fail("No new issues created")
		return expect.ValidationFeedback
	}

	issue := issues[0]

	expect.Assert(len(issues) < 2, "Multiple issues were created instead of one")

	expect.NotEmptyAndEqual(issue.Title, "My first Issue", "Issue title")
	expect.NotEmptyAndEqual(issue.Description, "I need to do some coding", "Issue description")
	expect.NotEmptyAndEqual(issue.Assignee, "Me", "Issue assignee")

	expect.Assert(issue.Status == models.StatusInProgress,
		fmt.Sprintf("Issue status should be 'In Progress', but was '%s'", issue.Status))

	return expect.Complete()
}
