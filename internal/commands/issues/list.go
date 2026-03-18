package issues

import (
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/shellcomp"
	"github.com/spf13/cobra"
)

// Variables for get-issues command flags.
var listFlags Flags

const (
	lsExamples = `pm list [id|title|description]
pm list --status open --type bug
pm list --title "New feature" --desc "feature description"
pm list -p 1 -l 10`
)

// ListCmd represents the get issues command.
var ListCmd = &cobra.Command{
	Use:     "list [query]",
	Short:   "List all issues",
	Long:    `List all issues in the project management system.`,
	Example: lsExamples,

	Aliases: []string{"ls", "search"},
	Args:    cobra.MinimumNArgs(0),
	RunE:    runGetIssuesCmd,
}

// runGetIssuesCmd executes the get issues command logic,
// which retrieves and displays a list of issues based on the provided search query and filters.
func runGetIssuesCmd(cmd *cobra.Command, args []string) error {
	queryArg := strings.Join(args, " ")

	filter := models.IssueFilter{
		TitleSearch:         listFlags.title,
		DescriptionContains: listFlags.description,
		Limit:               listFlags.limit,
	}

	// Only set filter fields if the corresponding flags
	// were explicitly provided by the user.
	if cmd.Flags().Changed("status") {
		filter.Status = new(models.Status(listFlags.status))
	}
	if cmd.Flags().Changed("type") {
		filter.IssueType = new(models.IssueType(listFlags.issueType))
	}
	if cmd.Flags().Changed("priority") {
		filter.Priority = &listFlags.priority
	}
	if cmd.Flags().Changed("assignee") {
		filter.Assignee = &listFlags.assignee
	}

	// Fetch issues based on the search query and filters.
	app := AppFromContext(cmd.Context())
	issuesPtr, err := app.Issues.SearchIssues(cmd.Context(), queryArg, filter)
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
	ListCmd.Flags().StringVar(&listFlags.title, "title", "", "Filter issues by title")
	ListCmd.Flags().StringVarP(&listFlags.description, "desc", "d", "", "Filter issues by description")
	ListCmd.Flags().StringVarP(&listFlags.status, "status", "s", "", "Filter issues by status (open, closed, in_progress, ready_to_sprint)")
	ListCmd.Flags().StringVarP(&listFlags.issueType, "type", "t", "", "Filter issues by type (bug, feature, task)")
	ListCmd.Flags().IntVarP(&listFlags.priority, "priority", "p", 0, "Filter issues by priority (0-4)")
	ListCmd.Flags().StringVarP(&listFlags.assignee, "assignee", "a", "", "Filter issues by assignee")

	ListCmd.Flags().IntVarP(&listFlags.limit, "limit", "l", 25, "Limit the number of issues returned")

	ListCmd.RegisterFlagCompletionFunc("status", shellcomp.CompletionFunc(statusOptions))
	ListCmd.RegisterFlagCompletionFunc("type", shellcomp.CompletionFunc(typeOptions))
	ListCmd.RegisterFlagCompletionFunc("priority", shellcomp.CompletionFunc(priorityRange))
}
