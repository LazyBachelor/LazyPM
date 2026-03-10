package survey

import (
	"github.com/LazyBachelor/LazyPM/internal/utils/shellcomp"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/spf13/cobra"
)

var (
	DevFlag       bool
	InterfaceType string
	Task          string
)

// StartCmd is the start command - RunE is set in cmd/survey/
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the user survey",
}

func init() {
	StartCmd.Flags().StringVarP(&Task, "task", "t", "", "Specify task.")
	StartCmd.Flags().StringVarP(&InterfaceType, "interface", "i", "", "Specify interface.")
	StartCmd.RegisterFlagCompletionFunc("task", shellcomp.CompletionFunc(task.ListTasks()))
	StartCmd.RegisterFlagCompletionFunc("interface", shellcomp.CompletionFunc(task.ListInterfaces()))
	StartCmd.Flags().BoolVar(&DevFlag, "dev", false, "Enable development mode, which skips database connection, submission and intro")
}
