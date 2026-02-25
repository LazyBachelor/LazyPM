package surveyCmd

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/spf13/cobra"
)

const appKey string = "app"
const long string = `Project Management Interface Survey

Thank you for participating in our survey!

We are gathering data on how users interact with different task management interfaces to better understand their preferences and compare their usability.
This survey will present you with a series of tasks to complete using various interfaces, including command-line, web-based, and terminal user interfaces.

Please answer the questions honestly and to the best of your ability.
Your responses will be kept confidential and used solely for research purposes.`

var app *service.App

// RootCmd is the base command for the survey CLI.
var RootCmd = &cobra.Command{
	Use:  "survey",
	Long: long,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if app != nil {
			cmd.SetContext(context.WithValue(cmd.Context(), appKey, app))
		}
	},
}

func SetApp(application *service.App) {
	app = application
}

func AppFromContext(ctx context.Context) *service.App {
	if a, ok := ctx.Value(appKey).(*service.App); ok {
		return a
	}
	return app
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
}
