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

func BaseDetails(interfaceType InterfaceType) TaskDetails {

	var interfaceDesc string

	switch interfaceType {
	case InterfaceTypeREPL:
		interfaceDesc = `How to use the REPL Interface
- The REPL interface allows you to interact with the task using a command-line interface.
- You can type commands to perform actions related to the task, such as creating issues, updating statuses, etc.
- The interface will provide prompts and feedback based on your inputs.`
	case InterfaceTypeTUI:
		interfaceDesc = `How to use the TUI Interface
- The TUI (Text User Interface) provides a more interactive experience in the terminal.
- You can navigate through menus, select options, and view task details in a structured format.
- Use keyboard shortcuts to perform actions and explore different sections of the interface.`
	case InterfaceTypeWeb:
		interfaceDesc = `How to use the Web Interface
- The Web interface allows you to interact with the task through a web browser.
- You can access the interface by navigating to the provided URL.
- The interface will have buttons, forms, and other interactive elements to help you complete the task.`
	default:
		interfaceDesc = "Unknown Interface"
	}

	return TaskDetails{
		Title:                "Base Task",
		Description:          "This is a base task.",
		TimeToComplete:       "10m",
		Difficulty:           "Easy",
		InterfaceType:        interfaceType,
		InterfaceDescription: interfaceDesc,
	}
}

func BaseConfig() Config {
	return models.BaseConfig
}

func ClearIssues(app *App) error {
	return app.Issues.DeleteIssues()
}

func BaseQuestions(interfaceType InterfaceType) Questions {
	return Questions{
		huh.NewGroup(
			huh.NewConfirm().Key("task_completed").
				Title("Where you able to complete the task?").
				Description("")),
		huh.NewGroup(
			huh.NewSelect[int]().Key("interface_difficulty").
				Title("How difficult was it to use the interface?").
				Description("By interface we mean the method of interaction with the task.").
				Options(
					huh.NewOption("Very easy", 1),
					huh.NewOption("Easy", 2),
					huh.NewOption("Moderate", 3),
					huh.NewOption("Hard", 4),
				),
		),
		huh.NewGroup(
			huh.NewSelect[int]().Key("task_difficulty").
				Title("How difficult did you find the task?").
				Description("Only rate the difficulty of the task itself. Not the usability of the interface").
				Options(
					huh.NewOption("Very easy", 1),
					huh.NewOption("Easy", 2),
					huh.NewOption("Moderate", 3),
					huh.NewOption("Hard", 4),
				),
		),
	}
}

type Keys []string

func BaseKeys() Keys {
	return []string{"task_completed", "interface_difficulty", "task_difficulty"}
}

func (k Keys) With(keys ...string) Keys {
	return append(k, keys...)
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
func FetchIssues(ctx context.Context, app *App, setupIssue *Issue) ([]*Issue, error) {
	issues, err := app.Issues.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return nil, err
	}

	var relevantIssues []*Issue
	for _, issue := range issues {
		if issue.ID != setupIssue.ID {
			relevantIssues = append(relevantIssues, issue)
		} else {
			*setupIssue = *issue
		}
	}

	return relevantIssues, nil
}
