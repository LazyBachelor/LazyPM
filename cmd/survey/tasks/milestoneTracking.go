package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils"
	"github.com/charmbracelet/huh"
)

const milestoneTrackingDescription = `You are tasked with managing a project milestone.

The "Q1 Release" milestone is approaching and you need to ensure all issues are on track.
Review the milestone issues and:

1. Identify issues at risk of missing the deadline
2. Update issue statuses to reflect current progress
3. Flag any blockers or dependencies causing delays
4. Close completed issues
5. Provide status updates for stakeholder visibility

The milestone deadline is 2 weeks away. Some issues have dependencies that need to be resolved first.`

type MilestoneTrackingTask struct {
	done       bool
	app        *App
	setupIssue *Issue
}

func NewMilestoneTrackingTask(app *App) *MilestoneTrackingTask {
	return &MilestoneTrackingTask{app: app, done: false}
}

func (t *MilestoneTrackingTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/milestone-task-stats.json")
}

func (t *MilestoneTrackingTask) Details() TaskDetails {
	return BaseDetails().
		WithTitle("Milestone Tracking Task").
		WithDescription(milestoneTrackingDescription).
		WithTimeToComplete("10m").
		WithDifficulty("Easy")
}

func (t *MilestoneTrackingTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("How many issues did you identify as at-risk?").
				Options(
					huh.NewOption("0", 0),
					huh.NewOption("1-2", 1),
					huh.NewOption("3+", 2),
				),
		),
	)
}

func (t *MilestoneTrackingTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	milestoneIssues := []*models.Issue{
		NewIssueBuilder().
			WithTitle("Core API endpoints").
			WithDescription("Implement REST API for core functionality. COMPLETED").
			WithPriority(1).
			WithStatus(models.StatusClosed).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Frontend dashboard").
			WithDescription("Create main dashboard UI. In progress - 80% complete").
			WithPriority(1).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("User management module").
			WithDescription("Depends on Core API. Add user CRUD operations. Currently blocked").
			WithPriority(2).
			WithStatus(models.StatusBlocked).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Analytics reporting").
			WithDescription("Generate usage analytics reports. Not started yet").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Email notifications").
			WithDescription("Setup email service for alerts. 50% complete").
			WithPriority(2).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, milestoneIssues, ""); err != nil {
		return err
	}

	t.setupIssue = NewIssueBuilder().
		WithTitle("Q1 Release Milestone Tracking").
		WithDescription(milestoneTrackingDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *MilestoneTrackingTask) Validate(ctx context.Context) ValidationFeedback {
	expect := utils.NewExpector()

	return expect.Complete()
}
