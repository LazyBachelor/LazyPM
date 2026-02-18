package commands

import (
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/spf13/cobra"
)

// replInstance holds a reference to the REPL for accessing validation feedback
var replInstance interface {
	GetCurrentFeedback() task.ValidationFeedback
}

// SetRepl sets the REPL instance for use by commands
func SetRepl(repl interface {
	GetCurrentFeedback() task.ValidationFeedback
}) {
	replInstance = repl
}

// StatusCmd displays the current task validation status
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check task validation status",
	Long:  "Displays the current task validation status and feedback.",
	RunE:  runStatusCmd,
}

func runStatusCmd(cmd *cobra.Command, args []string) error {
	if replInstance == nil {
		cmd.Println("No task validation available")
		return nil
	}

	feedback := replInstance.GetCurrentFeedback()
	if feedback.Message == "" {
		cmd.Println("No validation status available yet.")
		return nil
	}

	cmd.Print(feedback.Message)

	return nil
}
func init() {
	rootCmd.AddCommand(StatusCmd)
}
