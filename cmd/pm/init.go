package main

import (
	"context"
	"os"

	"github.com/LazyBachelor/LazyPM/cmd/pm/tasks"
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/commands/issues"
	"github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	godotenv.Load(".env")

	models.BaseConfig = models.BaseConfig.LoadFromEnv()

	if DB_URI != "" {
		models.BaseConfig = models.BaseConfig.WithDbUri(DB_URI)
	}

	task.RegisterInterface("tui", tui.New())
	task.RegisterInterface("web", web.New())
	task.RegisterInterface("repl", repl.New())

	task.RegisterTask("backlog_refinement", func(app *app.App) task.Tasker {
		return tasks.NewBacklogRefinementTask(app)
	})
	task.RegisterTask("create_issue", func(app *app.App) task.Tasker {
		return tasks.NewCreateIssueTask(app)
	})
	task.RegisterTask("coding_task", func(app *app.App) task.Tasker {
		return tasks.NewCodingTask(app)
	})
	task.RegisterTask("sprint_planning", func(app *app.App) task.Tasker {
		return tasks.NewSprintPlanningTask(app)
	})
	task.RegisterTask("priority_management", func(app *app.App) task.Tasker {
		return tasks.NewPriorityManagementTask(app)
	})
	task.RegisterTask("dependency_management", func(app *app.App) task.Tasker {
		return tasks.NewDependencyManagementTask(app)
	})
	task.RegisterTask("issue_review_cleanup", func(app *app.App) task.Tasker {
		return tasks.NewIssueReviewCleanupTask(app)
	})
	task.RegisterTask("git_task", func(app *app.App) task.Tasker {
		return tasks.NewGitTask(app)
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
	issues.DepCmd.GroupID = "issues"

	RootCmd.AddCommand(issues.CreateCmd)
	RootCmd.AddCommand(issues.GetCmd)
	RootCmd.AddCommand(issues.ListCmd)
	RootCmd.AddCommand(issues.UpdateCmd)
	RootCmd.AddCommand(issues.CloseCmd)
	RootCmd.AddCommand(issues.DeleteCmd)
	RootCmd.AddCommand(issues.DepCmd)

	RootCmd.AddGroup(&cobra.Group{ID: "comment", Title: "Comment Management"})
	issues.CommentCmd.GroupID = "comment"
	issues.CommentsCmd.GroupID = "comment"
	RootCmd.AddCommand(issues.CommentCmd)
	RootCmd.AddCommand(issues.CommentsCmd)

	RootCmd.AddGroup(&cobra.Group{ID: "survey", Title: "Survey Commands"})
	survey.StartCmd.GroupID = "survey"
	survey.StatusCmd.GroupID = "survey"
	survey.SubmitCmd.GroupID = "survey"
	survey.ListTasksCmd.GroupID = "survey"
	survey.ListInterfacesCmd.GroupID = "survey"

	survey.StartCmd.RunE = runStartCmd
	RootCmd.AddCommand(survey.StartCmd)
	RootCmd.AddCommand(survey.StatusCmd)
	RootCmd.AddCommand(survey.SubmitCmd)

	RootCmd.AddCommand(survey.ListTasksCmd)
	RootCmd.AddCommand(survey.ListInterfacesCmd)

	RootCmd.AddGroup(&cobra.Group{ID: "other", Title: "Additional Commands"})
	RootCmd.SetHelpCommandGroupID("other")
	RootCmd.AddCommand(replCmd)
	RootCmd.AddCommand(tuiCmd)
	RootCmd.AddCommand(webCmd)

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

var tuiCmd = &cobra.Command{
	Use:     "tui",
	GroupID: "other",
	Short:   "Start the interactive TUI interface",
	Long:    `Start the interactive Terminal User Interface (TUI) for managing your projects and issues in an interactive terminal environment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.New().Run(cmd.Context(), App.Config)
	},
}

var webCmd = &cobra.Command{
	Use:     "web",
	GroupID: "other",
	Short:   "Start the interactive web interface",
	Long:    `Start the interactive web interface for managing your projects and issues in a web browser.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return web.New().Run(cmd.Context(), App.Config)
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
	case "help", "completion":
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
		survey.SetApp(App)
		issues.SetApp(App)
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
