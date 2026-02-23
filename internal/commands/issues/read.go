package issuesCmd

import (
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
)

// GetCmd represents the get issue command.
var GetCmd = &cobra.Command{
	Use:   "describe [issue ID]",
	Short: "Get issue details",
	Long:  `Get issue details by ID`,

	ValidArgsFunction: completeIssues,

	Aliases: []string{"get", "read"},
	Args:    cobra.ExactArgs(1),
	RunE:    runGetCmd,
}

// runGetCmd executes the get issue command logic,
// which retrieves and displays issue details by its ID.
func runGetCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	// Fetch the issue details using the service layer.
	app := AppFromContext(cmd.Context())
	issuePtr, err := app.Issues.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return err
	}

	// If the issue is not found, inform the user.
	if issuePtr == nil {
		cmd.Printf("Issue with ID %s not found\n", issueID)
		return nil
	}

	// Display the issue details to the user.
	cmd.Println(models.IssueString(*issuePtr))

	return nil
}

// init function to set up the get issue command.
func init() {
}
