package tasks

import (
	"context"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
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
		interfaceDesc = `What is a REPL Interface?
- REPL stands for Read-Eval-Print Loop. It is an interactive programming environment that takes single user inputs (reads), executes them (eval), and returns the result to the user (print), then waits for the next input (loop).
- In this task, you will interact with the task through a REPL interface, which allows you to execute commands and receive immediate feedback in a command-line environment.
- You can type commands to perform actions related to the task, and the REPL will process those commands and provide responses based on your input.
- The REPL interface is designed to facilitate a more dynamic and interactive way of completing the task, allowing you to experiment and receive real-time feedback as you work through the task requirements.`
	case InterfaceTypeTUI:
		interfaceDesc = `What is a TUI Interface?
- TUI stands for Text User Interface. It is a user interface that uses text-based elements to allow users to interact with the application.
- The main way to interact with a TUI is through keyboard inputs, where you can navigate through menus, select options, and input data using the keyboard.
- In this task, you will interact with the task through a TUI interface, which provides a more structured and visually organized way to complete the task using text-based menus, forms, and other interactive elements.
- The TUI interface is designed to enhance usability and provide a more engaging experience while working through the task requirements, allowing you to navigate through options and input information in a more intuitive way.`
	case InterfaceTypeWeb:
		interfaceDesc = `What is a Web Interface?
- A Web Interface is a user interface that is accessed through a web browser. It allows users to interact with the application using graphical elements such as buttons, forms, and menus.
- In this task, you will interact with the task through a Web interface, which provides a more visually rich and user-friendly way to complete the task using a web-based platform.
- The Web interface is designed to enhance usability and provide a more engaging experience while working through the task requirements, allowing you to navigate through options and input information in a more intuitive way using a graphical interface.`
	default:
		interfaceDesc = "Unknown Interface"
	}

	return TaskDetails{
		Title:                "title not set",
		Description:          "description not set",
		TimeToComplete:       "not set",
		Difficulty:           "not set",
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
				Title("Were you able to complete the task?")),
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
