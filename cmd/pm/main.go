package main

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/app"
	issues "github.com/LazyBachelor/LazyPM/internal/commands/issues"
	"github.com/charmbracelet/fang"
)

var App *app.App
var appCleanup func()
var RootCmd = issues.RootCmd
var DB_URI string

func main() {
	ctx := context.Background()
	defer func() {
		if appCleanup != nil {
			appCleanup()
		}
	}()

	if err := fang.Execute(ctx, RootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme)); err != nil {
		return
	}
}
