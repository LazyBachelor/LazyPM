package issuesCmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CommentsCmd represents the command to list comments on an issue.
var CommentsCmd = &cobra.Command{
	Use:   "comments [issue ID]",
	Short: "List comments on an issue",
	Long:  `List all comments on an issue by ID.`,

	Args:              cobra.ExactArgs(1),
	RunE:              runCommentsCmd,
	ValidArgsFunction: completeIssues,
}

// runCommentsCmd executes the comments command logic.
func runCommentsCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	app := AppFromContext(cmd.Context())

	issue, err := app.Issues.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return fmt.Errorf("error fetching issue: %w", err)
	}
	if issue == nil {
		return fmt.Errorf("issue with ID %s not found", issueID)
	}

	comments, err := app.Issues.GetIssueComments(cmd.Context(), issue.ID)
	if err != nil {
		return fmt.Errorf("error fetching comments: %w", err)
	}

	cmd.Printf("Comments on %s (%s):\n\n", issueID, issue.Title)
	if len(comments) == 0 {
		cmd.Println("  No comments yet.")
		return nil
	}
	for _, c := range comments {
		if c != nil {
			cmd.Printf("  %s @ %s:\n    %s\n\n", c.Author, c.CreatedAt.Format("2006-01-02 15:04"), c.Text)
		}
	}
	return nil
}
