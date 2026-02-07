package commands

import (
	"fmt"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"

	"github.com/spf13/cobra"
)

// Variables to hold flag values for the create command.
var (
	createDescription string
	createStatus      string
	createType        string
	createPriority    int
)

const (
	createCmdExample = `pm create New issue -d "Description" -s open -t task -p 3
pm create Fix bug --desc "Bug description" --status in_progress --type bug --priority 5`
)

// createCmd represents the create command, which allows users to create a new issue with specified details.
var createCmd = &cobra.Command{
	Use:     "create [title]",
	Short:   "Create a new issue",
	Long:    `Create a new issue with the specified details.`,
	Example: createCmdExample,

	Aliases: []string{"add"},
	Args:    cobra.MinimumNArgs(1),
	RunE:    runCreateCmd,
}

// runCreateCmd executes the create command logic,
func runCreateCmd(cmd *cobra.Command, args []string) error {
	createTitle := strings.Join(args, " ")

	if createTitle == "" {
		return fmt.Errorf("issue title cannot be empty")
	}

	issue := &models.Issue{
		Title:       createTitle,
		Description: createDescription,
		Status:      models.Status(createStatus),
		IssueType:   models.IssueType(createType),
		Priority:    createPriority,
	}

	// Create the issue using the service layer.
	err := svc.Beads.CreateIssue(cmd.Context(), issue, "test_actor")
	if err != nil {
		return fmt.Errorf("error creating issue: %w", err)
	}

	// Build the output string with the created issue details.
	str := fmt.Sprintf("Created issue with ID: %s\n", issue.ID)

	if issue.Title != "" {
		str += fmt.Sprintf("Title: %s\n", issue.Title)
	}

	if issue.Description != "" {
		str += fmt.Sprintf("Description: %s\n", issue.Description)
	}

	if issue.Status != "" {
		str += fmt.Sprintf("Status: %s\n", issue.Status)
	}

	if issue.IssueType != "" {
		str += fmt.Sprintf("Type: %s\n", issue.IssueType)
	}

	if issue.Priority != 0 {
		str += fmt.Sprintf("Priority: %d\n", issue.Priority)
	}

	cmd.Print(str)

	return nil
}

// init function to set up the create command and its flags.
func init() {
	createCmd.Flags().StringVarP(&createDescription, "desc", "d", "", "Issue description")
	createCmd.Flags().StringVarP(&createStatus, "status", "s", "open", "Issue status(open, closed, in_progress)")
	createCmd.Flags().StringVarP(&createType, "type", "t", "task", "Issue type(bug, feature, task)")
	createCmd.Flags().IntVarP(&createPriority, "priority", "p", 0, "Issue priority(0-5)")

	createCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"bug", "feature", "task"}, cobra.ShellCompDirectiveDefault
	})

	createCmd.RegisterFlagCompletionFunc("status", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"open", "closed", "in_progress"}, cobra.ShellCompDirectiveDefault
	})

	createCmd.RegisterFlagCompletionFunc("priority", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"0", "1", "2", "3", "4", "5"}, cobra.ShellCompDirectiveDefault
	})

	rootCmd.AddCommand(createCmd)
}
