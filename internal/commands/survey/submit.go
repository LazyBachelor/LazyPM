package survey

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI = os.Getenv("MONGODB_URI")

var SubmitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit your survey responses",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, error := mongo.Connect(cmd.Context(), options.Client().ApplyURI(mongoURI))
		if error != nil {
			return fmt.Errorf("Failed to connect to MongoDB: %v", error)
		}

		go func() {
			if err := client.Disconnect(cmd.Context()); err != nil {
				fmt.Printf("Failed to disconnect MongoDB client: %v", err)
			}
		}()

		userStatscollection := client.Database("Responses").Collection("stats")
		taskMetricsCollection := client.Database("Responses").Collection("task_metrics")

		pmDir := "./.pm/"

		entries, err := os.ReadDir(pmDir)
		if err != nil {
			return fmt.Errorf("Failed to read .pm directory: %v", err)
		}

		if len(entries) == 0 {
			return fmt.Errorf("No files found in .pm directory")
		}

		statFile := pmDir + "stats.json"
		if _, err := os.Stat(statFile); os.IsNotExist(err) {
			return fmt.Errorf("stats.json not found in .pm directory")
		}

		stats, err := getStats(statFile)
		if err != nil {
			return fmt.Errorf("Failed to read stats.json: %v", err)
		}

		_, err = userStatscollection.InsertOne(cmd.Context(), stats)
		if err != nil {
			return fmt.Errorf("Failed to insert stats into database: %v", err)
		}

		metricFiles := []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			if entry.Name() == "stats.json" {
				continue
			}

			if strings.HasSuffix(entry.Name(), "-stats.json") {
				metricFiles = append(metricFiles, pmDir+entry.Name())
				continue
			}
		}

		for _, file := range metricFiles {
			metrics, err := getTaskMetrics(file)
			if err != nil {
				fmt.Printf("failed to read metrics from %s: %v", file, err)
				continue
			}

			_, err = taskMetricsCollection.InsertOne(cmd.Context(), metrics)
			if err != nil {
				fmt.Printf("failed to insert metrics from %s: %v", file, err)
				continue
			}
		}

		cmd.Printf("Successfully submitted survey responses and metrics to the database")

		return nil
	},
}

func getStats(file string) (*models.Statistics, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var stats models.Statistics
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func getTaskMetrics(file string) (*models.TaskMetricsFile, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var metrics models.TaskMetricsFile
	if err := json.Unmarshal(data, &metrics); err != nil {
		return nil, err
	}

	return &metrics, nil
}
