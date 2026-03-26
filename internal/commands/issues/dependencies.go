package issues

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
)

const depCmdExample = `pm dep ISSUE-1
pm dep view ISSUE-1
pm dep add ISSUE-1 ISSUE-2
pm dep remove ISSUE-1 ISSUE-2`

// DepCmd manages dependencies for issues.
// Default behavior: view dependencies of an issue.
var DepCmd = &cobra.Command{
	Use:     "dep [issue id]",
	Short:   "Manage dependencies for an issue",
	Long:    `Manage dependencies between issues.`,
	Aliases: []string{"dependencies"},
	Example: depCmdExample,

	// This command runs only when no subcommand is provided.
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeIssues,
	RunE:              runDepViewCmd,
}

var DepViewCmd = &cobra.Command{
	Use:               "view [issue id]",
	Short:             "View dependencies of an issue",
	Aliases:           []string{"show"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeIssues,
	RunE:              runDepViewCmd,
}

var DepAddCmd = &cobra.Command{
	Use:               "add [issue id] [depends-on id]",
	Short:             "Add a dependency: issue depends on depends-on issue",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeIssues,
	RunE:              runDepAddCmd,
}

var DepRemoveCmd = &cobra.Command{
	Use:               "remove [issue id] [depends-on id]",
	Short:             "Remove an existing dependency",
	Aliases:           []string{"rm"},
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeIssues,
	RunE:              runDepRemoveCmd,
}

func init() {
	DepCmd.AddCommand(DepViewCmd)
	DepCmd.AddCommand(DepAddCmd)
	DepCmd.AddCommand(DepRemoveCmd)
}

func runDepViewCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]
	app := AppFromContext(cmd.Context())

	deps, err := app.Issues.GetDependencies(cmd.Context(), issueID)
	if err != nil {
		return fmt.Errorf("failed to load dependencies for %s: %w", issueID, err)
	}

	if len(deps) == 0 {
		cmd.Printf("Issue %s has no dependencies.\n", issueID)
		return nil
	}

	cmd.Printf("Dependencies for %s:\n", issueID)
	for _, d := range deps {
		if d == nil {
			continue
		}
		if d.Title != "" {
			cmd.Printf("- %s: %s\n", d.ID, d.Title)
		} else {
			cmd.Printf("- %s\n", d.ID)
		}
	}

	return nil
}

func runDepAddCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]
	dependsOnID := args[1]

	app := AppFromContext(cmd.Context())

	dep := &models.Dependency{
		IssueID:     issueID,
		DependsOnID: dependsOnID,
		Type:        models.DepBlocks,
	}

	if err := app.Issues.AddDependency(cmd.Context(), dep, "cli"); err != nil {
		return fmt.Errorf("failed to add dependency: %s depends on %s: %w", issueID, dependsOnID, err)
	}

	cmd.Printf("Added dependency: %s depends on %s\n", issueID, dependsOnID)
	return runDepViewCmd(cmd, []string{issueID})
}

func runDepRemoveCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]
	dependsOnID := args[1]

	app := AppFromContext(cmd.Context())

	if err := app.Issues.RemoveDependency(cmd.Context(), issueID, dependsOnID, "cli"); err != nil {
		return fmt.Errorf("failed to remove dependency: %s depends on %s: %w", issueID, dependsOnID, err)
	}

	cmd.Printf("Removed dependency: %s no longer depends on %s\n", issueID, dependsOnID)
	return runDepViewCmd(cmd, []string{issueID})
}
