package commands

import (
	"context"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
)

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
	if svc == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	issues, err := svc.Beads.AllIssues(ctx)
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
