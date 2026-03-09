package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/LazyBachelor/LazyPM/cmd/pm/tasks"
	"github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/LazyBachelor/LazyPM/internal/storage"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func runStartCmd(cmd *cobra.Command, args []string) error {
	app := survey.AppFromContext(cmd.Context())
	if app == nil {
		if err := ensureAppInitialized(cmd.Context()); err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}
		app = survey.AppFromContext(cmd.Context())
		if app == nil {
			app = App
		}
	}

	if app == nil {
		return fmt.Errorf("application services are not available")
	}

	if len(args) == 0 {
		var noSubmit bool
		cmd.Println("No MongoDB password provided, survey responses will not be submitted.")
		cmd.Println("you can manually submit later with `pm survey submit <mongo-password>`.")

		huh.NewConfirm().Title("Are you sure you want to continue without submiting data?").
			Value(&noSubmit).RunAccessible(cmd.OutOrStdout(), cmd.InOrStdin())

		if !noSubmit {
			cmd.Println("Exiting. Please run the command again with a MongoDB password to submit your data.")
			return nil
		}
	} else {
		mongoPass := args[0]
		mongoStorage, err := storage.NewMongoStorage(app.Config.MongoURI, "participant", mongoPass)
		if err != nil {
			cmd.Println("Failed to connect to MongoDB, survey responses will not be submitted.")
			cmd.Println("You can manually submit later with `pm survey submit <password>`.")
			return nil
		} else {
			cmd.Println("Connected to MongoDB successfully.")

			defer mongoStorage.Close()

			ctx := cmd.Context()

			go func() {
				ticker := time.NewTicker(10 * time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						if err := mongoStorage.SubmitSurveyResponsesCmd(ctx); err != nil {
							cmd.Printf("Failed to submit survey responses: %v\n", err)
						}
					}
				}
			}()

			if err := mongoStorage.SubmitSurveyResponsesCmd(ctx); err != nil {
				cmd.Printf("Failed to submit survey responses: %v\n", err)
			}
		}
	}

	interfaces := initInterfaces()
	surveyTasks := initTasks(app)

	if cmd.Flags().Changed("interface") {
		if _, ok := interfaces[survey.InterfaceType]; !ok {
			return fmt.Errorf("invalid interface, valid are %v", task.ListInterfaces())
		}
		interfaces = map[string]task.Interface{
			survey.InterfaceType: interfaces[survey.InterfaceType],
		}
	}

	if cmd.Flags().Changed("task") {
		if surveyTask := surveyTasks[survey.Task]; surveyTask == nil {
			return fmt.Errorf("invalid task, valid are %v", task.ListTasks())
		}

		surveyTasks = map[string]task.Tasker{
			survey.Task: surveyTasks[survey.Task],
		}
	}

	if err := newIntroModel().Run(); err != nil {
		return returnIfUserQuit(err, "failed to run intro")
	}

	if err := taskLoop(cmd.Context(), app, surveyTasks, interfaces); err != nil {
		return returnIfUserQuit(err, "task loop failed")
	}
	return nil
}

func taskLoop(ctx context.Context, application *task.App, surveyTasks map[string]task.Tasker, interfaces map[string]task.Interface) error {
	var iNames []string
	for name := range interfaces {
		iNames = append(iNames, name)
	}

	if len(iNames) == 0 {
		return fmt.Errorf("no interfaces are available")
	}

	if len(surveyTasks) == 0 {
		return fmt.Errorf("no tasks are available")
	}

	rand.Shuffle(len(iNames), func(i, j int) {
		iNames[i], iNames[j] = iNames[j], iNames[i]
	})

	idx := 0
	for _, t := range surveyTasks {
		iIdx := idx % len(iNames)
		selected := interfaces[iNames[iIdx]]

		runner := task.NewTaskRunner(application)

		if err := runner.Run(ctx, t, selected, tasks.InterfaceToType(selected)); err != nil {
			return err
		}
		idx++
	}
	return nil
}

func returnIfUserQuit(err error, msg string) error {
	if errors.Is(err, task.ErrUserQuit) {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
