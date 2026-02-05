package commands

import (
	"context"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
)

var getIssueCmd = &cobra.Command{
	Use:               "describe [issue ID]",
	Short:             "Get issue details",
	Long:              `Get issue details by ID`,
	Aliases:           []string{"get", "read"},
	Args:              cobra.ExactArgs(1),
	RunE:              runGetCmd,
	ValidArgsFunction: completeIssues,
}

func runGetCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	issue, err := svc.Beads.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return err
	}

	if issue == nil {
		cmd.Printf("Issue with ID '%s' not found\n", issueID)
		return nil
	}

	cmd.Printf("Title: %s\n", issue.Title)
	cmd.Printf("Description: %s\n", issue.Description)
	cmd.Printf("Status: %s\n", issue.Status)
	cmd.Printf("Type: %s\n", issue.IssueType)
	cmd.Printf("Priority: %d\n", issue.Priority)

	return nil
}

func init() {
	rootCmd.AddCommand(getIssueCmd)
}

func completeIssues(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	issues, _ := GetIssueCompletions(cmd.Context(), toComplete)
	var ids []string
	for _, issue := range issues {
		ids = append(ids, issue.ID)
	}
	return ids, cobra.ShellCompDirectiveNoFileComp
}

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
