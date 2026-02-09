package commands

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// Variables for delete command flag.
var confirmDelete bool

// deleteCmd represents the delete command.
var deleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Short:   "Delete an existing issue",
	Long:    `Delete an existing issue by its ID.`,
	Example: `pm delete pm-abc`,

	Aliases: []string{"del", "remove", "rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runDeleteCmd,
}

// runDeleteCmd executes the delete command logic,
// which deletes an issue by its ID after confirming with the user.
func runDeleteCmd(cmd *cobra.Command, args []string) error {
	deleteID := strings.Join(args, " ")

	// Fetch the issue to ensure it exists before deletion.
	issue, err := svc.Beads.GetIssue(cmd.Context(), deleteID)
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
			Inline(true).WithTheme(huh.ThemeBase()).Run()
	}

	// If user did not confirm, cancel deletion.
	if !confirmDelete {
		cmd.Println("Deletion cancelled.")
		return nil
	}

	// Delete the issue.
	err = svc.Beads.DeleteIssue(cmd.Context(), deleteID)
	if err != nil {
		return fmt.Errorf("error deleting issue: %w", err)
	}

	cmd.Println("Deleted issue with ID:", deleteID)

	return nil
}

// init function to set up the delete command and its flags.
func init() {
	deleteCmd.Flags().BoolVarP(&confirmDelete, "yes", "y", true, "Confirm deletion without prompt")

	rootCmd.AddCommand(deleteCmd)
}
