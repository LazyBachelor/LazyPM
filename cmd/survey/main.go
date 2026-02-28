package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
	"github.com/LazyBachelor/LazyPM/internal/app"
	issues "github.com/LazyBachelor/LazyPM/internal/commands/issues"
	survey "github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/LazyBachelor/LazyPM/pkg/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
	"github.com/charmbracelet/fang"
)

func main() {
	ctx := context.Background()
	app, cleanup, err := initializeServices(ctx)
	if err != nil {
		return
	}
	defer cleanup()

	survey.SetApp(app)
	issues.SetApp(app)

	if err := fang.Execute(ctx, survey.RootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme)); err != nil {
		return
	}
}

func init() {
	task.RegisterInterface("tui", tui.New())
	task.RegisterInterface("web", web.New())
	task.RegisterInterface("repl", repl.New())

	task.RegisterTask("create_issue", func(app *app.App) task.Tasker {
		return tasks.NewCreateIssueTask(app)
	})
	task.RegisterTask("coding_task", func(app *app.App) task.Tasker {
		return tasks.NewCodingTask(app)
	})
	task.RegisterTask("git_task", func(app *app.App) task.Tasker {
		return tasks.NewGitTask(app)
	})
	task.RegisterTask("sprint_planning", func(app *app.App) task.Tasker {
		return tasks.NewSprintPlanningTask(app)
	})
	task.RegisterTask("issue_triage", func(app *app.App) task.Tasker {
		return tasks.NewIssueTriageTask(app)
	})
	task.RegisterTask("milestone_tracking", func(app *app.App) task.Tasker {
		return tasks.NewMilestoneTrackingTask(app)
	})
	task.RegisterTask("dependency_management", func(app *app.App) task.Tasker {
		return tasks.NewDependencyManagementTask(app)
	})
	task.RegisterTask("team_capacity", func(app *app.App) task.Tasker {
		return tasks.NewTeamCapacityTask(app)
	})
	task.RegisterTask("report_generation", func(app *app.App) task.Tasker {
		return tasks.NewReportGenerationTask(app)
	})
	task.RegisterTask("stakeholder_update", func(app *app.App) task.Tasker {
		return tasks.NewStakeholderUpdateTask(app)
	})
	task.RegisterTask("priority_management", func(app *app.App) task.Tasker {
		return tasks.NewPriorityManagementTask(app)
	})
	task.RegisterTask("backlog_refinement", func(app *app.App) task.Tasker {
		return tasks.NewBacklogRefinementTask(app)
	})

	// Basic survey commands
	survey.StartCmd.RunE = runStartCmd
	survey.RootCmd.AddCommand(survey.StartCmd)
	survey.RootCmd.AddCommand(survey.SubmitCmd)
	survey.RootCmd.AddCommand(survey.StatusCmd)
	survey.RootCmd.AddCommand(survey.ListTasksCmd)
	survey.RootCmd.AddCommand(survey.ListInterfacesCmd)
	survey.RootCmd.AddCommand(survey.IssuesCmd)

	// Issue related commands
	survey.IssuesCmd.AddCommand(issues.ListCmd)
	survey.IssuesCmd.AddCommand(issues.CreateCmd)
	survey.IssuesCmd.AddCommand(issues.UpdateCmd)
	survey.IssuesCmd.AddCommand(issues.DeleteCmd)
	survey.IssuesCmd.AddCommand(issues.CloseCmd)
	survey.IssuesCmd.AddCommand(issues.CommentCmd)
	survey.IssuesCmd.AddCommand(issues.CommentsCmd)

	// Issue commands for the REPL interface
	issues.RootCmd.AddCommand(issues.GetCmd)
	issues.RootCmd.AddCommand(issues.ListCmd)
	issues.RootCmd.AddCommand(issues.CloseCmd)
	issues.RootCmd.AddCommand(issues.CreateCmd)
	issues.RootCmd.AddCommand(issues.DeleteCmd)
	issues.RootCmd.AddCommand(issues.UpdateCmd)
	issues.RootCmd.AddCommand(issues.CommentCmd)
	issues.RootCmd.AddCommand(issues.CommentsCmd)
	issues.RootCmd.AddCommand(survey.StatusCmd)
}
