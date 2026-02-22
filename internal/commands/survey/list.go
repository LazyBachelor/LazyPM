package surveyCmd

import (
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/spf13/cobra"
)

var ListTasksCmd = &cobra.Command{
	Use:     "list-tasks",
	Aliases: []string{"ls-t"},
	Short:   "List available tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		for i, name := range task.ListTasks() {
			cmd.Printf("%d. %s\n", i+1, name)
		}
		return nil
	},
}

var ListInterfacesCmd = &cobra.Command{
	Use:     "list-interfaces",
	Aliases: []string{"ls-i"},
	Short:   "List available interfaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		for i, name := range task.ListInterfaces() {
			cmd.Printf("%d. %s\n", i+1, name)
		}
		return nil
	},
}
