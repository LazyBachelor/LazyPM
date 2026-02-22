package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/cmd/survey/tasks"
	"github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
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
	surveyCmd.StartCmd.RunE = runStartCmd
	surveyCmd.RootCmd.AddCommand(surveyCmd.StartCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.SubmitCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.StatusCmd)
	surveyCmd.RootCmd.AddCommand(surveyCmd.ListCmd)

	task.Register("create_issue", func(app *service.App) task.Tasker {
		return tasks.NewCreateIssueTask(app)
	})
	task.Register("coding_task", func(app *service.App) task.Tasker {
		return tasks.NewCodingTask(app)
	})
}
