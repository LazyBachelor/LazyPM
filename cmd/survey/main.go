package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
	issuesCmd "github.com/LazyBachelor/LazyPM/internal/commands/issues"
	surveyCmd "github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/LazyBachelor/LazyPM/internal/service"
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

	surveyCmd.SetApp(app)
	issuesCmd.SetApp(app)

	if err := fang.Execute(ctx, surveyCmd.RootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme)); err != nil {
		return
	}
}

func init() {
	task.RegisterInterface("tui", tui.NewTui())
	task.RegisterInterface("web", web.NewWeb())
	task.RegisterInterface("repl", repl.NewRepl())

	task.RegisterTask("create_issue", func(app *service.App) task.Tasker {
		return tasks.NewCreateIssueTask(app)
	})
	task.RegisterTask("coding_task", func(app *service.App) task.Tasker {
		return tasks.NewCodingTask(app)
	})
	task.RegisterTask("git_task", func(app *service.App) task.Tasker {
		return tasks.NewGitTask(app)
	})
	task.RegisterTask("sprint_planning", func(app *service.App) task.Tasker {
		return tasks.NewSprintPlanningTask(app)
	})
	task.RegisterTask("issue_triage", func(app *service.App) task.Tasker {
		return tasks.NewIssueTriageTask(app)
	})
	task.RegisterTask("milestone_tracking", func(app *service.App) task.Tasker {
		return tasks.NewMilestoneTrackingTask(app)
	})
	task.RegisterTask("dependency_management", func(app *service.App) task.Tasker {
		return tasks.NewDependencyManagementTask(app)
	})
	task.RegisterTask("team_capacity", func(app *service.App) task.Tasker {
		return tasks.NewTeamCapacityTask(app)
	})
	task.RegisterTask("report_generation", func(app *service.App) task.Tasker {
		return tasks.NewReportGenerationTask(app)
	})
	task.RegisterTask("stakeholder_update", func(app *service.App) task.Tasker {
		return tasks.NewStakeholderUpdateTask(app)
	})
	task.RegisterTask("priority_management", func(app *service.App) task.Tasker {
		return tasks.NewPriorityManagementTask(app)
	})
	task.RegisterTask("backlog_refinement", func(app *service.App) task.Tasker {
		return tasks.NewBacklogRefinementTask(app)
	})

	// Basic survey commands
	surveyCmd.StartCmd.RunE = runStartCmd
	surveyCmd.RootCmd.AddCommand(surveyCmd.StartCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.SubmitCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.StatusCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.ListTasksCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.ListInterfacesCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.IssuesCmd)

	// Issue related commands
	surveyCmd.IssuesCmd.AddCommand(issuesCmd.ListCmd)
	surveyCmd.IssuesCmd.AddCommand(issuesCmd.CreateCmd)
	surveyCmd.IssuesCmd.AddCommand(issuesCmd.UpdateCmd)
	surveyCmd.IssuesCmd.AddCommand(issuesCmd.DeleteCmd)
	surveyCmd.IssuesCmd.AddCommand(issuesCmd.CloseCmd)
	surveyCmd.IssuesCmd.AddCommand(issuesCmd.CommentCmd)
	surveyCmd.IssuesCmd.AddCommand(issuesCmd.CommentsCmd)

	// Issue commands for the REPL interface
	issuesCmd.RootCmd.AddCommand(issuesCmd.GetCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.ListCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.CloseCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.CreateCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.DeleteCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.UpdateCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.CommentCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.CommentsCmd)
	issuesCmd.RootCmd.AddCommand(surveyCmd.StatusCmd)
}
