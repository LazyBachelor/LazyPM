package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
)

func main() {
	ctx := context.Background()

	config := service.Config{
		IssuePrefix:           "pm",
		BeadsDBPath:           "./.pm/db.db",
		StatisticsStoragePath: "./.pm/stats.json",
	}
	svc, cleanup, err := service.NewServices(ctx, config)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer cleanup()

	issue := &models.Issue{
		IssueType:   models.TypeTask,
		Title:       "Sample Issue",
		Description: "This is a sample issue created for testing.",
		Status:      models.StatusOpen,
	}

	err = svc.Beads.CreateIssue(ctx, issue, "")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fetchedIssues, err := svc.Beads.SearchIssues(ctx, "", models.IssueFilter{})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for _, iss := range fetchedIssues {
		fmt.Printf("Issue ID: %s, Title: %s, Status: %s\n", iss.ID, iss.Title, iss.Status)
	}

	stats, err := svc.Statistics.GetStatistics()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("\nStatistics: %v\n", stats)
}
