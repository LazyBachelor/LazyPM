package main

import (
	"context"
	"log"

	"github.com/charmbracelet/fang"
)

func main() {
	ctx := context.Background()

	if err := fang.Execute(ctx, rootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme)); err != nil {
		log.Fatalf("Failed to execute command: %v\n", err)
	}
}
