package survey

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
)

type App = models.App

var app *App

const appKey string = "app"
const long string = `Project Management Interface Survey

Thank you for participating in our survey!

We are gathering data on how users interact with different task management interfaces to better understand their preferences and compare their usability.
This survey will present you with a series of tasks to complete using various interfaces, including command-line, web-based, and terminal user interfaces.

Please answer the questions honestly and to the best of your ability.
Your responses will be kept confidential and used solely for research purposes.`

// RootCmd is the base command for the survey CLI.
var RootCmd = &cobra.Command{
	Use:  "survey",
	Long: long,
}

func SetApp(application *App) {
	app = application
}

func AppFromContext(ctx context.Context) *App {
	if a, ok := ctx.Value(appKey).(*App); ok {
		return a
	}
	return app
}

func init() {
}
