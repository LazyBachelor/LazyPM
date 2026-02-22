package tasks

import (
	"context"
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
)

const description = `You are tasked with creating a new issue in the project management system.
This task will test your ability to use the issue creation workflow effectively.

Assign this task to yourself and start creating the issue.
Make sure to fill out all the necessary details, including the title, description, and assignee.`

type CreateIssueTask struct {
	app *service.App
}

func NewCreateIssueTask(app *service.App) *CreateIssueTask {
	return &CreateIssueTask{app: app}
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
	// Clear existing issues to ensure a clean state for the task
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	issue := models.NewBaseIssue().
		WithTitle("Create a New Issue").WithDescription(description).Build()

	return t.app.Issues.CreateIssue(ctx, &issue, "")
}

func (t *CreateIssueTask) Validate(ctx context.Context) (bool, error) {
	issues, err := t.app.Issues.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return false, err
	}

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

	return true, nil
}
