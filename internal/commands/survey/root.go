package surveyCmd

import (
	"github.com/spf13/cobra"
)

var (
	InterfaceType string
	Task          int
)

// RootCmd is the base command for the survey CLI.
var RootCmd = &cobra.Command{
	Use: "survey",
	Long: `Project Management Interface Survey

Thank you for participating in our survey!

We are gathering data on how users interact with different task management interfaces to better understand their preferences and compare their usability.
This survey will present you with a series of tasks to complete using various interfaces, including command-line, web-based, and terminal user interfaces.

Please answer the questions honestly and to the best of your ability.
Your responses will be kept confidential and used solely for research purposes.`}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	StartCmd.Flags().StringVarP(&InterfaceType, "interface", "i", "tui", "Specify interface.")
	StartCmd.Flags().IntVarP(&Task, "task", "t", 1, "Run task directly")
}
