package main

import (
	"context"
	"os"

	"github.com/LazyBachelor/LazyPM/cmd/pm/tasks"
	"github.com/LazyBachelor/LazyPM/internal/app"
	issues "github.com/LazyBachelor/LazyPM/internal/commands/issues"
	survey "github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/LazyBachelor/LazyPM/pkg/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
	"github.com/spf13/cobra"
)

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

	cobra.EnableCommandSorting = false

	if os.Getenv("DEV") == "True" {
		issues.RootCmd.SetCompletionCommandGroupID("other")
	} else {
		issues.RootCmd.CompletionOptions.DisableDefaultCmd = true
	}

	setupLazyInitialization()

	RootCmd.AddGroup(&cobra.Group{ID: "issues", Title: "Issue Management"})
	issues.CreateCmd.GroupID = "issues"
	issues.GetCmd.GroupID = "issues"
	issues.ListCmd.GroupID = "issues"
	issues.UpdateCmd.GroupID = "issues"
	issues.CloseCmd.GroupID = "issues"
	issues.DeleteCmd.GroupID = "issues"

	RootCmd.AddCommand(issues.CreateCmd)
	RootCmd.AddCommand(issues.GetCmd)
	RootCmd.AddCommand(issues.ListCmd)
	RootCmd.AddCommand(issues.UpdateCmd)
	RootCmd.AddCommand(issues.CloseCmd)
	RootCmd.AddCommand(issues.DeleteCmd)

	RootCmd.AddGroup(&cobra.Group{ID: "comment", Title: "Comment Management"})
	issues.CommentCmd.GroupID = "comment"
	issues.CommentsCmd.GroupID = "comment"
	RootCmd.AddCommand(issues.CommentCmd)
	RootCmd.AddCommand(issues.CommentsCmd)

	var SurveyRootCmd = survey.RootCmd
	SurveyRootCmd.GroupID = "survey"
	RootCmd.AddGroup(&cobra.Group{ID: "survey", Title: "Survey Commands"})

	survey.StartCmd.RunE = runStartCmd
	SurveyRootCmd.AddCommand(survey.StartCmd)
	SurveyRootCmd.AddCommand(survey.StatusCmd)
	SurveyRootCmd.AddCommand(survey.SubmitCmd)
	SurveyRootCmd.AddCommand(survey.ListTasksCmd)
	SurveyRootCmd.AddCommand(survey.ListInterfacesCmd)

	RootCmd.AddCommand(SurveyRootCmd)

	RootCmd.AddGroup(&cobra.Group{ID: "other", Title: "Additional Commands"})
	RootCmd.SetHelpCommandGroupID("other")
	RootCmd.AddCommand(replCmd)

}

var replCmd = &cobra.Command{
	Use:     "repl",
	GroupID: "other",
	Short:   "Start the interactive REPL interface",
	Long:    `Start the interactive Read-Eval-Print Loop (REPL) for managing your projects and issues in an interactive terminal environment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return repl.New().Run(cmd.Context(), App.Config)
	},
}

func initializeApp(ctx context.Context) (*app.App, func(), error) {
	return app.New(ctx, tasks.BaseConfig().WithAutoInit(true))
}

func setupLazyInitialization() {
	prevPreRun := RootCmd.PersistentPreRun
	prevPreRunE := RootCmd.PersistentPreRunE

	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if commandNeedsApp(cmd) {
			if err := ensureAppInitialized(cmd.Context()); err != nil {
				return err
			}
		}

		if prevPreRunE != nil {
			if err := prevPreRunE(cmd, args); err != nil {
				return err
			}
		}

		if prevPreRun != nil {
			prevPreRun(cmd, args)
		}

		return nil
	}
}

func commandNeedsApp(cmd *cobra.Command) bool {
	name := cmd.Name()
	switch name {
	case "help", "completion", "__complete", "__completeNoDesc":
		return false
	}

	if cmd == RootCmd || name == "survey" {
		return false
	}

	if parent := cmd.Parent(); parent != nil && parent.Name() == "survey" {
		switch name {
		case "tasks", "interfaces":
			return false
		}
	}

	return true
}

func ensureAppInitialized(ctx context.Context) error {
	if App != nil {
		return nil
	}

	application, cleanup, err := initializeApp(ctx)
	if err != nil {
		return err
	}

	App = application
	appCleanup = cleanup
	survey.SetApp(App)
	issues.SetApp(App)

	return nil
}

func initInterfaces() map[string]task.Interface {
	interfaces := make(map[string]task.Interface)
	for _, name := range task.ListInterfaces() {
		i, err := task.GetInterface(name)
		if err != nil {
			continue
		}
		interfaces[name] = i
	}
	return interfaces
}

func initTasks(app *app.App) map[string]task.Tasker {
	taskMap := make(map[string]task.Tasker)
	for _, name := range task.ListTasks() {
		t, err := task.GetTask(name, app)
		if err != nil {
			continue
		}
		taskMap[name] = t
	}

	return taskMap
}
