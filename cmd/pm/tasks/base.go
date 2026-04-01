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
		interfaceDesc = `About our REPL Interface?

- REPL stands for Read-Eval-Print Loop.
  It's an interactive programming environment that takes single user inputs (reads),
  executes them (eval), and returns the result to the user (print),
  then waits for the next input (loop).

- In a traditional CLI, you would use the help command to see available commands.

- You can type 'pm help' to get a list of available commands and their descriptions.

- You can run shell commands directly from the REPL, like "git" or "ls".

- When you type commands in this interface, suggestions will appear as you type,
  which you can select to auto-complete your command.

- Write "exit" to skip task during the survey.`
	case InterfaceTypeTUI:
		interfaceDesc = `About our TUI Interface?

- TUI stands for Terminal User Interface.
  It is a user interface that uses text-based elements to interact with the application.

- The main way to interact with a TUI is through keyboard inputs,
  where you can navigate through menus, select options, and input data using the keyboard.

- At the bottom of the interface you will find the help menu, which lists available keybinds.

- The interface has parts:
  - List View: This view shows the available issues in a list format,
    allowing you to browse through them and select one to work on.
  - Kanban: This view allows you to manage and track the progress of different issues,
    this also they main way of creating and managing sprints.

- Press "q" to quit task during the survey.`
	case InterfaceTypeWeb:
		interfaceDesc = `About our Web Interface?

- A Web Interface is a user interface that is accessed through a web browser.
  Users to interact with the application using graphical elements such as buttons and forms.

- We designed this interface to only be interactive through mouse clicks.

- The interface has parts:
  - Issues: This view shows the available issues in a list format.
  - Kanban: This view allows you to manage and track the progress of different issues.

- Press q in the server terminal to skip the task.`
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
	return []string{"interface_difficulty", "task_difficulty"}
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

// FetchIssues retrieves all issues from the app and returns those that are relevant for validation
func FetchIssues(ctx context.Context, app *App) ([]*Issue, error) {
	issues, err := app.Issues.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return nil, err
	}

	var relevantIssues []*Issue
	for _, issue := range issues {
		relevantIssues = append(relevantIssues, issue)
	}

	return relevantIssues, nil
}
