package commands

import (
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
)

// Variables for get-issues command flags.
var (
	titleFlag       string
	descriptionFlag string
	statusFlag      string
	typeFlag        string
	priorityFlag    int
	limit           int = 25
)

const (
	lsExamples = `pm ls [id|title|description]
pm ls --status open --type bug
pm ls --title "New feature" --desc "feature description"
pm ls -p 1 -l 10`
)

// getIssuesCmd represents the get issues command.
var getIssuesCmd = &cobra.Command{
	Use:     "ls [search query]",
	Short:   "List all issues",
	Long:    `List all issues in the project management system.`,
	Example: lsExamples,

	Aliases: []string{"list", "search"},
	Args:    cobra.MinimumNArgs(0),
	RunE:    runGetIssuesCmd,
}

// runGetIssuesCmd executes the get issues command logic,
// which retrieves and displays a list of issues based on the provided search query and filters.
func runGetIssuesCmd(cmd *cobra.Command, args []string) error {
	queryArg := strings.Join(args, " ")

	filter := models.IssueFilter{
		TitleSearch:         titleFlag,
		DescriptionContains: descriptionFlag,
		Limit:               limit,
	}

	// Only set filter fields if the corresponding flags
	// were explicitly provided by the user.
	if cmd.Flags().Changed("status") {
		s := models.Status(statusFlag)
		filter.Status = &s
	}
	if cmd.Flags().Changed("type") {
		t := models.IssueType(typeFlag)
		filter.IssueType = &t
	}
	if cmd.Flags().Changed("priority") {
		filter.Priority = &priorityFlag
	}

	// Fetch issues based on the search query and filters.
	issuesPtr, err := svc.Beads.SearchIssues(cmd.Context(), queryArg, filter)
	if err != nil {
		return err
	}

	// Convert the returned issue pointers to issue values and print them.
	issues := models.IssuesPtrToIssues(issuesPtr)
	models.PrintIssues(issues)

	return nil
}

// init function to set up the get issues command and its flags.
func init() {
	getIssuesCmd.Flags().StringVar(&titleFlag, "title", "", "Filter issues by title")
	getIssuesCmd.Flags().StringVarP(&descriptionFlag, "desc", "d", "", "Filter issues by description")
	getIssuesCmd.Flags().StringVarP(&statusFlag, "status", "s", "", "Filter issues by status (open, closed, in_progress)")
	getIssuesCmd.Flags().StringVarP(&typeFlag, "type", "t", "", "Filter issues by type (bug, feature, task)")
	getIssuesCmd.Flags().IntVarP(&priorityFlag, "priority", "p", 0, "Filter issues by priority (0-5)")
	getIssuesCmd.Flags().IntVarP(&limit, "limit", "l", 25, "Limit the number of issues returned")

	getIssuesCmd.RegisterFlagCompletionFunc("status", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"open", "closed", "in_progress"}, cobra.ShellCompDirectiveDefault
	})

	getIssuesCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"bug", "feature", "task"}, cobra.ShellCompDirectiveDefault
	})

	getIssuesCmd.RegisterFlagCompletionFunc("priority", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"0", "1", "2", "3", "4", "5"}, cobra.ShellCompDirectiveDefault
	})

	rootCmd.AddCommand(getIssuesCmd)
}
