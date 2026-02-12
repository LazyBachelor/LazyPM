package tasks

import (
	"context"
	"errors"

	"github.com/LazyBachelor/LazyPM/cmd/survey/ui"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/charmbracelet/huh"
)

func NewCreateIssueTask(interfaceType Interface) *Task {
	aboutScreen := ui.NewTaskModel(createIssueDetails())
	questionare := ui.NewQuestionnaireModel(createIssueQuestionare())

	task := NewTask(interfaceType, aboutScreen, questionare)
	task.SetDbStateFunc(createIssueDbState)
	task.SetValidateFunc(createIssueValidate)

	return task
}

func createIssueDetails() ui.TaskDetails {
	return ui.TaskDetails{
		Title:          "Create Issue Task",
		Description:    `Description for create issue task`,
		TimeToComplete: "15m",
		Difficulty:     "Hard",
	}
}

func createIssueQuestionare() ui.Questions {
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
		return false, errors.New("No issues found. Please create an issue to proceed.")
	}

	return true, nil
}
