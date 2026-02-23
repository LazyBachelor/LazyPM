package issuesCmd

import (
	"context"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
)

// Variables for completion options and functions.
var (
	typeOptions   = []string{"bug", "feature", "task"}
	statusOptions = []string{"open", "closed", "in_progress"}
	priorityRange = []string{"0", "1", "2", "3", "4"}
)

// completionFunc returns a function that provides shell completion for the given options.
func completionFunc(options []string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return options, cobra.ShellCompDirectiveDefault
	}
}

// completeIssues provides shell completion for issue IDs and titles.
func completeIssues(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	issues, _ := GetIssueCompletions(cmd.Context(), toComplete)
	var ids []string
	for _, issue := range issues {
		ids = append(ids, issue.ID)
	}
	return ids, cobra.ShellCompDirectiveNoFileComp
}

// GetIssueCompletions fetches issues matching the toComplete string for shell completion.
func GetIssueCompletions(ctx context.Context, toComplete string) ([]models.Issue, cobra.ShellCompDirective) {
	app := AppFromContext(ctx)
	if app == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	issues, err := app.Issues.AllIssues(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []models.Issue
	for _, issue := range issues {
		if strings.HasPrefix(issue.ID, toComplete) {
			completions = append(completions, issue)
		} else if strings.HasPrefix(issue.Title, toComplete) {
			completions = append(completions, issue)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
