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

func NewCreateIssueTask(svc *service.Services) *task.Task {
	aboutScreen := ui.NewTaskModel(createIssueDetails())
	questionnaire := ui.NewQuestionnaireModel(createIssueQuestionnaire())

	task := task.NewTask(svc, aboutScreen, questionnaire)
	task.SetConfigFunc(createIssueConfig)
	task.SetDbStateFunc(createIssueDbState)
	task.SetValidateFunc(createIssueValidate)
	return task
}

func createIssueConfig() task.TaskConfig {
	return task.TaskConfig{
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/task-1-stats.json",
		WebAddress:            "localhost:8080",
	}
}

func createIssueDetails() ui.TaskDetails {
	return ui.TaskDetails{
		Title:          "Create Issue Task",
		Description:    "Create a new issue in the project management system to test the issue creation workflow.",
		TimeToComplete: "15m",
		Difficulty:     "Hard",
	}
}

func createIssueQuestionnaire() ui.Questions {
	return ui.Questions{
		huh.NewGroup(
			huh.NewConfirm().Title("Was this good"),
		),
		huh.NewGroup(
			huh.NewSelect[int]().Options(
				huh.NewOption("Very good", 1),
				huh.NewOption("Very Bad", 2),
			).Title("How good was it?"),
		),
	}
}

func createIssueDbState(ctx context.Context, svc *service.Services) error {
	// Clear existing issues to ensure a clean state for the task
	if err := svc.DeleteIssues(); err != nil {
		return err
	}

	issue := models.Issue{
		ID:    "pm-abc",
		Title: "Create A New Issue", Description: description,
		IssueType: models.TypeTask, Status: models.StatusOpen,
	}

	if err := svc.Beads.CreateIssue(ctx, &issue, ""); err != nil {
		return err
	}

	return nil
}

func createIssueValidate(ctx context.Context, svc *service.Services) (ok bool, errorMsg error) {
	// Fetches issues, indexed with latest first
	issues, err := svc.Beads.SearchIssues(ctx, "", models.IssueFilter{})
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
