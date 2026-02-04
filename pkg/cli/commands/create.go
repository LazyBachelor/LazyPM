package commands

import (
	"fmt"
	"os"

	"github.com/LazyBachelor/LazyPM/internal/models"

	"github.com/spf13/cobra"
)

var (
	createTitle       string
	createDescription string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long:  `Create a new issue with the specified details.`,
	RunE:  runCreateCmd,
}

func init() {
	createCmd.Flags().StringVar(&createTitle, "title", "", "Issue title")
	_ = createCmd.MarkFlagRequired("title")
	createCmd.Flags().StringVar(&createDescription, "desc", "", "Issue description")
}

func runCreateCmd(cmd *cobra.Command, args []string) error {
	issue := &models.Issue{
		Title:       createTitle,
		Description: createDescription,
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
