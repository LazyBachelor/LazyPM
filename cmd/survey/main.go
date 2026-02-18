package main

import (
	"context"

	"github.com/charmbracelet/fang"
)

func main() {
	ctx := context.Background()

	if err := fang.Execute(ctx, rootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme)); err != nil {
		return
	}
}
