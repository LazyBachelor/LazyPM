package commands

import (
	"fmt"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/charmbracelet/huh"

	"github.com/spf13/cobra"
)

// createFlags holds the flag values for the create command
var createFlags Flags

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

	Args:    cobra.MinimumNArgs(0),
	Aliases: []string{"add"},
	RunE:    runCreateCmd,
}

// runCreateCmd executes the create command logic,
func runCreateCmd(cmd *cobra.Command, args []string) error {
	createFlags.title = strings.Join(args, " ")

	// Run interactive if flag is set
	if createFlags.interactive {
		if err := runCreateInteractive(); err != nil {
			return err
		}
	}

	if createFlags.title == "" {
		return fmt.Errorf("issue title cannot be empty")
	}

	issue := &models.Issue{
		Title:       createFlags.title,
		Description: createFlags.description,
		Status:      models.Status(createFlags.status),
		IssueType:   models.IssueType(createFlags.issueType),
		Priority:    createFlags.priority,
	}

	// Create the issue using the service layer.
	err := svc.Beads.CreateIssue(cmd.Context(), issue, "test_actor")
	if err != nil {
		return fmt.Errorf("error creating issue: %w", err)
	}

	// Display the created issue details to the user.
	cmd.Printf("Created issue:\n%s", models.IssueString(*issue))

	return nil
}

func runCreateInteractive() error {
	form := huh.NewForm(

		huh.NewGroup(
			huh.NewInput().Value(&createFlags.title).Title("Title"),
			huh.NewText().Value(&createFlags.description).Title("Description"),
		).Title("Issue Details"),

		huh.NewGroup(
			huh.NewSelect[string]().
				Options(
					huh.NewOption("Open", "open"),
					huh.NewOption("Closed", "closed"),
					huh.NewOption("In Progress", "in_progress"),
				).Value(&createFlags.status).Title("Status"),

			huh.NewSelect[string]().
				Options(
					huh.NewOption("Bug", "bug"),
					huh.NewOption("Feature", "feature"),
					huh.NewOption("Task", "task"),
				).Value(&createFlags.issueType).Title("Type"),

			huh.NewSelect[int]().
				Options(
					huh.NewOption("0", 0),
					huh.NewOption("1", 1),
					huh.NewOption("2", 2),
					huh.NewOption("3", 3),
					huh.NewOption("4", 4),
					huh.NewOption("5", 5),
				).Value(&createFlags.priority).Title("Priority"),
		).Title("Create New Issue").WithTheme(huh.ThemeBase()),
	)

	return form.Run()
}

// init function to set up the create command and its flags.
func init() {
	createCmd.Flags().BoolVarP(&createFlags.interactive, "interactive", "i", false, "Create issue interactively")
	createCmd.Flags().StringVarP(&createFlags.description, "desc", "d", "", "Issue description")
	createCmd.Flags().StringVarP(&createFlags.status, "status", "s", "open", "Issue status(open, closed, in_progress)")
	createCmd.Flags().StringVarP(&createFlags.issueType, "type", "t", "task", "Issue type(bug, feature, task)")
	createCmd.Flags().IntVarP(&createFlags.priority, "priority", "p", 0, "Issue priority(0-5)")

	createCmd.RegisterFlagCompletionFunc("type", completionFunc(typeOptions))
	createCmd.RegisterFlagCompletionFunc("status", completionFunc(statusOptions))
	createCmd.RegisterFlagCompletionFunc("priority", completionFunc(priorityRange))

	rootCmd.AddCommand(createCmd)
}
