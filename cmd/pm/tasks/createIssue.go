package tasks

import (
	"context"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const description = `You are tasked with creating a new issue in the project management system.

Your task:
1. Create a new issue with the title "My first Issue"
2. Add this detailed description "I need to do some coding"
3. Assign the issue to yourself as "Me"
4. Mark the issue as In Progress when you are done.

Make sure to fill out all the necessary details to help others understand the work item.`

type CreateIssueTask struct {
	done       bool
	app        *App
	setupIssue *Issue
}

func NewCreateIssueTask(app *App) *CreateIssueTask {
	return &CreateIssueTask{app: app, done: false}
}

func (t *CreateIssueTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/create-issue-stats.json")
}

func (t *CreateIssueTask) Details(interfaceType InterfaceType) TaskDetails {
	return BaseDetails(interfaceType).
		WithTitle("Create Issue Task").
		WithDescription(description).
		WithDifficulty("Medium").
		WithTimeToComplete("3m")
}

func (t *CreateIssueTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().Key("context_difficulty").
				Title("How difficult was it remembering what needed to be done for the task?").
				Description("We are interested in how difficult it was keeping track of the task requirements. And if one interface helps you keep context").
				Options(
					huh.NewOption("Very Easy", 1),
					huh.NewOption("Easy", 2),
					huh.NewOption("Moderate", 3),
					huh.NewOption("Difficult", 4),
					huh.NewOption("Very Difficult", 5),
				),
		),
	)
}

func (t *CreateIssueTask) QuestionnaireKeys(_ InterfaceType) []string {
	return BaseKeys().With("context_difficulty")
}

func (t *CreateIssueTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	return nil
}

func (t *CreateIssueTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app)
	if err != nil {
		return expect.Fatal("Could not fetch issues")
	}

	if len(issues) == 0 {
		expect.Fail("No issues were created")
		return expect.ValidationFeedback
	}

	issue := issues[0]

	expect.NotEmptyAndEqual(issue.Title, "My first Issue", "Issue title")
	expect.NotEmptyAndEqual(issue.Description, "I need to do some coding", "Issue description")
	expect.NotEmptyAndEqual(issue.Assignee, "Me", "Issue assignee")
	expect.Equal(issue.Status, models.StatusInProgress, "Issue status")

	return expect.Complete()
}
