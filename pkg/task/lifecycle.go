package task

import (
	"log/slog"

	"github.com/LazyBachelor/LazyPM/internal/models"
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

	var store MetricsStore
	if config.StatisticsStoragePath != "" {
		store = NewFileMetricsStore(config.StatisticsStoragePath, logger)
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
