package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/steveyegge/beads"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long:  `Create a new issue with the specified details.`,
	RunE:  runCreateCmd,
}

func runCreateCmd(cmd *cobra.Command, args []string) error {
	issue := &beads.Issue{
		Title:       "test",
		Description: "This is a test issue created by the CLI.",
		Status:      beads.StatusOpen,
		IssueType:   beads.TypeBug,
		Priority:    0,
	}

	err := svc.CreateIssue(cmd.Context(), issue, "test_actor")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating issue: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created issue: %s\n", issue.ID)
	return nil
}
