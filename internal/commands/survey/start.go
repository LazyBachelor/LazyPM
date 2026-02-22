package surveyCmd

import "github.com/spf13/cobra"

// StartCmd is the start command - RunE is set in cmd/survey/
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the user survey",
}
