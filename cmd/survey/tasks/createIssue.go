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
	svc *service.Services
}

func NewCreateIssueTask(svc *service.Services) *CreateIssueTask {
	return &CreateIssueTask{svc: svc}
}

func (t *CreateIssueTask) Config() task.TaskConfig {
	config := BaseConfig()
	config.StatisticsStoragePath = "./.pm/create-issue-stats.json"
	return config
}

func (t *CreateIssueTask) Details() taskui.TaskDetails {
	details := BaseDetails()
	details.Title = "Create Issue Task"
	details.Description = description
	return details
}

func (t *CreateIssueTask) Questions(interfaceType task.InterfaceType) taskui.Questions {
	questions := BaseQuestions(interfaceType)
	return questions
}

func (t *CreateIssueTask) Setup(ctx context.Context) error {
	// Clear existing issues to ensure a clean state for the task
	if err := t.svc.DeleteIssues(); err != nil {
		return err
	}

	issue := models.Issue{
		ID:          "pm-abc",
		Title:       "Create A New Issue",
		Description: description,
		IssueType:   models.TypeTask,
		Status:      models.StatusOpen,
	}

	return t.svc.Beads.CreateIssue(ctx, &issue, "")
}

func (t *CreateIssueTask) Validate(ctx context.Context) (bool, error) {
	issues, err := t.svc.Beads.SearchIssues(ctx, "", models.IssueFilter{})
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
