package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
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

	if err := fang.Execute(ctx, surveyCmd.RootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme)); err != nil {
		return
	}
}

func init() {
	task.RegisterInterface("tui", tui.NewTui())
	task.RegisterInterface("web", web.NewWeb())
	task.RegisterInterface("repl", repl.NewRepl())

	surveyCmd.StartCmd.RunE = runStartCmd
	surveyCmd.RootCmd.AddCommand(surveyCmd.StartCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.SubmitCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.StatusCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.ListTasksCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.ListInterfacesCmd)

	task.RegisterTask("create_issue", func(app *service.App) task.Tasker {
		return tasks.NewCreateIssueTask(app)
	})
	task.RegisterTask("coding_task", func(app *service.App) task.Tasker {
		return tasks.NewCodingTask(app)
	})
	task.RegisterTask("git-task", func(app *service.App) task.Tasker {
		return tasks.NewGitTask(app)
	})

}
