package surveyCmd

import "github.com/spf13/cobra"

var IssuesCmd = &cobra.Command{
	Use:     "issues",
	Aliases: []string{"i"},
	Short:   "Commands related to issues",
}
