package main

import (
	"errors"
	"log"
	"os"

	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "survey",
	Short: "Run the user survey",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the user survey",
	Run: func(cmd *cobra.Command, args []string) {
		if err := newIntroModel().Run(); err != nil {
			if errors.Is(err, ErrUserQuit) {
				os.Exit(0)
			}
			log.Fatalf("Failed to run intro screen: %v\n", err)
		}

		svc, close, err := initializeServices(cmd.Context())
		if err != nil {
			log.Fatalf("Failed to initialize services: %v\n", err)
		}
		defer close()

		surveyTasks := initTasks()
		interfaces := initInterfaces()

		if err := taskLoop(cmd.Context(), svc, surveyTasks, interfaces); err != nil {
			if errors.Is(err, task.ErrUserQuit) {
				os.Exit(0)
			}
			log.Fatalf("Task loop failed: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
