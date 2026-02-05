package commands

import (
	"fmt"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Short:   "Delete existing issue",
	Long:    `Delete an existing issue by its ID.`,
	Example: `pm delete issue_id`,
	RunE:    runDeleteCmd,
	Aliases: []string{"del"},
	Args:    cobra.MinimumNArgs(1),
}

func runDeleteCmd(cmd *cobra.Command, args []string) error {
	deleteID := strings.Join(args, " ")

	issue := &models.Issue{
		ID: deleteID,
	}

	err := svc.Beads.DeleteIssue(cmd.Context(), issue.ID)
	if err != nil {
		return fmt.Errorf("error deleting issue: %w", err)
	}

	str := fmt.Sprintf("Deleted issue with ID: %s\n", issue.ID)

	fmt.Print(str)

	return nil
}
