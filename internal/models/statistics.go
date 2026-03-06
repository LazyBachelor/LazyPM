package models

import (
	"time"
)

type Statistics struct {
	ID        int           `json:"id"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`

	InterfaceType  InterfaceType `json:"interface_type"`
	TaskRuns       int           `json:"task_runs"`
	TasksCompleted int           `json:"tasks_completed"`
	TasksFailed    int           `json:"tasks_failed"`

	TotalDurationMs   int64 `json:"total_duration_ms"`
	AverageDurationMs int64 `json:"average_duration_ms"`

	TotalUserActions        int `json:"total_user_actions"`
	QuestionnairesCompleted int `json:"questionnaires_completed"`
	QuestionnairesAbandoned int `json:"questionnaires_abandoned"`

	ValidationAttempts     int `json:"validation_attempts"`
	ValidationSuccesses    int `json:"validation_successes"`
	ValidationFailures     int `json:"validation_failures"`
	ValidationChecksPassed int `json:"validation_checks_passed"`
	ValidationChecksFailed int `json:"validation_checks_failed"`

	LastTaskName string       `json:"last_task_name"`
	LastRunID    int          `json:"last_run_id"`
}

type TaskMetricsFile struct {
	TaskName  string           `json:"task_name"`
	UpdatedAt time.Time        `json:"updated_at"`
	Summary   TaskStatsSummary `json:"summary"`
	Runs      []TaskRunMetrics `json:"runs"`
}

type TaskStatsSummary struct {
	TotalRuns               int           `json:"total_runs"`
	CompletedRuns           int           `json:"completed_runs"`
	IncompleteRuns          int           `json:"incomplete_runs"`
	TotalDurationMs         int64         `json:"total_duration_ms"`
	AverageDurationMs       int64         `json:"average_duration_ms"`
	TotalUserActions        int           `json:"total_user_actions"`
	QuestionnairesCompleted int           `json:"questionnaires_completed"`
	QuestionnairesAbandoned int           `json:"questionnaires_abandoned"`
	ValidationAttempts      int           `json:"validation_attempts"`
	ValidationSuccesses     int           `json:"validation_successes"`
	ValidationFailures      int           `json:"validation_failures"`
	ValidationChecksPassed  int           `json:"validation_checks_passed"`
	ValidationChecksFailed  int           `json:"validation_checks_failed"`
	LastInterfaceType       InterfaceType `json:"last_interface_type"`
	FirstRunStartedAt       time.Time     `json:"first_run_started_at"`
	LastRunStartedAt        time.Time     `json:"last_run_started_at"`
	LastRunEndedAt          time.Time     `json:"last_run_ended_at"`
}

type TaskRunMetrics struct {
	RunID                  int            `json:"run_id"`
	TaskName               string         `json:"task_name"`
	InterfaceType          InterfaceType  `json:"interface_type"`
	StartedAt              time.Time      `json:"started_at"`
	EndedAt                time.Time      `json:"ended_at"`
	DurationMs             int64          `json:"duration_ms"`
	Completed              bool           `json:"completed"`
	ValidationAttempts     int            `json:"validation_attempts"`
	ValidationSuccesses    int            `json:"validation_successes"`
	ValidationFailures     int            `json:"validation_failures"`
	ValidationChecksPassed int            `json:"validation_checks_passed"`
	ValidationChecksFailed int            `json:"validation_checks_failed"`
	LastValidationMessage  string         `json:"last_validation_message,omitempty"`
	QuestionnaireCompleted bool           `json:"questionnaire_completed"`
	QuestionnaireUserQuit  bool           `json:"questionnaire_user_quit"`
	QuestionnaireAnswers   map[string]any `json:"questionnaire_answers,omitempty"`
	Logs                   []TaskLogEntry `json:"logs"`
	Error                  string         `json:"error,omitempty"`
}

type TaskLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source,omitempty"`
	Action    string    `json:"action,omitempty"`
	Target    string    `json:"target,omitempty"`
	Result    string    `json:"result,omitempty"`
	Attempt   int       `json:"attempt,omitempty"`
}
