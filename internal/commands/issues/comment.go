package issues

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
)

// commentFlags holds the flag values for the comment command.
var commentFlags struct {
	message string
	author  string
}

const commentCmdExample = `pm comment ISSUE-1 Fix the login bug
pm comment ISSUE-1 LGTM --author alice
pm comment ISSUE-1 -m "Needs review"`

var CommentCmd = &cobra.Command{
	Use:     "comment [id] [message...]",
	Short:   "Add a comment on an issue",
	Long:    `Add a comment on an issue by ID. All arguments after the issue ID form the message (no quotes needed), or use --message.`,
	Example: commentCmdExample,

	Args:              cobra.RangeArgs(1, 100), // issue ID + up to 99 words for the message
	RunE:              runCommentCmd,
	ValidArgsFunction: completeIssues,
}

// runCommentCmd executes the comment command logic,
// which adds a comment on an issue by its ID.
func runCommentCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]
	var text string
	if len(args) > 1 {
		text = strings.Join(args[1:], " ")
	} else {
		text = commentFlags.message
	}
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("comment text cannot be empty (use --message or pass words after the issue ID)")
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
		author = defaultCommentAuthor()
	}

	comment, err := app.Issues.AddIssueComment(cmd.Context(), issue.ID, author, text)
	if err != nil {
		return fmt.Errorf("error adding comment: %w", err)
	}

	cmd.Printf("Comment added on %s:\n", issueID)
	cmd.Printf("  %s @ %s: %s\n", comment.Author, comment.CreatedAt.Format("2006-01-02 15:04"), comment.Text)
	return nil
}

func defaultCommentAuthor() string {
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	if s := os.Getenv("USER"); s != "" {
		return s
	}
	if s := os.Getenv("USERNAME"); s != "" {
		return s
	}
	return "user"
}

func init() {
	CommentCmd.Flags().StringVarP(&commentFlags.message, "message", "m", "", "Comment text (alternative to positional arguments)")
	CommentCmd.Flags().StringVarP(&commentFlags.author, "author", "a", "", "Author name (default: current OS user)")
}
