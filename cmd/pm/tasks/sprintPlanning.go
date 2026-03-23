package tasks

import (
	"context"
	"fmt"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/check"
)

const sprintPlanningDescription = `You are tasked with sprint planning.

A new sprint is starting and you need to organize the backlog.
You will be given a list of issues with different priorities and dependencies.

Your task:
1. Review the backlog issues
2. Select which issues to include in the sprint
3. Add issues to the sprint
4. Address any blocked or dependent issues
5. Prioritize high-priority items

The goal is to create a realistic sprint plan that delivers value while respecting team capacity.`

type SprintPlanningTask struct {
	done      bool
	app       *App
	sprintNum int
}

func NewSprintPlanningTask(app *App) *SprintPlanningTask {
	return &SprintPlanningTask{app: app, done: false}
}

func (t *SprintPlanningTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/sprint-planning-stats.json")
}

func (t *SprintPlanningTask) Details(interfaceType InterfaceType) TaskDetails {
	return BaseDetails(interfaceType).
		WithTitle("Sprint Planning Task").
		WithDescription(sprintPlanningDescription).
		WithTimeToComplete("3m").
		WithDifficulty("Medium")
}

func (t *SprintPlanningTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Key("sprint-planning-selected-issues").
				Title("How difficult was it selecting issues for the sprint?").
				Description("We are interested in how intuitive and time-consuming the selection process was.").
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

func (t *SprintPlanningTask) QuestionnaireKeys(_ InterfaceType) []string {
	return BaseKeys().With("sprint-planning-selected-issues")
}

func (t *SprintPlanningTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	sprintNum, err := t.app.Issues.AddSprint(ctx)
	if err != nil {
		return err
	}
	t.sprintNum = sprintNum

	backlogIssues := []*models.Issue{
		NewIssueBuilder().
			WithTitle("Implement user authentication").
			WithDescription("Add login/logout functionality. Priority: High").
			WithPriority(4).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Design database schema").
			WithDescription("Create tables for users and orders. Priority: High").
			WithPriority(4).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Setup CI/CD pipeline").
			WithDescription("Configure automated testing and deployment. Currently blocked by server setup").
			WithPriority(3).
			WithStatus(models.StatusBlocked).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Create API documentation").
			WithDescription("Document all REST endpoints. Priority: Low").
			WithPriority(1).
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

	return t.app.Issues.CreateIssues(ctx, backlogIssues, "")
}

func (t *SprintPlanningTask) Validate(ctx context.Context) ValidationFeedback {
	expect := check.NewExpector()

	if t.sprintNum == 0 {
		return expect.Fatal("No sprint was created. Create a sprint first.")
	}

	sprintIssues, err := t.app.Issues.GetIssuesBySprint(ctx, t.sprintNum)
	if err != nil {
		return expect.Fatal("Could not fetch sprint issues")
	}

	allIssues, err := FetchIssues(ctx, t.app)
	if err != nil {
		return expect.Fatal("Could not fetch issues")
	}

	expect.Assert(len(sprintIssues) > 0, fmt.Sprintf("Sprint %d has no issues. Move issues from Backlog to the sprint.", t.sprintNum))

	sorted := make([]*models.Issue, len(allIssues))
	copy(sorted, allIssues)
	for i := range sorted {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Priority > sorted[i].Priority {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	top3 := sorted[:min(len(sorted), 3)]
	sprintIssueIDs := make(map[string]bool)
	for _, issue := range sprintIssues {
		sprintIssueIDs[issue.ID] = true
	}

	var topInSprint int
	for _, issue := range top3 {
		if sprintIssueIDs[issue.ID] {
			topInSprint++
		}
	}

	expect.Assert(topInSprint >= 2,
		fmt.Sprintf("Only %d of 3 high-priority issues in sprint.\n High priority: \n- %s \n- %s \n- %s",
			topInSprint, top3[0].Title, top3[1].Title, top3[2].Title))

	expect.Assert(len(sprintIssues) >= 3,
		fmt.Sprintf("Sprint %d only has %d issues. Add at least %d more from backlog.",
			t.sprintNum, len(sprintIssues), 3-len(sprintIssues)))

	return expect.Complete()
}
