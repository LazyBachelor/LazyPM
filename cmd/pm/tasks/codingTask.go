package tasks

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const codingDescription = `You are tasked with doing a chore in the codebase.

This task will test your ability to read and understand instructions, change text, and save it to a file.

The MongoDB Driver dependency in the file is outdated and needs to be updated to the latest version.
This is a common task for developers, and it requires attention to detail and the ability to follow instructions carefully.

Your task:
1. Assign this Issue to yourself as "Me" and mark it as "In Progress"
2. Create a New Issue and give it these details:
	- Title: "Upgrade MongoDB Driver Dependency"
	- Description: "We need to upgrade the MongoDB Driver dependency to the latest version."
	- Status: "In Progress"
	- Issue Type: "Chore"
3. A file will appear on your desktop named "code.txt".
   Open it and follow the instructions inside. And save the file after you are done.
4. When you are done, mark this and the issue you made as "Closed".`

var textFileDescription = `
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
	done       bool
	setupIssue *Issue
	app        *App
}

func NewCodingTask(app *App) *CodingTask {
	return &CodingTask{app: app, done: false}
}

func (t *CodingTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/coding-task-stats.json")
}

func (t *CodingTask) Details() TaskDetails {
	return BaseDetails().WithTitle("Coding Task").WithDescription(codingDescription)
}

func (t *CodingTask) Questions(interfaceType InterfaceType) (questions Questions) {
	return BaseQuestions(interfaceType)
}

func (t *CodingTask) QuestionnaireKeys(interfaceType InterfaceType) []string {
	keys := []string{"task_completed", "task_difficulty", "final_confirmation"}
	return keys
}

func (t *CodingTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	t.setupIssue = NewIssueBuilder().
		WithTitle("Coding Task - Upgrade MongoDB Driver").
		WithDescription(codingDescription).
		WithStatus(models.StatusOpen).
		WithIssueType(models.TypeTask).
		Build()

	if err := t.app.Issues.CreateIssue(ctx, t.setupIssue, "LazyPM"); err != nil {
		return err
	}

	if err := os.WriteFile("./code.txt", []byte(textFileContent), 0644); err != nil {
		return err
	}

	return nil
}

func (t *CodingTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app, t.setupIssue)
	if err != nil {
		return expect.ValidationFeedback
	}

	expect.Assert(t.setupIssue.Assignee == "Me", "The original issue should be assigned to 'Me'. Please assign the issue to yourself.")
	expect.Assert(t.setupIssue.Status == models.StatusInProgress, "The original issue should be marked as In Progress before starting work. Please update the issue status to In Progress.")
	expect.Assert(len(issues) > 0, "No new issues created. Please create an issue with the specified details.")

	if len(issues) == 0 {
		expect.Fail("No new issues created")
		return expect.ValidationFeedback
	}

	issue := issues[0]

	expect.Assert(len(issues) < 2, "Multiple issues were created instead of one. Delete the extra issues and try again.")

	expect.NotEmptyString(issue.Title, "Issue title should not be empty")
	expect.Assert(issue.Title == "Upgrade MongoDB Driver Dependency",
		fmt.Sprintf("Issue title does not match the expected value 'Upgrade MongoDB Driver Dependency', but was '%s'", issue.Title))

	expect.NotEmptyString(issue.Description, "Issue description should not be empty")
	expect.Assert(issue.Description == "We need to upgrade the MongoDB Driver dependency to the latest version.",
		fmt.Sprintf("Issue description does not match the expected value 'We need to upgrade the MongoDB Driver dependency to the latest version.', but was '%s'", issue.Description))

	expect.Assert(issue.IssueType == models.TypeChore,
		fmt.Sprintf("Issue type should be 'Chore', but was '%s'", issue.IssueType))

	expect.Assert(issue.Assignee == "Me",
		fmt.Sprintf("Issue should be assigned to 'Me', but was assigned to '%s'", issue.Assignee))

	if issue.Status == models.StatusInProgress || isInProgress {
		isInProgress = true
	} else {
		expect.Fail("Issue should be marked as in-progress when work starts")
	}

	if _, err := os.Stat("./code.txt"); os.IsNotExist(err) {
		expect.Fail("The code.txt file should exist on the desktop. Please create the file with the specified content.")
		return expect.ValidationFeedback
	}

	fileContent, err := os.ReadFile("./code.txt")
	if err != nil {
		expect.Fail("Error reading code.txt file: " + err.Error())
		return expect.ValidationFeedback
	}

	code, ok := strings.CutPrefix(string(fileContent), codingDescription+textFileDescription+"\n")
	if !ok {
		expect.Fail("The content of code.txt does not match the expected format. Please ensure the file contains the correct instructions and code.")
		return expect.ValidationFeedback
	}

	expect.Assert(strings.Contains(code, "go.mongodb.org/mongo-driver v1.17.9"),
		"The go.mod snippet in code.txt should contain the updated MongoDB Driver version (v1.17.9). Please make sure you have updated the dependency in the snippet accordingly.")

	if !isInProgress {
		return expect.ValidationFeedback
	} else if issue.Status != models.StatusClosed {
		expect.Fail("Issue should be set to Closed once the work is completed")
	} else {
		expect.Assert(t.setupIssue.Status == models.StatusClosed,
			"The original setup issue should be set to Closed once the work is completed")
	}

	return expect.Complete()
}
