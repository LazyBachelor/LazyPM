package tasks

import (
	"context"
	"os"
	"strings"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const codingDescription = `You are tasked with doing a upgrading a dependency in the codebase.

Your task:
1. Assign the given issue to yourself as 'Me'.
2. A text file will appear in the this directory:
   - Open it and follow the instructions inside.
   - Save the file after you are done.
3. When you are done, mark this task as "Closed".`

var textFileDescription = `

# Instructions for the coding task

Please upgrade the MongoDB Driver dependency in the go.mod file to the latest version.
It should be v1.17.9. After you are done, save the file and mark the task as completed.
############################################################`

var code = `
require (
	charm.land/lipgloss/v2 v2.0.0-beta.3.0.20251106193318-19329a3e8410
	github.com/go-git/go-git/v6 v6.0.0-20260222090600-424e9964d3a3
	github.com/muesli/reflow v0.3.0
	github.com/steveyegge/beads v0.49.6
	go.mongodb.org/mongo-driver v1.17.8
	go.mongodb.org/mongo-driver/v2 v2.5.0
)

require (
	github.com/c-bata/go-prompt v0.2.6
	github.com/charmbracelet/bubbles v0.21.1
	github.com/charmbracelet/bubbletea v1.3.10
	github.com/charmbracelet/fang v0.4.4
	github.com/charmbracelet/huh v0.8.0
	github.com/charmbracelet/lipgloss v1.1.1-0.20250404203927-76690c660834
	github.com/spf13/cobra v1.10.2
	golang.org/x/term v0.40.0
)

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/a-h/templ v0.3.977
	github.com/donseba/go-htmx v1.12.1
	github.com/go-chi/chi/v5 v5.2.5
	github.com/go-playground/form/v4 v4.3.0
	github.com/go-playground/validator/v10 v10.30.1
	github.com/rs/cors v1.11.1
)

tool (
	github.com/a-h/templ/cmd/templ
	github.com/haatos/goshipit/cmd/gsi
)
`

var textFileContent = codingDescription + textFileDescription + "\n" + code

type CodingTask struct {
	done  bool
	app   *App
	issue *Issue
}

func NewCodingTask(app *App) *CodingTask {
	return &CodingTask{app: app, done: false}
}

func (t *CodingTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/coding-task-stats.json")
}

func (t *CodingTask) Details(interfaceType InterfaceType) TaskDetails {
	return BaseDetails(interfaceType).
		WithTitle("Coding Task").
		WithDescription(codingDescription).
		WithDifficulty("Hard").WithTimeToComplete("5m")
}

func (t *CodingTask) Questions(interfaceType InterfaceType) (questions Questions) {
	return BaseQuestions(interfaceType).
		With(
			huh.NewGroup(
				huh.NewSelect[int]().
					Key("coding_interface_friction").
					Title("How much friction did you feel switching between editing code and using the interface?").
					Description("By friction we mean the difficulty or inconvenience of switching between these two activities.").
					Options(
						huh.NewOption("Very low", 1),
						huh.NewOption("Low", 2),
						huh.NewOption("Moderate", 3),
						huh.NewOption("High", 4),
						huh.NewOption("Very high", 5),
					),
			),
		)
}

func (t *CodingTask) QuestionnaireKeys(interfaceType InterfaceType) []string {
	return BaseKeys().With("coding_interface_friction")
}

func (t *CodingTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	if err := os.WriteFile("./code.txt", []byte(textFileContent), 0644); err != nil {
		return err
	}

	t.issue = NewIssueBuilder().
		WithTitle("Upgrade MongoDB Driver").
		WithPriority(4).WithDescription(codingDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.issue, "")
}

func (t *CodingTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issue, err := t.app.Issues.GetIssue(ctx, t.issue.ID)
	if err != nil {
		return expect.Fatal("Could not fetch issues")
	}
	if issue == nil {
		return expect.Fatal("Issue was deleted or could not be found")
	}

	if !expect.Equal(issue.Assignee, "Me", "Issue Assignee").Valid() {
		return expect.ValidationFeedback
	}

	if _, err := os.Stat("./code.txt"); os.IsNotExist(err) {
		expect.Fail("The code.txt file should exist on the desktop.")
		return expect.ValidationFeedback
	}

	fileContent, err := os.ReadFile("./code.txt")
	if err != nil {
		expect.Fail("Error reading code.txt file: " + err.Error())
		return expect.ValidationFeedback
	}

	code, ok := strings.CutPrefix(string(fileContent), codingDescription+textFileDescription+"\n")
	if !ok {
		expect.Fail("The content of code.txt does not match the expected format.")
		return expect.ValidationFeedback
	}

	expect.Contains(code, "go.mongodb.org/mongo-driver v1.17.9", "MongoDB Driver version")
	if !expect.Valid() {
		return expect.ValidationFeedback
	}

	expect.Equal(issue.Status, models.StatusClosed, "Issue Status")

	return expect.Complete()
}
