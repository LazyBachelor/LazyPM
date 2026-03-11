package survey

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/storage"
	"github.com/spf13/cobra"
)

var SubmitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit your survey responses",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := AppFromContext(cmd.Context())

		if app == nil {
			return fmt.Errorf("application context not initialized")
		}

		if app.Config.DbUri == "" {
			cmd.Println("No database URI provided in environment, survey responses will not be submitted.")
			return nil
		}

		db, err := storage.NewMongoStorageInteractive(cmd.Context(), app.Config.DbUri)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		defer db.Close()

		if err := db.SubmitSurveyResponsesCmd(cmd.Context(), app.Config.AppDir); err != nil {
			return fmt.Errorf("failed to submit survey responses: %w", err)
		}

		cmd.Println("Successfully submitted survey responses and metrics to the database")
		return nil
	},
}
