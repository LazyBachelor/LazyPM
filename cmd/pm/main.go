package main

import (
	"context"

	"charm.land/fang/v2"
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/commands/issues"
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
