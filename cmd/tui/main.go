package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
	bv_ui "github.com/Dicklesworthstone/beads_viewer/pkg/ui"
	"github.com/LazyBachelor/LazyPM/internal/models"
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
		fmt.Fprintf(os.Stderr, "Error initializing services: %v\n", err)
		os.Exit(1)
	}
	defer close()

	beadsIssues, err := svc.Beads.AllIssues(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading issues: %v\n", err)
		os.Exit(1)
	}

	issues, err := convertToModelIssues(beadsIssues)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting issues: %v\n", err)
		os.Exit(1)
	}

	model := bv_ui.NewModel(issues, nil, config.BeadsDBPath)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func convertToModelIssues(beadsIssues []models.Issue) ([]model.Issue, error) {
	result := make([]model.Issue, 0, len(beadsIssues))

	for _, issue := range beadsIssues {

		jsonData, err := json.Marshal(issue)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal issue ID %s: %w", issue.ID, err)
		}

		var modelIssue model.Issue
		if err := json.Unmarshal(jsonData, &modelIssue); err != nil {
			return nil, fmt.Errorf("failed to unmarshal issue ID %s: %w", issue.ID, err)
		}

		result = append(result, modelIssue)
	}
	return result, nil
}
