package survey

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/storage"
	"github.com/spf13/cobra"
)

var SubmitCmd = &cobra.Command{
	Use:   "submit <mongo-password>",
	Short: "Submit your survey responses",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := AppFromContext(cmd.Context())

		if app == nil {
			return fmt.Errorf("application context not initialized")
		}

		if len(args) == 0 {
			cmd.Println("No MongoDB password provided, survey responses will not be submitted.")
			return nil
		}

		mongoPassword := args[0]
		mongoClient, err := storage.NewMongoStorage(app.Config.MongoURI, "participant", mongoPassword)
		if err != nil {
			return fmt.Errorf("failed to connect to MongoDB: %w", err)
		}
		defer mongoClient.Close()

		if err := mongoClient.SubmitSurveyResponsesCmd(cmd.Context()); err != nil {
			return fmt.Errorf("failed to submit survey responses: %w", err)
		}

		cmd.Println("Successfully submitted survey responses and metrics to the database")
		return nil
	},
}
