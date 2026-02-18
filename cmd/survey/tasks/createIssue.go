package tasks

import (
	"context"
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	"github.com/charmbracelet/huh"
)

const description = `You are tasked with creating a new issue in the project management system.
This task will test your ability to use the issue creation workflow effectively.

Assign this task to yourself and start creating the issue.
Make sure to fill out all the necessary details, including the title, description, and assignee.`

func init() {
	task.Register("create_issue", func(svc *service.Services) task.Tasker {
		return NewCreateIssueTask(svc)
	})
}

type CreateIssueTask struct {
	svc *service.Services
}

func NewCreateIssueTask(svc *service.Services) *CreateIssueTask {
	return &CreateIssueTask{svc: svc}
}

func (t *CreateIssueTask) Config() task.TaskConfig {
	return task.TaskConfig{
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/task-1-stats.json",
		WebAddress:            "localhost:8080",
	}
}

func (t *CreateIssueTask) Details() taskui.TaskDetails {
	return taskui.TaskDetails{
		Title:          "Create Issue Task",
		Description:    "Create a new issue in the project management system...",
		TimeToComplete: "15m",
		Difficulty:     "Hard",
	}
}

func (t *CreateIssueTask) Questions(interfaceType task.InterfaceType) taskui.Questions {
	return taskui.Questions{
		huh.NewGroup(huh.NewConfirm().Title("Was this good")),
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(
					huh.NewOption("Very good", 1),
					huh.NewOption("Very Bad", 2),
				).
				Title("How good was it?"),
		),
	}
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

	if issues[0].Title == "" {
		return false, fmt.Errorf("issue title is empty")
	}

	if issues[0].Description == "" {
		return false, fmt.Errorf("issue description is empty")
	}

	return true, nil
}
