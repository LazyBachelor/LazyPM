package issues

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/shellcomp"
	"github.com/spf13/cobra"
)

var updateFlags Flags

var UpdateCmd = &cobra.Command{
	Use:               "update [id]",
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

	app := AppFromContext(cmd.Context())

	issue, err := app.Issues.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return fmt.Errorf("error getting issue: %w", err)
	}

	if issue == nil {
		return fmt.Errorf("issue with ID %s not found", issueID)
	}

	updates, err := getUpdateValues(cmd)
	if err != nil {
		return fmt.Errorf("error getting update values: %w", err)
	}

	err = app.Issues.UpdateIssue(cmd.Context(), issueID, updates, "test_actor")
	if err != nil {
		return fmt.Errorf("error updating issue: %w", err)
	}

	updatedIssue, err := app.Issues.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return fmt.Errorf("error getting updated issue: %w", err)
	}

	cmd.Printf("Updated issue to:\n%s", models.IssueString(*updatedIssue))

	return nil
}

func init() {
	UpdateCmd.Flags().StringVar(&updateFlags.title, "title", "", "New issue title")
	UpdateCmd.Flags().StringVarP(&updateFlags.description, "desc", "d", "", "New issue description")
	UpdateCmd.Flags().StringVarP(&updateFlags.status, "status", "s", "", "New issue status(open, closed, in_progress)")
	UpdateCmd.Flags().StringVarP(&updateFlags.issueType, "type", "t", "", "New issue type(bug, feature, task)")
	UpdateCmd.Flags().IntVarP(&updateFlags.priority, "priority", "p", 0, "New issue priority(0-4)")
	UpdateCmd.Flags().StringVarP(&updateFlags.assignee, "assignee", "a", "", "New issue assignee")

	UpdateCmd.RegisterFlagCompletionFunc("type", shellcomp.CompletionFunc(typeOptions))
	UpdateCmd.RegisterFlagCompletionFunc("status", shellcomp.CompletionFunc(statusOptions))
	UpdateCmd.RegisterFlagCompletionFunc("priority", shellcomp.CompletionFunc(priorityRange))
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

	if cmd.Flags().Changed("assignee") {
		updates["assignee"] = updateFlags.assignee
	}

	if len(updates) == 0 {
		return updates, fmt.Errorf("no updates specified")
	}

	return updates, nil
}
