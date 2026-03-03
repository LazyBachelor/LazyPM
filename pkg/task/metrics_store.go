package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
)

type MetricsStore interface {
	Append(ctx context.Context, taskName string, run models.TaskRunMetrics) error
}

type FileMetricsStore struct {
	path   string
	logger *slog.Logger
}

func NewFileMetricsStore(path string, logger *slog.Logger) *FileMetricsStore {
	return &FileMetricsStore{
		path:   path,
		logger: logger,
	}
}

func (s *FileMetricsStore) Append(ctx context.Context, taskName string, run models.TaskRunMetrics) error {

	if s.path == "" {
		return nil
	}

	dir := filepath.Dir(s.path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create metrics directory: %w", err)
		}
	}

	metrics := models.TaskMetricsFile{
		TaskName: taskName,
		Runs:     []models.TaskRunMetrics{},
	}

	if bytes, err := os.ReadFile(s.path); err == nil {
		if len(bytes) > 0 {
			if err := json.Unmarshal(bytes, &metrics); err != nil {
				return fmt.Errorf("parse metrics file: %w", err)
			}
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read metrics file: %w", err)
	}

	if metrics.TaskName == "" {
		metrics.TaskName = taskName
	}

	run.RunID = len(metrics.Runs) + 1
	metrics.Runs = append(metrics.Runs, run)
	metrics.Summary = buildTaskStatsSummary(metrics.Runs)
	metrics.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return fmt.Errorf("encode metrics: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("write metrics file: %w", err)
	}

	if s.logger != nil {
		s.logger.Info(
			"task metrics persisted",
			"path", s.path,
			"task", taskName,
			"run_id", run.RunID,
		)
	}

	return nil
}
