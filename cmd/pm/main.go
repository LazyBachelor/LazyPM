package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/app"
	issues "github.com/LazyBachelor/LazyPM/internal/commands/issues"
	survey "github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/charmbracelet/fang"
)

var App *app.App
var RootCmd = issues.RootCmd

func main() {
	ctx := context.Background()
	app, cleanup, err := initializeServices(ctx)
	if err != nil {
		return
	}
	defer cleanup()

	App = app
	survey.SetApp(App)
	issues.SetApp(App)

	if err := fang.Execute(ctx, RootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme)); err != nil {
		return
	}
}
