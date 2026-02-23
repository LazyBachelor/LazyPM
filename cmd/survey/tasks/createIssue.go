package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
)

const description = `You are tasked with creating a new issue in the project management system.
This task will test your ability to use the issue creation workflow effectively.

Assign this task to yourself and start creating a new issue.
Make sure to fill out all the necessary details, including the title and description.
Once you have created the issue, mark it as closed to complete the task.`

type CreateIssueTask struct {
	done      bool
	app       *service.App
	setupTask *models.Issue
}

func NewCreateIssueTask(app *service.App) *CreateIssueTask {
	return &CreateIssueTask{app: app, done: false}
}

func (t *CreateIssueTask) Config() task.Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/create-issue-stats.json")
}

func (t *CreateIssueTask) Details() taskui.TaskDetails {
	return BaseDetails().WithTitle("Create Issue Task").WithDescription(description)
}

func (t *CreateIssueTask) Questions(interfaceType task.InterfaceType) taskui.Questions {
	return BaseQuestions(interfaceType)
}

func (t *CreateIssueTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	t.setupTask = models.NewBaseIssue().
		WithTitle("Create a New Issue").WithDescription(description).Build()

	return t.app.Issues.CreateIssue(ctx, t.setupTask, "")
}

func (t *CreateIssueTask) Validate(ctx context.Context) (bool, error) {
	issues, err := t.app.Issues.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return false, err
	}

	/* Not implemented yet, as assigning to self is not supported in the current implementation
	if t.setupTask.Assignee != "Me" {
		return false, fmt.Errorf("issue not assigned to self")
	}
	*/

	if len(issues) < 2 {
		return false, fmt.Errorf("issue not created")
	}

	var createdIssue *models.Issue
	for i := range issues {
		if issues[i].ID != "pm-abc" {
			createdIssue = issues[i]
			break
		}
	}

	if createdIssue == nil {
		return false, fmt.Errorf("new issue not found")
	}

	if createdIssue.Title == "" {
		return false, fmt.Errorf("issue title is empty")
	}

	if createdIssue.Description == "" {
		return false, fmt.Errorf("issue description is empty")
	}

	/* Not implemented yet, as assigning to self is not supported in the current implementation
	if createdIssue.Assignee != "Me" {
		return false, fmt.Errorf("issue not assigned to self")
	}
	*/

	if createdIssue.Status != models.StatusClosed {
		return false, fmt.Errorf("issue status is not closed")
	}

	if t.setupTask.Status != models.StatusClosed {
		return false, fmt.Errorf("setup issue status is not closed")
	}

	return EndTaskWithTimeout(&t.done, "Task completed!", 5*time.Second)
}
