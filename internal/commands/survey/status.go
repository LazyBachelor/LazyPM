package survey

import (
	"github.com/LazyBachelor/LazyPM/internal/commands/issues"
	"github.com/spf13/cobra"
)

// StatusCmd displays the current task validation status
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check task validation status",
	Long:  "Displays the current task validation status and feedback.",
	RunE:  runStatusCmd,
}

func runStatusCmd(cmd *cobra.Command, args []string) error {
	app := issues.AppFromContext(cmd.Context())
	if app == nil || app.CurrentFeedback == nil {
		cmd.Println("No validation status available.")
		return nil
	}

	if app.CurrentFeedback.Message == "" {
		cmd.Println("No validation status available yet.")
		return nil
	}

	cmd.Println(app.CurrentFeedback.Message)
	for _, check := range app.CurrentFeedback.Checks {
		if check.Valid {
			cmd.Printf("✅ %s\n", check.Message)
		} else {
			cmd.Printf("❌ %s\n", check.Message)
		}
	}
	return nil
}
