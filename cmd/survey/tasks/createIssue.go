package tasks

import (
	"context"
	"errors"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	ui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	"github.com/charmbracelet/huh"
)

func NewCreateIssueTask() *task.Task {
	aboutScreen := ui.NewTaskModel(createIssueDetails())
	questionnaire := ui.NewQuestionnaireModel(createIssueQuestionnaire())

	task := task.NewTask(aboutScreen, questionnaire)
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
		WebAddress:            ":8080",
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
	if err := svc.DeleteIssues(); err != nil {
		return err
	}

	issues := []*models.Issue{
		{Title: "Test Issue", Description: "Long Description", IssueType: models.TypeBug, Status: models.StatusBlocked},
	}

	if err := svc.Beads.CreateIssues(ctx, issues, "actor"); err != nil {
		return err
	}

	return nil
}

func createIssueValidate(ctx context.Context, svc *service.Services) (ok bool, errorMsg error) {
	issues, err := svc.Beads.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return false, err
	}

	if len(issues) == 0 {
		return false, errors.New("no issues found. Please create an issue to proceed.")
	}

	return true, nil
}
