package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
	"github.com/Dicklesworthstone/beads_viewer/pkg/ui"
	"github.com/LazyBachelor/LazyPM/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config := service.Config{
		StatisticsStoragePath: "./.pm/stats.json",
		BeadsDBPath:           "./.pm/db.db",
		IssuePrefix:           "pm",
	}

	svc, close, err := service.NewServices(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	beadsIssues, err := svc.Beads.AllIssues(context.Background())
	if err != nil {
		panic(err)
	}

	modelIssues := make([]model.Issue, 0, len(beadsIssues))
	for _, issue := range beadsIssues {
		jsonData, err := json.Marshal(issue)
		if err != nil {
			panic(err)
		}

		var modelIssue model.Issue
		if err := json.Unmarshal(jsonData, &modelIssue); err != nil {
			panic(err)
		}
		modelIssues = append(modelIssues, modelIssue)
	}

	model := ui.NewModel(modelIssues, nil, config.BeadsDBPath)

	if err := tea.NewProgram(model, tea.WithAltScreen(),
		tea.WithMouseAllMotion()); err != nil {
		panic(err)
	}

}
