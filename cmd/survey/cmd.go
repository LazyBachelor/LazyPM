package main

import (
	"github.com/LazyBachelor/LazyPM/pkg/cli/repl"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/LazyBachelor/LazyPM/pkg/tui"
	"github.com/LazyBachelor/LazyPM/pkg/web"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "survey",
	Short: "This application exists to gather metrics and feedback on task management interfaces.",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the user survey",
	RunE:  runStartCmd,
	Args:  cobra.MinimumNArgs(0),
}

func runStartCmd(cmd *cobra.Command, args []string) error {
	if err := newIntroModel().Run(); err != nil {
		return returnIfUserQuit(err, "failed to run intro")
	}

	svc, cleanup, err := initializeServices(cmd.Context())
	if err != nil {
		return returnIfUserQuit(err, "failed to initialize services")
	}
	defer cleanup()

	surveyTasks := initTasks(svc)
	interfaces := initInterfaces()

	if args[0] == "tui" {
		interfaces = []task.Interface{tui.NewTui()}
	}
	if args[0] == "repl" {
		interfaces = []task.Interface{repl.NewRepl()}
	}
	if args[0] == "web" {
		interfaces = []task.Interface{web.NewWeb()}
	}

	if err := taskLoop(cmd.Context(), surveyTasks, interfaces); err != nil {
		return returnIfUserQuit(err, "task loop failed")
	}
	return nil
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(startCmd)
}
