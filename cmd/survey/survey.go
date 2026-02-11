package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/LazyBachelor/LazyPM/cmd/survey/forms"
	"github.com/LazyBachelor/LazyPM/pkg"
	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
	"github.com/charmbracelet/huh"
)

//go:embed assets/*
var assetsFS embed.FS

func main() {
	intro, err := forms.NewIntroduction(assetsFS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading introduction: %v\n", err)
		os.Exit(1)
	}

	if err := intro.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running introduction: %v\n", err)
		os.Exit(1)
	}

	var selected string
	prompt := huh.NewSelect[string]().Options(
		huh.NewOption("Start CLI in REPL Mode", "repl"),
		huh.NewOption("Start Web Interface", "web"),
		huh.NewOption("Start TUI Interface", "tui"),
	).Value(&selected).WithTheme(huh.ThemeBase16())

	if err := prompt.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running prompt: %v\n", err)
		os.Exit(1)
	}

	config := pkg.SurveyConfig{
		RootCmd:               "pm",
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/stats.json",
		WebAddress:            "localhost:8080",
	}

	ctx := context.Background()
	switch selected {
	case "repl":
		fmt.Println("Starting CLI in REPL mode...")
		err = repl.RunREPL(ctx, config)
	case "web":
		fmt.Println("Starting Web Interface...")
		err = web.Run(ctx, config)
	case "tui":
		fmt.Println("Starting TUI Interface...")
		_, err = tui.Run(ctx, config)
	default:
		fmt.Println("Invalid selection. Exiting.")
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running selected interface: %v\n", err)
		os.Exit(1)
	}

}
