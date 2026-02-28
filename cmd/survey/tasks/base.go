package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	taskui "github.com/LazyBachelor/LazyPM/pkg/task/ui"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
	"github.com/charmbracelet/huh"
)

type ValidationFeedback = models.ValidationFeedback
type InterfaceType = models.InterfaceType
type Interface = task.Interface

const (
	InterfaceTypeCLI  = models.InterfaceTypeCLI
	InterfaceTypeTUI  = models.InterfaceTypeTUI
	InterfaceTypeWeb  = models.InterfaceTypeWeb
	InterfaceTypeREPL = models.InterfaceTypeREPL
)

func InterfaceToType(it Interface) InterfaceType {
	switch it.(type) {
	case *repl.REPL:
		return InterfaceTypeREPL
	case *tui.Tui:
		return InterfaceTypeTUI
	case *web.Web:
		return InterfaceTypeWeb
	default:
		return InterfaceType("unknown")
	}
}

func BaseDetails() taskui.TaskDetails {
	return taskui.TaskDetails{
		Title:          "Base Task",
		Description:    "This is a base task.",
		TimeToComplete: "10m",
		Difficulty:     "Easy",
	}
}

func BaseConfig() task.Config {
	return models.BaseConfig
}

func ClearIssues(app *service.App) error {
	return app.Issues.DeleteIssues()
}

func BaseQuestions(interfaceType task.InterfaceType) taskui.Questions {
	var taskRating int
	return taskui.Questions{
		huh.NewGroup(
			huh.NewConfirm().
				Title("Did you complete the task?"),
		),
		huh.NewGroup(
			huh.NewSelect[int]().Value(&taskRating).
				Options(
					huh.NewOption("Very easy", 1),
					huh.NewOption("Easy", 2),
					huh.NewOption("Moderate", 3),
					huh.NewOption("Hard", 4),
				).
				Title("How difficult was the task?"),
		),
	}
}

func Question(fields ...huh.Field) *huh.Group {
	return huh.NewGroup(fields...)
}

func ReplQuestion(interfaceType task.InterfaceType, fields ...huh.Field) *huh.Group {
	if interfaceType != InterfaceTypeREPL {
		return nil
	}
	return huh.NewGroup(fields...)
}

func WebQuestion(interfaceType task.InterfaceType, fields ...huh.Field) *huh.Group {
	if interfaceType != InterfaceTypeWeb {
		return nil
	}
	return huh.NewGroup(fields...)
}

func TUIQuestion(interfaceType task.InterfaceType, fields ...huh.Field) *huh.Group {
	if interfaceType != InterfaceTypeTUI {
		return nil
	}
	return huh.NewGroup(fields...)
}

// FetchIssues retrives all issues from the app and returns those that are relevant for validation,
// excluding the setup issue. It also updates the setup issue with the latest data from the app.
func FetchIssues(ctx context.Context, app *service.App, setupIssue *models.Issue) ([]*models.Issue, error) {
	issues, err := app.Issues.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return nil, err
	}

	var relevantIssues []*models.Issue
	for _, issue := range issues {
		if issue.ID != setupIssue.ID {
			relevantIssues = append(relevantIssues, issue)
		} else {
			*setupIssue = *issue
		}
	}

	return relevantIssues, nil
}
