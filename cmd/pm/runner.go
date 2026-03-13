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

	if !cmd.Flags().Changed("dev") {
		var mongoStorage *storage.MongoStorage
		var continueWithoutSubmitting bool

		for {

			if app.Config.DbUri == "" {
				cmd.Println("No database URI provided in environment, survey responses will not be submitted.")
				break
			}

			db, err := storage.NewMongoStorageInteractive(cmd.Context(), app.Config.DbUri)
			if err == nil {
				mongoStorage = db
				break
			}

			cmd.Println("Failed to connect to database, survey responses will not be submitted.")

			if err := huh.NewConfirm().
				Title("Do you want to continue without submitting your responses?").
				Description("You can fix your database connection and submit your responses later with the submit command.").
				Value(&continueWithoutSubmitting).
				WithTheme(huh.ThemeBase16()).
				RunAccessible(cmd.OutOrStdout(), cmd.InOrStdin()); err != nil {
				return fmt.Errorf("failed to read user input: %w", err)
			}

			if continueWithoutSubmitting {
				break
			}
		}

		if mongoStorage != nil {
			cmd.Println("Connected to Database Successfully. Starting survey...")
			time.Sleep(2 * time.Second)

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
						if err := mongoStorage.SubmitSurveyResponsesCmd(ctx, app.Config.AppDir); err != nil {
							cmd.Printf("Failed to submit survey responses: %v\n", err)
						}
					}
				}
			}()

			if err := mongoStorage.SubmitSurveyResponsesCmd(ctx, app.Config.AppDir); err != nil {
				cmd.Printf("Failed to submit survey responses: %v\n", err)
			}

			defer func() {
				if err := mongoStorage.SubmitSurveyResponsesCmd(context.Background(), app.Config.AppDir); err != nil {
					cmd.Printf("Failed to submit survey responses on shutdown: %v\n", err)
				}
			}()

		} else {
			cmd.Println("Starting survey without database connection. Your responses will not be submitted...")
			time.Sleep(2 * time.Second)
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
	if !cmd.Flags().Changed("dev") {
		if err := newIntroModel().Run(); err != nil {
			return returnIfUserQuit(err, "failed to run intro")
		}

		introAnswers, err := newIntroQuestionnaire().Run()
		if err != nil {
			return returnIfUserQuit(err, "failed to run intro questionnaire")
		}

		if app != nil && app.Stats != nil && introAnswers != nil {
			if err := app.Stats.RecordIntroQuestionnaireAnswers(introAnswers); err != nil {
				cmd.Printf("Failed to record intro questionnaire answers: %v\n", err)
			}
		}
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

	taskNames := task.ListTasks()
	if len(taskNames) == 0 {
		return fmt.Errorf("no tasks are available")
	}

	rand.Shuffle(len(iNames), func(i, j int) {
		iNames[i], iNames[j] = iNames[j], iNames[i]
	})

	idx := 0
	for _, taskName := range taskNames {
		t, ok := surveyTasks[taskName]
		if !ok {
			continue
		}
		iIdx := idx % len(iNames)
		selected := interfaces[iNames[iIdx]]

		if selected == nil {
			return fmt.Errorf("interface %q is nil (available: %v)", iNames[iIdx], iNames)
		}

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
