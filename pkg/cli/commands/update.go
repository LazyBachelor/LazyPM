package commands

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
)

var updateFlags Flags

var updateCmd = &cobra.Command{
	Use:               "update [issue ID]",
	Short:             "Update an existing issue",
	Long:              `Update an existing issue by its ID with the specified details.`,
	Example:           `pm update pm-001 --title "New title" -d "Description" -s in_progress --type task -p 3`,
	RunE:              runUpdateCmd,
	Aliases:           []string{"edit"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeIssues,
}

func runUpdateCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	issue, err := svc.Beads.GetIssue(cmd.Context(), issueID)
	if err != nil || issue == nil {
		return fmt.Errorf("error getting issue: %w", err)
	}

	updates, err := getUpdateValues(cmd)
	if err != nil {
		return fmt.Errorf("error getting update values: %w", err)
	}

	err = svc.Beads.UpdateIssue(cmd.Context(), issueID, updates, "test_actor")
	if err != nil {
		return fmt.Errorf("error updating issue: %w", err)
	}

	updatedIssue, err := svc.Beads.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return fmt.Errorf("error getting updated issue: %w", err)
	}

	cmd.Printf("Updated issue to:\n%s", models.IssueString(*updatedIssue))

	return nil
}

func init() {
	updateCmd.Flags().StringVar(&updateFlags.title, "title", "", "New issue title")
	updateCmd.Flags().StringVarP(&updateFlags.description, "desc", "d", "", "New issue description")
	updateCmd.Flags().StringVarP(&updateFlags.status, "status", "s", "", "New issue status(open, closed, in_progress)")
	updateCmd.Flags().StringVarP(&updateFlags.issueType, "type", "t", "", "New issue type(bug, feature, task)")
	updateCmd.Flags().IntVarP(&updateFlags.priority, "priority", "p", 0, "New issue priority(0-5)")

	updateCmd.RegisterFlagCompletionFunc("type", completionFunc(typeOptions))
	updateCmd.RegisterFlagCompletionFunc("status", completionFunc(statusOptions))
	updateCmd.RegisterFlagCompletionFunc("priority", completionFunc(priorityRange))

	rootCmd.AddCommand(updateCmd)
}

func getUpdateValues(cmd *cobra.Command) (map[string]interface{}, error) {
	updates := make(map[string]interface{})

	if cmd.Flags().Changed("title") {
		if updateFlags.title == "" {
			return updates, fmt.Errorf("issue title cannot be empty")
		}
		updates["title"] = updateFlags.title
	}

	if cmd.Flags().Changed("desc") {
		updates["description"] = updateFlags.description
	}

	if cmd.Flags().Changed("status") {
		updates["status"] = updateFlags.status
	}

	if cmd.Flags().Changed("type") {
		updates["issue_type"] = updateFlags.issueType
	}

	if cmd.Flags().Changed("priority") {
		updates["priority"] = updateFlags.priority
	}

	if len(updates) == 0 {
		return updates, fmt.Errorf("no updates specified")
	}

	return updates, nil
}
