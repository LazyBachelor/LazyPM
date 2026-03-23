package issues

import (
	"fmt"
	"strconv"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
)

// sprintFlags holds the flag values for sprint commands
type sprintFlags struct {
	sprintNum int
	issueID   string
}

var sprintCmdFlags sprintFlags

// SprintCmd represents the sprint management command
var SprintCmd = &cobra.Command{
	Use:   "sprint",
	Short: "Manage sprints",
	Long:  `Manage sprints - create, list, add/remove issues, and view sprint contents.`,
}

// SprintListCmd lists all sprints
var SprintListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all sprints",
	Long:    `List all sprints in the project.`,
	Aliases: []string{"ls"},
	RunE:    runSprintListCmd,
}

// SprintCreateCmd creates a new sprint
var SprintCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new sprint",
	Long:    `Create a new sprint for organizing issues.`,
	Aliases: []string{"new"},
	RunE:    runSprintCreateCmd,
}

// SprintIssuesCmd lists issues in a sprint
var SprintIssuesCmd = &cobra.Command{
	Use:     "issues [sprint-num]",
	Short:   "List issues in a sprint",
	Long:    `List all issues assigned to a specific sprint. Use 'backlog' or omit to view the backlog.`,
	Aliases: []string{"show", "view"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runSprintIssuesCmd,
}

// SprintAddCmd adds an issue to a sprint
var SprintAddCmd = &cobra.Command{
	Use:   "add <issue-id> [sprint-num]",
	Short: "Add an issue to a sprint",
	Long:  `Add an issue to a sprint. If sprint-num is omitted, adds to backlog.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runSprintAddCmd,

	ValidArgsFunction: completeIssues,
}

// SprintRemoveCmd removes an issue from a sprint
var SprintRemoveCmd = &cobra.Command{
	Use:     "remove <issue-id> [sprint-num]",
	Short:   "Remove an issue from a sprint",
	Long:    `Remove an issue from a sprint. If sprint-num is omitted, removes from backlog.`,
	Aliases: []string{"rm"},
	Args:    cobra.RangeArgs(1, 2),
	RunE:    runSprintRemoveCmd,

	ValidArgsFunction: completeIssues,
}

// SprintBacklogCmd shows the backlog sprint
var SprintBacklogCmd = &cobra.Command{
	Use:   "backlog",
	Short: "Show backlog issues",
	Long:  `Show all issues in the backlog sprint.`,
	RunE:  runSprintBacklogCmd,
}

// SprintDeleteCmd deletes a sprint
var SprintDeleteCmd = &cobra.Command{
	Use:     "delete <sprint-num>",
	Short:   "Delete a sprint",
	Long:    `Delete a sprint by its number. Issues in the sprint will not be deleted.`,
	Aliases: []string{"del", "rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runSprintDeleteCmd,
}

func runSprintListCmd(cmd *cobra.Command, args []string) error {
	app := AppFromContext(cmd.Context())

	sprints, err := app.Issues.GetSprints(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get sprints: %w", err)
	}

	if len(sprints) == 0 {
		cmd.Println("No sprints found.")
		return nil
	}

	backlogNum, _ := app.Issues.GetBacklogSprint(cmd.Context())

	cmd.Println("Sprints:")
	for _, sprintNum := range sprints {
		label := ""
		if sprintNum == backlogNum {
			label = " (backlog)"
		}
		cmd.Printf("  Sprint %d%s\n", sprintNum, label)
	}

	return nil
}

func runSprintCreateCmd(cmd *cobra.Command, args []string) error {
	app := AppFromContext(cmd.Context())

	sprintNum, err := app.Issues.AddSprint(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to create sprint: %w", err)
	}

	cmd.Printf("Created sprint %d\n", sprintNum)
	return nil
}

func runSprintIssuesCmd(cmd *cobra.Command, args []string) error {
	app := AppFromContext(cmd.Context())
	ctx := cmd.Context()

	var sprintNum int
	var isBacklog bool
	var err error

	if len(args) == 0 || args[0] == "backlog" {
		sprintNum, err = app.Issues.GetBacklogSprint(ctx)
		if err != nil {
			return fmt.Errorf("failed to get backlog sprint: %w", err)
		}
		isBacklog = true
	} else {
		sprintNum, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid sprint number: %s", args[0])
		}
		isBacklog = false
	}

	var issues []*models.Issue
	if isBacklog {
		issues, err = app.Issues.GetIssuesNotInAnySprint(ctx)
	} else {
		issues, err = app.Issues.GetIssuesBySprint(ctx, sprintNum)
	}
	if err != nil {
		return fmt.Errorf("failed to get issues in sprint %d: %w", sprintNum, err)
	}

	if isBacklog {
		cmd.Printf("Backlog (%d issues):\n", len(issues))
	} else {
		cmd.Printf("Sprint %d (%d issues):\n", sprintNum, len(issues))
	}

	if len(issues) == 0 {
		cmd.Println("  No issues in this sprint.")
		return nil
	}

	issuesList := models.IssuesPtrToIssues(issues)
	models.PrintIssues(issuesList)

	return nil
}

func runSprintAddCmd(cmd *cobra.Command, args []string) error {
	app := AppFromContext(cmd.Context())
	ctx := cmd.Context()

	issueID := args[0]

	var sprintNum int
	var err error

	if len(args) == 1 {
		sprintNum, err = app.Issues.GetBacklogSprint(ctx)
		if err != nil {
			return fmt.Errorf("failed to get backlog sprint: %w", err)
		}
	} else {
		sprintNum, err = strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid sprint number: %s", args[1])
		}
	}

	err = app.Issues.AddIssueToSprint(ctx, issueID, sprintNum)
	if err != nil {
		return fmt.Errorf("failed to add issue to sprint: %w", err)
	}

	backlogNum, _ := app.Issues.GetBacklogSprint(ctx)

	if sprintNum == backlogNum {
		cmd.Printf("Added issue %s to backlog\n", issueID)
	} else {
		cmd.Printf("Added issue %s to sprint %d\n", issueID, sprintNum)
	}

	return nil
}

func runSprintRemoveCmd(cmd *cobra.Command, args []string) error {
	app := AppFromContext(cmd.Context())
	ctx := cmd.Context()

	issueID := args[0]

	var sprintNum int
	var err error

	if len(args) == 1 {
		sprintNum, err = app.Issues.GetBacklogSprint(ctx)
		if err != nil {
			return fmt.Errorf("failed to get backlog sprint: %w", err)
		}
	} else {
		sprintNum, err = strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid sprint number: %s", args[1])
		}
	}

	err = app.Issues.RemoveIssueFromSprint(ctx, issueID, sprintNum)
	if err != nil {
		return fmt.Errorf("failed to remove issue from sprint: %w", err)
	}

	backlogNum, _ := app.Issues.GetBacklogSprint(ctx)

	if sprintNum == backlogNum {
		cmd.Printf("Removed issue %s from backlog\n", issueID)
	} else {
		cmd.Printf("Removed issue %s from sprint %d\n", issueID, sprintNum)
	}

	return nil
}

func runSprintBacklogCmd(cmd *cobra.Command, args []string) error {
	app := AppFromContext(cmd.Context())
	ctx := cmd.Context()

	issues, err := app.Issues.GetIssuesNotInAnySprint(ctx)
	if err != nil {
		return fmt.Errorf("failed to get backlog issues: %w", err)
	}

	cmd.Printf("Backlog (%d issues):\n", len(issues))

	if len(issues) == 0 {
		cmd.Println("  No issues in backlog.")
		return nil
	}

	issuesList := models.IssuesPtrToIssues(issues)
	models.PrintIssues(issuesList)

	return nil
}

func runSprintDeleteCmd(cmd *cobra.Command, args []string) error {
	app := AppFromContext(cmd.Context())
	ctx := cmd.Context()

	sprintNum, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid sprint number: %s", args[0])
	}

	backlogNum, _ := app.Issues.GetBacklogSprint(ctx)
	if sprintNum == backlogNum {
		return fmt.Errorf("cannot delete the backlog sprint")
	}

	err = app.Issues.RemoveSprint(ctx, sprintNum)
	if err != nil {
		return fmt.Errorf("failed to delete sprint: %w", err)
	}

	cmd.Printf("Deleted sprint %d\n", sprintNum)
	return nil
}

func init() {
	SprintCmd.AddCommand(SprintListCmd)
	SprintCmd.AddCommand(SprintCreateCmd)
	SprintCmd.AddCommand(SprintIssuesCmd)
	SprintCmd.AddCommand(SprintAddCmd)
	SprintCmd.AddCommand(SprintRemoveCmd)
	SprintCmd.AddCommand(SprintBacklogCmd)
	SprintCmd.AddCommand(SprintDeleteCmd)

	RootCmd.AddCommand(SprintCmd)
}
