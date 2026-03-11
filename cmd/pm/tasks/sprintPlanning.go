package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
	"github.com/charmbracelet/huh"
)

const sprintPlanningDescription = `You are tasked with sprint planning.

A new sprint is starting and you need to organize the backlog.
You will be given a list of issues with different priorities and dependencies.

Your task:
1. Review the backlog issues
2. Select which issues to include in the sprint
3. Update issue statuses to move items into the sprint
4. Address any blocked or dependent issues
5. Prioritize high-priority items

The goal is to create a realistic sprint plan that delivers value while respecting team capacity.`

type SprintPlanningTask struct {
	done       bool
	app        *App
	setupIssue *Issue
}

func NewSprintPlanningTask(app *App) *SprintPlanningTask {
	return &SprintPlanningTask{app: app, done: false}
}

func (t *SprintPlanningTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/sprint-planning-stats.json")
}

func (t *SprintPlanningTask) Details() TaskDetails {
	return BaseDetails().
		WithTitle("Sprint Planning Task").
		WithDescription(sprintPlanningDescription).
		WithTimeToComplete("15m").
		WithDifficulty("Medium")
}

func (t *SprintPlanningTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Key("sprint-planning-selected-issues").
				Title("How many issues did you select for the sprint?").
				Options(
					huh.NewOption("2-3 issues", 1),
					huh.NewOption("4-5 issues", 2),
					huh.NewOption("6+ issues", 3),
				),
		),
	)
}

func (t *SprintPlanningTask) QuestionnaireKeys(_ InterfaceType) []string {
	return []string{"task_completed", "task_difficulty", "sprint-planning-selected-issues"}
}

func (t *SprintPlanningTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	backlogIssues := []*models.Issue{
		NewIssueBuilder().
			WithTitle("Implement user authentication").
			WithDescription("Add login/logout functionality. Priority: High").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Design database schema").
			WithDescription("Create tables for users and orders. Priority: High").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Setup CI/CD pipeline").
			WithDescription("Configure automated testing and deployment. Currently blocked by server setup").
			WithPriority(2).
			WithStatus(models.StatusBlocked).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Create API documentation").
			WithDescription("Document all REST endpoints. Priority: Low").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Implement search functionality").
			WithDescription("Depends on database schema. Priority: Medium").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, backlogIssues, ""); err != nil {
		return err
	}

	t.setupIssue = NewIssueBuilder().
		WithTitle("Sprint Planning - Week 1").
		WithDescription(sprintPlanningDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *SprintPlanningTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	issues, err := FetchIssues(ctx, t.app, t.setupIssue)
	if err != nil {
		return expect.ValidationFeedback
	}

	if len(issues) == 0 {
		expect.Fail("No backlog issues found to plan a sprint with.")
		return expect.ValidationFeedback
	}

	// Sort by priority ascending (0 is highest priority).
	sorted := make([]*models.Issue, len(issues))
	copy(sorted, issues)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Priority < sorted[i].Priority {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	topN := 5
	if len(sorted) < topN {
		topN = len(sorted)
	}
	top := sorted[:topN]

	var plannedCount int
	var readyToSprintCount int
	for _, issue := range top {
		if issue.Status == models.StatusReadyToSprint ||
			issue.Status == models.StatusInProgress ||
			issue.Status == models.StatusClosed {
			plannedCount++
		}
		if issue.Status == models.StatusReadyToSprint {
			readyToSprintCount++
		}
	}

	expect.Assert(plannedCount >= 3,
		"Expected at least 3 of the 5 highest-priority issues to be moved into 'ready_to_sprint', 'in_progress', or 'closed' for the sprint.")
	expect.Assert(readyToSprintCount >= 1,
		"Expected at least one of the highest-priority issues to be marked as 'ready_to_sprint' to indicate it is planned for the sprint.")

	return expect.Complete()
}
