package commands

import (
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
)

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

var getIssuesCmd = &cobra.Command{
	Use:     "ls [search query]",
	Short:   "List all issues",
	Long:    `List all issues in the project management system.`,
	Aliases: []string{"list", "search"},
	Example: lsExamples,
	Args:    cobra.MinimumNArgs(0),
	RunE:    runGetIssuesCmd,
}

func runGetIssuesCmd(cmd *cobra.Command, args []string) error {

	queryArg := strings.Join(args, " ")

	filter := models.IssueFilter{
		TitleSearch:         titleFlag,
		DescriptionContains: descriptionFlag,
		Limit:               limit,
	}

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

	issuesPtr, err := svc.Beads.SearchIssues(cmd.Context(), queryArg, filter)
	if err != nil {
		return err
	}

	issues := models.IssuesPtrToIssues(issuesPtr)

	models.PrintIssues(issues)

	return nil
}

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
