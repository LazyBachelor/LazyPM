package task

import (
	"log/slog"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RunLifecycle struct {
	collector    *taskRunCollector
	config       Config
	details      models.TaskDetails
	app          *App
	logger       *slog.Logger
	metricsStore MetricsStore
}

func NewRunLifecycle(app *App, config Config, details models.TaskDetails, iType InterfaceType, logger *slog.Logger) *RunLifecycle {

	collector := newTaskRunCollector(details.Title, iType, logger)

	config = config.WithActionLogger(func(action string) {
		collector.recordUserAction(action)
	})

	var participantID bson.ObjectID
	var store MetricsStore
	if config.StatisticsStoragePath != "" {
		participantID = app.Stats.GetParticipantID()
		store = NewFileMetricsStore(config.StatisticsStoragePath, participantID, logger)
	}

	return &RunLifecycle{
		collector:    collector,
		config:       config,
		details:      details,
		app:          app,
		logger:       logger,
		metricsStore: store,
	}
}
