package tasks

import (
	"context"
	"sort"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
	"github.com/charmbracelet/huh"
)

const sprintPlanningDescription = `Spring planning task

You are tasked with sprint planning.

A new sprint is starting and you need to organize the backlog.
You will be given a list of issues with different priorities and dependencies.

Your task: 

1. Open the backlog issues
2. Select the issue with the highest priority firstly to include in the sprint (1 is highest)
3. Update issue statuses for the 5 highest priority issues to "Ready to sprint" and save
4. Find the blocked issue and add issue to the sprint
5. Prioritize 6 high-priority items in the sprint

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
				Key("sprint-selected-issues").
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
	return []string{"task_completed", "task_difficulty", "sprint-selected-issues"}
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
		expect.Fail("No issues found for sprint planning")
		return expect.ValidationFeedback
	}

	var blockedIssue *models.Issue
	for _, issue := range issues {
		if strings.Contains(strings.ToLower(issue.Description), "currently blocked") {
			blockedIssue = issue
			break
		}
	}

	expect.Assert(blockedIssue != nil, "Expected a blocked issue to exist and be added to the sprint")

	if blockedIssue != nil {
		expect.Assert(blockedIssue.Status != models.StatusBlocked,
			"The blocked issue should have been added to the sprint and no longer be blocked")
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Priority < issues[j].Priority
	})

	topN := 5
	if len(issues) < topN {
		topN = len(issues)
	}

	readyStatus := models.Status("ready_to_sprint")
	readyCount := 0

	for i := 0; i < topN; i++ {
		if issues[i].Status == readyStatus {
			readyCount++
		}
	}

	expect.Assert(readyCount == topN,
		"Expected the 6 highest-priority issues to have status 'Ready to sprint'")

	return expect.Complete()
}
