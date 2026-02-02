package commands

import (
	"github.com/LazyBachelor/LazyPM/internal/models"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long:  `Create a new issue with the specified details.`,
	RunE:  runCreateCmd,
}

func runCreateCmd(cmd *cobra.Command, args []string) error {
	issue := &models.Issue{
		Title:       "test",
		Description: "This is a test issue created by the CLI.",
		Status:      models.StatusOpen,
		IssueType:   models.TypeBug,
		Priority:    0,
	}

	err := svc.Beads.CreateIssue(cmd.Context(), issue, "test_actor")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating issue: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created issue: %s\n", issue.ID)
	return nil
}
