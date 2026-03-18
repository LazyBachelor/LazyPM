package issues

import (
	"context"
	"fmt"
	"strings"

	"charm.land/huh/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/style"
	"github.com/spf13/cobra"
)

// Variables for delete command flag.
var (
	confirmDelete     bool
	deleteIDs         []string
	deleteInteractive bool
)

// DeleteCmd represents the delete command.
var DeleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Short:   "Delete an existing issue",
	Long:    `Delete an existing issue by its ID.`,
	Example: `pm delete pm-abc`,

	ValidArgsFunction: completeIssues,

	Aliases: []string{"del", "remove", "rm"},
	RunE:    runDeleteCmd,
}

// runDeleteCmd executes the delete command logic,
// which deletes an issue by its ID after confirming with the user.
func runDeleteCmd(cmd *cobra.Command, args []string) error {
	deleteID := strings.Join(args, " ")

	if deleteInteractive {
		if err := runDeleteInteractive(cmd.Context()); err != nil {
			return err
		}
		return nil
	}

	if deleteID == "" {
		return fmt.Errorf("issue ID cannot be empty")
	}

	app := AppFromContext(cmd.Context())

	// Fetch the issue to ensure it exists before deletion.
	issue, err := app.Issues.GetIssue(cmd.Context(), deleteID)
	if err != nil {
		return fmt.Errorf("error fetching issue: %w", err)
	}

	if issue == nil {
		return fmt.Errorf("issue with ID %s not found", deleteID)
	}

	// Prompt for confirmation if not already confirmed via flag.
	if !cmd.Flags().Changed("yes") {
		huh.NewConfirm().Value(&confirmDelete).
			Title("You want to delete this issue?").
			Inline(true).WithTheme(style.BaseTheme{}).Run()
	}

	// If user did not confirm, cancel deletion.
	if !confirmDelete {
		cmd.Println("Deletion cancelled.")
		return nil
	}

	// Delete the issue.
	err = app.Issues.DeleteIssue(cmd.Context(), deleteID)
	if err != nil {
		return fmt.Errorf("error deleting issue: %w", err)
	}

	cmd.Println("Deleted issue with ID:", deleteID)

	return nil
}

// runDeleteInteractive runs the interactive mode for deleting issues,
// allowing users to select multiple issues for deletion.
func runDeleteInteractive(ctx context.Context) error {
	app := AppFromContext(ctx)
	options := []huh.Option[string]{}

	issues, err := app.Issues.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return fmt.Errorf("error fetching issues: %w", err)
	}

	for _, issue := range issues {
		desc := fmt.Sprintf("%s: %s", issue.ID, issue.Title)
		options = append(options, huh.NewOption(desc, issue.ID))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(options...).Value(&deleteIDs).
				Title("Select issues to delete"))).WithTheme(style.BaseTheme{})

	if err := form.Run(); err != nil {
		return fmt.Errorf("error running interactive form: %w", err)
	}

	if len(deleteIDs) == 0 {
		return fmt.Errorf("no issues selected for deletion")
	}

	for _, id := range deleteIDs {
		err := app.Issues.DeleteIssue(ctx, id)
		if err != nil {
			return fmt.Errorf("error deleting issue with ID %s: %w", id, err)
		}
		fmt.Printf("Deleted issue with ID: %s\n", id)
	}

	return nil
}

// init function to set up the delete command and its flags.
func init() {
	DeleteCmd.Flags().BoolVarP(&deleteInteractive, "interactive", "i", false, "Delete issues interactively")
	DeleteCmd.Flags().BoolVarP(&confirmDelete, "yes", "y", true, "Confirm deletion without prompt")
}
