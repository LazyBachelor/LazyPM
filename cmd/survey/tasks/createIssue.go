package tasks

import (
	"context"
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	ui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	"github.com/charmbracelet/huh"
)

const description = `You are tasked with creating a new issue in the project management system.
This task will test your ability to use the issue creation workflow effectively.

Assign this task to yourself and start creating the issue.
Make sure to fill out all the necessary details, including the title, description, and assignee.`

type CreateIssueTask struct {
	*task.Task
	svc *service.Services
}

func NewCreateIssueTask(svc *service.Services) *CreateIssueTask {
	return &CreateIssueTask{
		svc: svc,
	}
}

func (t *CreateIssueTask) Init() *task.Task {
	task := task.NewTask(t.svc, t.Details(), t.QuestionsFunc())
	task.SetConfigFunc(t.Config)
	task.SetDbStateFunc(t.DbStateFunc)
	task.SetValidateFunc(t.ValidateFunc)
	return task
}

func (t *CreateIssueTask) Config() task.TaskConfig {
	return task.TaskConfig{
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/task-1-stats.json",
		WebAddress:            "localhost:8080",
	}
}

func (t *CreateIssueTask) Details() ui.TaskDetails {
	return ui.TaskDetails{
		Title:          "Create Issue Task",
		Description:    "Create a new issue in the project management system to test the issue creation workflow.",
		TimeToComplete: "15m",
		Difficulty:     "Hard",
	}
}

func (t *CreateIssueTask) QuestionsFunc() task.QuestionsFunc {
	return func(interfaceType task.InterfaceType) ui.Questions {
		questions := ui.Questions{}

		// Add an extra question for TUI
		if interfaceType == task.InterfaceTUI {
			questions = append(questions,
				huh.NewGroup(
					huh.NewConfirm().Title("Did you complete?"),
				),
			)
		}

		questions = append(questions,
			huh.NewGroup(
				huh.NewConfirm().Title("Was this good"),
			),
			huh.NewGroup(
				huh.NewSelect[int]().Options(
					huh.NewOption("Very good", 1),
					huh.NewOption("Very Bad", 2),
				).Title("How good was it?"),
			),
		)

		return questions
	}
}

func (t *CreateIssueTask) DbStateFunc(ctx context.Context) error {
	// Clear existing issues to ensure a clean state for the task
	if err := t.svc.DeleteIssues(); err != nil {
		return err
	}

	issue := models.Issue{
		ID:    "pm-abc",
		Title: "Create A New Issue", Description: description,
		IssueType: models.TypeTask, Status: models.StatusOpen,
	}

	if err := t.svc.Beads.CreateIssue(ctx, &issue, ""); err != nil {
		return err
	}

	return nil
}

func (t *CreateIssueTask) ValidateFunc(ctx context.Context) (ok bool, errorMsg error) {
	// Fetches issues, indexed with latest first
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
