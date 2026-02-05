package commands

import (
	"strings"

	"github.com/spf13/cobra"
)

var getIssueCmd = &cobra.Command{
	Use:               "describe [issue ID]",
	Aliases:           []string{"get", "read"},
	Short:             "Gets issue details",
	Long:              `Gets issue details by ID`,
	RunE:              runGetCmd,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeIssues,
}

func runGetCmd(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	issue, err := svc.Beads.GetIssue(cmd.Context(), issueID)
	if err != nil {
		return err
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
	if svc == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	issues, err := svc.Beads.AllIssues(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	for _, issue := range issues {

		if strings.HasPrefix(issue.Title, toComplete) {
			completions = append(completions, issue.ID)
		}

		if strings.HasPrefix(issue.ID, toComplete) {
			completions = append(completions, issue.ID)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
