package surveyCmd

import "github.com/spf13/cobra"

var (
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
}
