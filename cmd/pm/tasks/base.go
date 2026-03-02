package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
	"github.com/charmbracelet/huh"
)

type App = app.App
type Config = models.Config
type ValidationFeedback = models.ValidationFeedback

type Issue = models.Issue
type IssueFilter = models.IssueFilter

type Questions = models.Questions
type TaskDetails = models.TaskDetails

type Interface = task.Interface
type InterfaceType = models.InterfaceType

func NewIssueBuilder() *models.IssueBuilder {
	return models.NewIssueBuilder().
		WithStatus(models.StatusOpen).
		WithIssueType(models.TypeTask)
}

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

func BaseDetails() TaskDetails {
	return TaskDetails{
		Title:          "Base Task",
		Description:    "This is a base task.",
		TimeToComplete: "10m",
		Difficulty:     "Easy",
	}
}

func BaseConfig() Config {
	return models.BaseConfig
}

func ClearIssues(app *App) error {
	return app.Issues.DeleteIssues()
}

func BaseQuestions(interfaceType InterfaceType) Questions {
	var taskRating int
	return Questions{
		huh.NewGroup(
			huh.NewConfirm().
				Key("task_completed").
				Title("Did you complete the task?"),
		),
		huh.NewGroup(
			huh.NewSelect[int]().Value(&taskRating).
				Key("task_difficulty").
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

func ReplQuestion(interfaceType InterfaceType, fields ...huh.Field) *huh.Group {
	if interfaceType != InterfaceTypeREPL {
		return nil
	}
	return huh.NewGroup(fields...)
}

func WebQuestion(interfaceType InterfaceType, fields ...huh.Field) *huh.Group {
	if interfaceType != InterfaceTypeWeb {
		return nil
	}
	return huh.NewGroup(fields...)
}

func TUIQuestion(interfaceType InterfaceType, fields ...huh.Field) *huh.Group {
	if interfaceType != InterfaceTypeTUI {
		return nil
	}
	return huh.NewGroup(fields...)
}

// FetchIssues retrieves all issues from the app and returns those that are relevant for validation,
// excluding the setup issue. It also updates the setup issue with the latest data from the app.
func FetchIssues(ctx context.Context, app *app.App, setupIssue *models.Issue) ([]*models.Issue, error) {
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
