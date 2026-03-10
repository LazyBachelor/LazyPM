package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/charmbracelet/huh"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client *mongo.Client
}

func NewMongoStorage(ctx context.Context, uri, username, password string) (*MongoStorage, error) {
	credentials := options.Credential{
		Username: username,
		Password: password,
	}

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(uri).SetAuth(credentials))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("cannot reach MongoDB: %v", err)
	}

	return &MongoStorage{client: client}, nil
}

func NewMongoStorageInteractive(ctx context.Context, uri string) (*MongoStorage, error) {
	var username, password string
	if os.Getenv("DB_USER") == "" {
		if err := huh.NewInput().
			Title("Enter the Database Username").
			Value(&username).
			WithTheme(huh.ThemeBase16()).Run(); err != nil {
			return nil, fmt.Errorf("failed to read username: %w", err)
		}
	} else {
		username = os.Getenv("DB_USER")
	}

	if username == "" {
		return nil, fmt.Errorf("No username provided.")
	}

	if os.Getenv("DB_PASSWORD") == "" {
		if err := huh.NewInput().
			Title("Enter the Survey Password").
			EchoMode(huh.EchoModePassword).
			Value(&password).
			WithTheme(huh.ThemeBase16()).Run(); err != nil {
			return nil, fmt.Errorf("failed to read password: %w", err)
		}
	} else {
		password = os.Getenv("DB_PASSWORD")
	}

	if password == "" {
		return nil, fmt.Errorf("No password provided.")

	}

	mongoClient, err := NewMongoStorage(ctx, uri, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return mongoClient, nil
}

func (s *MongoStorage) Close() error {
	return s.client.Disconnect(context.Background())
}

func (s *MongoStorage) SubmitSurveyResponsesCmd(ctx context.Context) error {
	pmDir := "./.pm/"
	userStatscollection := s.client.Database("Responses").Collection("stats")
	taskMetricsCollection := s.client.Database("Responses").Collection("metrics")

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

	_, err = userStatscollection.UpdateOne(ctx,
		bson.M{"_id": stats.ID},
		bson.M{"$set": stats},
		options.Update().SetUpsert(true))

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

		_, err = taskMetricsCollection.UpdateOne(ctx, bson.M{"_id": metrics.ID}, bson.M{"$set": metrics},
			options.Update().SetUpsert(true))

		if err != nil {
			fmt.Printf("failed to insert metrics from %s: %v", file, err)
			continue
		}
	}

	return nil
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
