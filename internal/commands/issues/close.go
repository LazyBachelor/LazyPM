package issues

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// CloseCmd represents the close command,
// which allows users to close an existing issue by its ID.
var CloseCmd = &cobra.Command{
	Use:     "close [id]",
	Short:   "Close an existing issue",
	Long:    `Close an existing issue by its ID.`,
	Example: `pm close pm-abc`,
	

	Args: cobra.ExactArgs(1),
	RunE: runCloseCmd,

	ValidArgsFunction: completeIssues,
}

// runCloseCmd executes the close command logic,
// which closes an issue by its ID after confirming with the user.
func runCloseCmd(cmd *cobra.Command, args []string) error {
	closeID := args[0]

	if closeID == "" {
		return fmt.Errorf("issue ID cannot be empty")
	}

	app := AppFromContext(cmd.Context())

	// Fetch the issue to ensure it exists before closing.
	issue, err := app.Issues.GetIssue(cmd.Context(), closeID)
	if err != nil {
		return fmt.Errorf("error fetching issue: %w", err)
	}

	if issue == nil {
		return fmt.Errorf("issue with ID %s not found", closeID)
	}

	// Ask for closing reason
	if err = huh.NewInput().Value(&issue.CloseReason).
		Title("Reason for closing the issue?").WithTheme(huh.ThemeBase()).Run(); err != nil {
		return fmt.Errorf("error getting close reason: %w", err)
	}

	// Close the issue.
	err = app.Issues.CloseIssue(cmd.Context(), closeID, issue.CloseReason, "", "")
	if err != nil {
		return fmt.Errorf("error closing issue: %w", err)
	}

	cmd.Println("Closed issue with ID:", closeID)

	return nil
}
