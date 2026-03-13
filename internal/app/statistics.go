package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return models.Statistics{}, fmt.Errorf("statistics data not initialized")
	}
	return *s.storage.Data, nil
}

func (s *StatisticsService) GetParticipantID() primitive.ObjectID {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storage.Data == nil {
		return primitive.NilObjectID
	}
	return s.storage.Data.ID
}

func (s *StatisticsService) RecordTaskRun(ctx context.Context, run models.TaskRunMetrics) error {
	_ = ctx

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storage.Data == nil {
		return fmt.Errorf("statistics data not initialized")
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
	stats.DurationMs = stats.EndTime.Sub(stats.StartTime).Milliseconds()
	stats.LastInterfaceType = run.InterfaceType

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

func (s *StatisticsService) RecordIntroQuestionnaireAnswers(answers map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storage.Data == nil {
		return fmt.Errorf("statistics data not initialized")
	}

	s.storage.Data.IntroQuestionnaireAnswers = answers

	if err := s.storage.Save(); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to save intro questionnaire answers", "error", err)
		}
		return err
	}

	if s.logger != nil {
		s.logger.Info("intro questionnaire answers saved",
			"answers_count", len(answers),
		)
	}

	return nil
}
