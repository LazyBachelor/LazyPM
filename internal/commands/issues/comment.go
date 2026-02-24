package issuesCmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// commentFlags holds the flag values for the comment command.
var commentFlags struct {
	message string
	author  string
}

const commentCmdExample = `pm comment ISSUE-1 "Fix the login bug"
pm comment ISSUE-1 "LGTM" --author "alice"
pm comment ISSUE-1 -m "Needs review"`

// CommentCmd represents the comment command,
// which allows users to add a comment on an existing issue by its ID.
var CommentCmd = &cobra.Command{
	Use:     "comment [issue ID] [message]",
	Short:   "Add a comment on an issue",
	Long:    `Add a comment on an issue by ID. Message can be passed as an argument or via --message.`,
	Example: commentCmdExample,

	Args:              cobra.RangeArgs(1, 2),
	RunE:              runCommentCmd,
	ValidArgsFunction: completeIssues,
}

// runCommentCmd executes the comment command logic,
// which adds a comment on an issue by its ID.
func runCommentCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]
	var text string
	if len(args) >= 2 {
		text = args[1]
	} else {
		text = commentFlags.message
	}
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("comment text cannot be empty (use --message or pass as argument)")
	}

	app := AppFromContext(cmd.Context())

	// Ensure issue exists.
	issue, err := app.Issues.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return fmt.Errorf("error fetching issue: %w", err)
	}
	if issue == nil {
		return fmt.Errorf("issue with ID %s not found", issueID)
	}

	author := commentFlags.author
	if author == "" {
		author = "cli"
	}

	comment, err := app.Issues.AddIssueComment(cmd.Context(), issue.ID, author, text)
	if err != nil {
		return fmt.Errorf("error adding comment: %w", err)
	}

	cmd.Printf("Comment added on %s:\n", issueID)
	cmd.Printf("  %s @ %s: %s\n", comment.Author, comment.CreatedAt.Format("2006-01-02 15:04"), comment.Text)
	return nil
}

func init() {
	CommentCmd.Flags().StringVarP(&commentFlags.message, "message", "m", "", "Comment text (alternative to positional argument)")
	CommentCmd.Flags().StringVarP(&commentFlags.author, "author", "a", "cli", "Author name for the comment")
}
