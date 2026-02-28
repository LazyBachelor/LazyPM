package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils"
	"github.com/charmbracelet/huh"
)

const issueTriageDescription = `You are tasked with triaging incoming issues.

The support team has submitted several bug reports and feature requests that need to be reviewed and prioritized. Your job is to:

1. Review each incoming issue
2. Assign appropriate priority
3. Set the correct status
4. Identify if it's a bug, feature, or task
5. Leave comments explaining your decisions

Make decisions quickly but thoughtfully. Not everything is high priority!`

type IssueTriageTask struct {
	done       bool
	app        *App
	setupIssue *models.Issue
}

func NewIssueTriageTask(app *App) *IssueTriageTask {
	return &IssueTriageTask{app: app, done: false}
}

func (t *IssueTriageTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/triage-task-stats.json")
}

func (t *IssueTriageTask) Details() TaskDetails {
	return BaseDetails().
		WithTitle("Issue Triage Task").
		WithDescription(issueTriageDescription).
		WithTimeToComplete("12m").
		WithDifficulty("Medium")
}

func (t *IssueTriageTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("How many Critical priority issues did you identify?").
				Options(
					huh.NewOption("0", 0),
					huh.NewOption("1", 1),
					huh.NewOption("2", 2),
					huh.NewOption("3+", 3),
				),
		),
	)
}

func (t *IssueTriageTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	triageIssues := []*models.Issue{
		models.NewIssueBuilder().
			WithTitle("App crashes on login").
			WithDescription("Users report the app crashes immediately after entering credentials. Affects all users on Android 12. Needs urgent attention.").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Dark mode support").
			WithDescription("Users requesting dark mode theme for better night time usage. Nice to have feature.").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Database timeout errors").
			WithDescription("Intermittent timeouts when querying large datasets. Needs investigation.").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Add export to CSV feature").
			WithDescription("Sales team needs ability to export reports to CSV format.").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		models.NewIssueBuilder().
			WithTitle("Security: Password reset vulnerability").
			WithDescription("Reported by security audit - password reset tokens don't expire. Critical security issue.").
			WithPriority(0).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, triageIssues, ""); err != nil {
		return err
	}

	t.setupIssue = models.NewBaseIssue().
		WithTitle("Issue Triage Queue").
		WithDescription(issueTriageDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *IssueTriageTask) Validate(ctx context.Context) ValidationFeedback {
	expect := utils.NewExpector()

	return expect.Complete()
}
