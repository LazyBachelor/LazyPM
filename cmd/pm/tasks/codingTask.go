package tasks

import (
	"context"
	"os"
	"strings"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const codingDescription = `You are tasked with fixing a logical error in the code.

Your task:
1. A text file will appear in the this directory:
   - Open it and follow the instructions inside.
   - Save the file after you are done.
2. When you are done, mark this task as "Closed".`

var textFileDescription = `

# Instructions for the coding task

There is a major logical error in this code, you need to fix it.
Change the function logic so that it correctly adds two numbers together instead of subtracting them.
############################################################`

var code = `
function Add(a, b int) int {
	return a - b
}
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
		WithTitle("Fix major error").
		WithDescription(codingDescription).
		WithStatus(models.StatusInProgress).
		WithPriority(4).
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

	expect.Contains(code, "a + b", "Function logic")
	if !expect.Valid() {
		return expect.ValidationFeedback
	}

	expect.Equal(issue.Status, models.StatusClosed, "Issue Status")

	return expect.Complete()
}
