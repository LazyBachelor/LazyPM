package surveyCmd

import (
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		for i, name := range task.List() {
			cmd.Printf("%d. %s\n", i+1, name)
		}
		return nil
	},
}
