package issues

import (
	"fmt"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/style"
	"github.com/spf13/cobra"
)

var CloseCmd = &cobra.Command{
	Use:     "close [id]",
	Short:   "Close an existing issue",
	Long:    `Close an existing issue by its ID.`,
	Example: `pm close pm-abc`,

	Args: cobra.ExactArgs(1),
	RunE: runCloseCmd,

	ValidArgsFunction: completeIssues,
}

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

	closeReason := ""
	if err = huh.NewSelect[string]().Value(&closeReason).
		Title("Reason for closing the issue?").
		Options(
			huh.NewOption("Done", "Done"),
			huh.NewOption("Duplicate issue", "Duplicate issue"),
			huh.NewOption("Won't fix", "Won't fix"),
			huh.NewOption("Obsolete", "Obsolete"),
			huh.NewOption("Other", "Other"),
		).WithTheme(style.BaseTheme{}).Run(); err != nil {
		return fmt.Errorf("error getting close reason: %w", err)
	}

	if closeReason == "Other" {
		if err = huh.NewInput().Value(&closeReason).
			Title("Enter closing reason:").WithTheme(style.BaseTheme{}).Run(); err != nil {
			return fmt.Errorf("error getting close reason: %w", err)
		}
		if closeReason == "" {
			return fmt.Errorf("closing reason cannot be empty when selecting 'Other'")
		}
	}

	err = app.Issues.CloseIssue(cmd.Context(), closeID, closeReason, "", "")
	if err != nil {
		return fmt.Errorf("error closing issue: %w", err)
	}

	cmd.Println("Closed issue with ID:", closeID)

	return nil
}
