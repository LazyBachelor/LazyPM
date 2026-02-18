package main

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/spf13/cobra"
)

var interfaceType string
var stage int

var rootCmd = &cobra.Command{
	Use: "survey",
	Long: `Project Management Interface Survey

Thank you for participating in our survey!

We are gathering data on how users interact with different task management interfaces to better understand their preferences and compare their usability.
This survey will present you with a series of tasks to complete using various interfaces, including command-line, web-based, and terminal user interfaces.

Please answer the questions honestly and to the best of your ability.
Your responses will be kept confidential and used solely for research purposes.`}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the user survey",
	RunE:  runStartCmd,
}

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit your survey responses",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Submitting responses and metrics...")
		return nil
	},
}

func runStartCmd(cmd *cobra.Command, args []string) error {
	interfaces := initInterfaces()

	svc, cleanup, err := initializeServices(cmd.Context())
	if err != nil {
		return returnIfUserQuit(err, "failed to initialize services")
	}
	defer cleanup()

	surveyTasks := initTasks(svc)

	if cmd.Flags().Changed("interface") {
		if _, ok := interfaces[interfaceType]; !ok {
			return fmt.Errorf("invalid interface, valid are (tui, repl, web)")
		}
		interfaces = map[string]task.Interface{
			interfaceType: interfaces[interfaceType],
		}
	}

	if cmd.Flags().Changed("stage") {
		if stage < 1 || stage > len(surveyTasks) {
			return fmt.Errorf("invalid stage")
		}
		if err := runTask(cmd.Context(), surveyTasks[stage-1], interfaces[interfaceType]); err != nil {
			return err
		}
		return nil
	}

	if err := newIntroModel().Run(); err != nil {
		return returnIfUserQuit(err, "failed to run intro")
	}

	if err := taskLoop(cmd.Context(), surveyTasks, interfaces); err != nil {
		return returnIfUserQuit(err, "task loop failed")
	}
	return nil
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	startCmd.Flags().StringVarP(&interfaceType, "interface", "i", "tui", "Specify interface.")
	startCmd.Flags().IntVarP(&stage, "stage", "s", 1, "Run stage directly")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(submitCmd)
}
