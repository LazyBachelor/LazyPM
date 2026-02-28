package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils"
	"github.com/charmbracelet/huh"
)

const stakeholderUpdateDescription = `You are tasked with preparing stakeholder updates.

A key stakeholder has requested an update on the progress of their requested features.
You need to:

1. Identify issues related to the stakeholder's requests (marked with "stakeholder" in description)
2. Review the current status of each issue
3. Provide clear, non-technical status updates via comments
4. Highlight any blockers or delays
5. Set realistic expectations for delivery
6. Close completed items with completion notes

The stakeholder is interested in the Dashboard Enhancement and Export Features specifically.`

type StakeholderUpdateTask struct {
	done       bool
	app        *App
	setupIssue *Issue
}

func NewStakeholderUpdateTask(app *App) *StakeholderUpdateTask {
	return &StakeholderUpdateTask{app: app, done: false}
}

func (t *StakeholderUpdateTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/stakeholder-task-stats.json")
}

func (t *StakeholderUpdateTask) Details() TaskDetails {
	return BaseDetails().
		WithTitle("Stakeholder Update Task").
		WithDescription(stakeholderUpdateDescription).
		WithTimeToComplete("10m").
		WithDifficulty("Easy")
}

func (t *StakeholderUpdateTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("How satisfied will the stakeholder be with this update?").
				Options(
					huh.NewOption("Very satisfied", 1),
					huh.NewOption("Satisfied", 2),
					huh.NewOption("Neutral", 3),
					huh.NewOption("Concerned", 4),
				),
		),
	)
}

func (t *StakeholderUpdateTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	stakeholderIssues := []*models.Issue{
		NewIssueBuilder().
			WithTitle("Dashboard Enhancement - Charts").
			WithDescription("Add interactive charts to the main dashboard (stakeholder request). COMPLETED").
			WithPriority(1).
			WithStatus(models.StatusClosed).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Dashboard Enhancement - Filters").
			WithDescription("Add date range filters to dashboard (stakeholder request). In Progress").
			WithPriority(2).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Export to PDF").
			WithDescription("Allow exporting reports to PDF format (stakeholder request). In Progress").
			WithPriority(1).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Export to Excel").
			WithDescription("Allow exporting data to Excel format (stakeholder request). BLOCKED - waiting for library approval").
			WithPriority(2).
			WithStatus(models.StatusBlocked).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, stakeholderIssues, ""); err != nil {
		return err
	}

	t.setupIssue = NewIssueBuilder().
		WithTitle("Stakeholder Update Preparation").
		WithDescription(stakeholderUpdateDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *StakeholderUpdateTask) Validate(ctx context.Context) ValidationFeedback {
	expect := utils.NewExpector()

	return expect.Complete()
}
