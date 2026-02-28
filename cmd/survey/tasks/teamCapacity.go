package tasks

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils"
	"github.com/charmbracelet/huh"
)

const teamCapacityDescription = `You are tasked with managing team capacity.

You have a team of 4 developers with varying skill and availability.
Review the upcoming sprint workload and:

1. Review assigned issues and their priorities
2. Identify overloaded team members based on issue assignments
3. Rebalance workload to match capacity
4. Consider vacations and time off mentioned in issue descriptions
5. Ensure critical tasks have coverage

Team member Alice is on vacation next week (mentioned in her issues). Bob can only work 50% time due to other commitments.`

type TeamCapacityTask struct {
	done       bool
	app        *App
	setupIssue *Issue
}

func NewTeamCapacityTask(app *App) *TeamCapacityTask {
	return &TeamCapacityTask{app: app, done: false}
}

func (t *TeamCapacityTask) Config() Config {
	return BaseConfig().WithStatisticsStoragePath("./.pm/capacity-task-stats.json")
}

func (t *TeamCapacityTask) Details() TaskDetails {
	return BaseDetails().
		WithTitle("Team Capacity Management").
		WithDescription(teamCapacityDescription).
		WithTimeToComplete("12m").
		WithDifficulty("Medium")
}

func (t *TeamCapacityTask) Questions(interfaceType InterfaceType) Questions {
	return BaseQuestions(interfaceType).With(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("How many issues did you reassign or update?").
				Options(
					huh.NewOption("0", 0),
					huh.NewOption("1-2", 1),
					huh.NewOption("3+", 2),
				),
		),
	)
}

func (t *TeamCapacityTask) Setup(ctx context.Context) error {
	if err := ClearIssues(t.app); err != nil {
		return err
	}

	capacityIssues := []*models.Issue{
		NewIssueBuilder().
			WithTitle("Authentication module").
			WithDescription("Implement OAuth2 flow. Assigned to: Alice (on vacation next week)").
			WithPriority(1).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Payment integration").
			WithDescription("Integrate Stripe API. Assigned to: Alice (on vacation next week)").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Dashboard widgets").
			WithDescription("Create reusable widget components. Assigned to: Bob (50% capacity)").
			WithPriority(2).
			WithStatus(models.StatusInProgress).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("API rate limiting").
			WithDescription("Add rate limiting middleware. Assigned to: Bob (50% capacity)").
			WithPriority(2).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Data migration").
			WithDescription("Migrate legacy data to new schema. Assigned to: Charlie (full capacity)").
			WithPriority(3).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
		NewIssueBuilder().
			WithTitle("Bug fixes batch").
			WithDescription("Fix reported bugs from QA. Assigned to: Diana (full capacity)").
			WithPriority(1).
			WithStatus(models.StatusOpen).
			WithIssueType(models.TypeTask).
			Build(),
	}

	if err := t.app.Issues.CreateIssues(ctx, capacityIssues, ""); err != nil {
		return err
	}

	t.setupIssue = NewIssueBuilder().
		WithTitle("Team Capacity Planning").
		WithDescription(teamCapacityDescription).
		Build()

	return t.app.Issues.CreateIssue(ctx, t.setupIssue, "")
}

func (t *TeamCapacityTask) Validate(ctx context.Context) ValidationFeedback {
	expect := utils.NewExpector()

	return expect.Complete()
}
