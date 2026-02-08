package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	updateDescription string
	updateStatus      string
	updateType        string
	updatePriority    int
	updateTitle       string
)

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

func init() {
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "New issue title")
	updateCmd.Flags().StringVarP(&updateDescription, "desc", "d", "", "New issue description")
	updateCmd.Flags().StringVarP(&updateStatus, "status", "s", "", "New issue status(open, closed, in_progress)")
	updateCmd.Flags().StringVarP(&updateType, "type", "", "", "New issue type(bug, feature, task)")
	updateCmd.Flags().IntVarP(&updatePriority, "priority", "p", -1, "New issue priority(0-5)")

	updateCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"bug", "feature", "task"}, cobra.ShellCompDirectiveDefault
	})

	updateCmd.RegisterFlagCompletionFunc("status", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"open", "closed", "in_progress"}, cobra.ShellCompDirectiveDefault
	})

	updateCmd.RegisterFlagCompletionFunc("priority", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"0", "1", "2", "3", "4", "5"}, cobra.ShellCompDirectiveDefault
	})

	rootCmd.AddCommand(updateCmd)
}

func runUpdateCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	updates := make(map[string]interface{})

	if cmd.Flags().Changed("title") {
		if updateTitle == "" {
			return fmt.Errorf("issue title cannot be empty")
		}
		updates["title"] = updateTitle
	}

	if cmd.Flags().Changed("desc") {
		updates["description"] = updateDescription
	}

	if cmd.Flags().Changed("status") {
		updates["status"] = updateStatus
	}

	if cmd.Flags().Changed("type") {
		updates["issue_type"] = updateType
	}

	if cmd.Flags().Changed("priority") {
		updates["priority"] = updatePriority
	}

	if len(updates) == 0 {
		return fmt.Errorf("no updates specified")
	}

	issue, err := svc.Beads.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return fmt.Errorf("error getting issue: %w", err)
	}

	if issue == nil {
		return fmt.Errorf("issue with ID '%s' not found", issueID)
	}

	err = svc.Beads.UpdateIssue(cmd.Context(), issueID, updates, "test_actor")
	if err != nil {
		return fmt.Errorf("error updating issue: %w", err)
	}

	updatedIssue, err := svc.Beads.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return fmt.Errorf("error getting updated issue: %w", err)
	}

	str := fmt.Sprintf("Updated issue with ID: %s\n", issueID)

	if updatedIssue.Title != "" {
		str += fmt.Sprintf("Title: %s\n", updatedIssue.Title)
	}

	if updatedIssue.Description != "" {
		str += fmt.Sprintf("Description: %s\n", updatedIssue.Description)
	}

	if updatedIssue.Status != "" {
		str += fmt.Sprintf("Status: %s\n", updatedIssue.Status)
	}

	if updatedIssue.IssueType != "" {
		str += fmt.Sprintf("Type: %s\n", updatedIssue.IssueType)
	}

	if cmd.Flags().Changed("priority") {
		str += fmt.Sprintf("Priority: %d\n", updatedIssue.Priority)
	}

	fmt.Print(str)

	return nil
}
