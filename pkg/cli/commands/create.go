package commands

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"

	"github.com/spf13/cobra"
)

var (
	createDescription string
	createStatus      string
	createType        string
	createPriority    int
)

var createCmd = &cobra.Command{
	Use:   "issue [title]",
	Short: "Create a new issue",
	Long:  `Create a new issue with the specified details.`,
	RunE:  runCreateCmd,
	Args:  cobra.ExactArgs(1),
}

func init() {
	createCmd.Flags().StringVarP(&createDescription, "desc", "d", "", "Issue description")
	createCmd.Flags().StringVarP(&createStatus, "status", "s", "open", "Issue status(open, closed, in_progress)")
	createCmd.Flags().StringVarP(&createType, "type", "t", "task", "Issue type(bug, feature, task)")
	createCmd.Flags().IntVarP(&createPriority, "priority", "p", 0, "Issue priority(0-5)")
}

func runCreateCmd(cmd *cobra.Command, args []string) error {
	createTitle := args[0]

	if createTitle == "" {
		return fmt.Errorf("issue title cannot be empty")
	}

	issue := &models.Issue{
		Title:       createTitle,
		Description: createDescription,
		Status:      models.Status(createStatus),
		IssueType:   models.IssueType(createType),
		Priority:    createPriority,
	}

	err := svc.Beads.CreateIssue(cmd.Context(), issue, "test_actor")
	if err != nil {
		return fmt.Errorf("error creating issue: %w", err)
	}

	str := fmt.Sprintf("Created issue with ID: %s\n", issue.ID)

	if issue.Title != "" {
		str += fmt.Sprintf("Title: %s\n", issue.Title)
	}

	if issue.Description != "" {
		str += fmt.Sprintf("Description: %s\n", issue.Description)
	}

	if issue.Status != "" {
		str += fmt.Sprintf("Status: %s\n", issue.Status)
	}

	if issue.IssueType != "" {
		str += fmt.Sprintf("Type: %s\n", issue.IssueType)
	}

	if issue.Priority != 0 {
		str += fmt.Sprintf("Priority: %d\n", issue.Priority)
	}

	fmt.Print(str)

	return nil
}
