package app

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/storage"
)

type StatisticsService struct {
	storage *storage.Storage[models.Statistics]
	logger  *slog.Logger
	mu      sync.Mutex
}

func NewStatisticsService(storage *storage.Storage[models.Statistics], logger *slog.Logger) (*StatisticsService, error) {
	if err := storage.Init(); err != nil {
		return nil, err
	}
	return &StatisticsService{
		storage: storage,
		logger:  logger,
	}, nil
}

func (s *StatisticsService) Load(ctx context.Context) error {
	return s.storage.Load()
}

func (s *StatisticsService) Save(ctx context.Context) error {
	return s.storage.Save()
}

func (s *StatisticsService) GetStatistics() (models.Statistics, error) {
	if s.storage.Data == nil {
		return models.Statistics{}, errors.New("statistics data not initialized")
	}
	return *s.storage.Data, nil
}

func (s *StatisticsService) RecordTaskRun(ctx context.Context, run models.TaskRunMetrics) error {
	_ = ctx

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storage.Data == nil {
		return errors.New("statistics data not initialized")
	}

	stats := s.storage.Data
	now := time.Now()

	if stats.StartTime.IsZero() {
		stats.StartTime = run.StartedAt
	}
	if run.StartedAt.Before(stats.StartTime) {
		stats.StartTime = run.StartedAt
	}

	stats.EndTime = now
	stats.Duration = stats.EndTime.Sub(stats.StartTime)
	stats.InterfaceType = run.InterfaceType

	stats.TaskRuns++
	stats.LastTaskName = run.TaskName
	stats.LastRunID = run.RunID

	if run.Completed {
		stats.TasksCompleted++
	} else {
		stats.TasksFailed++
	}

	stats.TotalDurationMs += run.DurationMs
	if stats.TaskRuns > 0 {
		stats.AverageDurationMs = stats.TotalDurationMs / int64(stats.TaskRuns)
	}

	stats.ValidationAttempts += run.ValidationAttempts
	stats.ValidationSuccesses += run.ValidationSuccesses
	stats.ValidationFailures += run.ValidationFailures
	stats.ValidationChecksPassed += run.ValidationChecksPassed
	stats.ValidationChecksFailed += run.ValidationChecksFailed

	userActions := 0
	for _, log := range run.Logs {
		if log.Level == "user_action" {
			userActions++
		}
	}
	stats.TotalUserActions += userActions
	stats.ButtonClicks.Clicks += userActions
	if run.QuestionnaireCompleted {
		stats.QuestionnairesCompleted++
	}
	if run.QuestionnaireUserQuit {
		stats.QuestionnairesAbandoned++
	}

	if err := s.storage.Save(); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to save global statistics", "error", err)
		}
		return err
	}

	if s.logger != nil {
		s.logger.Info("global statistics updated",
			"task", run.TaskName,
			"run_id", run.RunID,
			"task_runs", stats.TaskRuns,
			"completed", stats.TasksCompleted,
			"failed", stats.TasksFailed,
		)
	}

	return nil
}
