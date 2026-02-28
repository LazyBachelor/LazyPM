package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils"
	"github.com/charmbracelet/huh"
)

const reportGenerationDescription = `You are tasked with generating a project status report.

The stakeholders need a weekly status update. Review the current project state and:

1. Identify completed issues since last report
2. Count issues in progress and their status
3. Note any blocked items and blockers
4. Calculate velocity
5. Flag any risks or concerns
6. Update the project status

Add summary comments to at least 3 key issues that stakeholders should know about.`

type ReportGenerationTask struct {
	done       bool
	app        *App
	setupIssue *Issue
}

func NewReportGenerationTask(app *App) *ReportGenerationTask {
	return &ReportGenerationTask{app: app, done: false}
}

func (t *ReportGenerationTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/report-task-stats.json")
}

func (t *ReportGenerationTask) Details() TaskDetails {
	return BaseDetails().
		WithTitle("Status Report Generation").
		WithDescription(reportGenerationDescription).
		WithTimeToComplete("10m").
		WithDifficulty("Easy")
}

func (t *ReportGenerationTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("What is the overall project status?").
				Options(
					huh.NewOption("On Track", 1),
					huh.NewOption("At Risk", 2),
					huh.NewOption("Off Track", 3),
				),
		),
	)
}

func (t *ReportGenerationTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	reportIssues := []*models.Issue{
		NewIssueBuilder().
			WithTitle("User login feature").
			WithDescription("Allow users to login with email/password.").
			WithPriority(1).
			WithStatus(models.StatusClosed).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Password reset").
			WithDescription("Email-based password reset flow. In Progress").
			WithPriority(2).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Database optimization").
			WithDescription("Optimize slow queries identified in profiling.").
			WithPriority(1).
			WithStatus(models.StatusBlocked).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Mobile responsive design").
			WithDescription("Make UI work on mobile devices.").
			WithPriority(2).
			WithStatus(models.StatusClosed).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Third-party API integration").
			WithDescription("Waiting for vendor API documentation.").
			WithPriority(1).
			WithStatus(models.StatusBlocked).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, reportIssues, ""); err != nil {
		return err
	}

	t.setupIssue = NewIssueBuilder().
		WithTitle("Weekly Status Report").
		WithDescription(reportGenerationDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *ReportGenerationTask) Validate(ctx context.Context) ValidationFeedback {
	expect := utils.NewExpector()

	return expect.Complete()
}
