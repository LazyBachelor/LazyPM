package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Statistics struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	StartTime  time.Time          `bson:"start_time" json:"start_time"`
	EndTime    time.Time          `bson:"end_time" json:"end_time"`
	DurationMs int64              `bson:"duration_ms" json:"duration_ms"`

	LastInterfaceType InterfaceType `bson:"last_interface_type" json:"last_interface_type"`
	TaskRuns          int           `bson:"task_runs" json:"task_runs"`
	TasksCompleted    int           `bson:"tasks_completed" json:"tasks_completed"`
	TasksFailed       int           `bson:"tasks_failed" json:"tasks_failed"`

	TotalDurationMs   int64 `bson:"total_duration_ms" json:"total_duration_ms"`
	AverageDurationMs int64 `bson:"average_duration_ms" json:"average_duration_ms"`

	TotalUserActions        int `bson:"total_user_actions" json:"total_user_actions"`
	QuestionnairesCompleted int `bson:"questionnaires_completed" json:"questionnaires_completed"`
	QuestionnairesAbandoned int `bson:"questionnaires_abandoned" json:"questionnaires_abandoned"`

	ValidationAttempts     int `bson:"validation_attempts" json:"validation_attempts"`
	ValidationSuccesses    int `bson:"validation_successes" json:"validation_successes"`
	ValidationFailures     int `bson:"validation_failures" json:"validation_failures"`
	ValidationChecksPassed int `bson:"validation_checks_passed" json:"validation_checks_passed"`
	ValidationChecksFailed int `bson:"validation_checks_failed" json:"validation_checks_failed"`

	LastTaskName string `bson:"last_task_name" json:"last_task_name"`
	LastRunID    int    `bson:"last_run_id" json:"last_run_id"`
}

type TaskMetricsFile struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	ParticipantID primitive.ObjectID `bson:"participant_id" json:"participant_id"`
	TaskName      string             `bson:"task_name" json:"task_name"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	Summary       TaskStatsSummary   `bson:"summary" json:"summary"`
	Runs          []TaskRunMetrics   `bson:"runs" json:"runs"`
}

type TaskStatsSummary struct {
	TotalRuns               int           `bson:"total_runs" json:"total_runs"`
	CompletedRuns           int           `bson:"completed_runs" json:"completed_runs"`
	IncompleteRuns          int           `bson:"incomplete_runs" json:"incomplete_runs"`
	TotalDurationMs         int64         `bson:"total_duration_ms" json:"total_duration_ms"`
	AverageDurationMs       int64         `bson:"average_duration_ms" json:"average_duration_ms"`
	TotalUserActions        int           `bson:"total_user_actions" json:"total_user_actions"`
	QuestionnairesCompleted int           `bson:"questionnaires_completed" json:"questionnaires_completed"`
	QuestionnairesAbandoned int           `bson:"questionnaires_abandoned" json:"questionnaires_abandoned"`
	ValidationAttempts      int           `bson:"validation_attempts" json:"validation_attempts"`
	ValidationSuccesses     int           `bson:"validation_successes" json:"validation_successes"`
	ValidationFailures      int           `bson:"validation_failures" json:"validation_failures"`
	ValidationChecksPassed  int           `bson:"validation_checks_passed" json:"validation_checks_passed"`
	ValidationChecksFailed  int           `bson:"validation_checks_failed" json:"validation_checks_failed"`
	LastInterfaceType       InterfaceType `bson:"last_interface_type" json:"last_interface_type"`
	FirstRunStartedAt       time.Time     `bson:"first_run_started_at" json:"first_run_started_at"`
	LastRunStartedAt        time.Time     `bson:"last_run_started_at" json:"last_run_started_at"`
	LastRunEndedAt          time.Time     `bson:"last_run_ended_at" json:"last_run_ended_at"`
}

type TaskRunMetrics struct {
	RunID                  int            `bson:"run_id" json:"run_id"`
	TaskName               string         `bson:"task_name" json:"task_name"`
	InterfaceType          InterfaceType  `bson:"interface_type" json:"interface_type"`
	StartedAt              time.Time      `bson:"started_at" json:"started_at"`
	EndedAt                time.Time      `bson:"ended_at" json:"ended_at"`
	DurationMs             int64          `bson:"duration_ms" json:"duration_ms"`
	Completed              bool           `bson:"completed" json:"completed"`
	ValidationAttempts     int            `bson:"validation_attempts" json:"validation_attempts"`
	ValidationSuccesses    int            `bson:"validation_successes" json:"validation_successes"`
	ValidationFailures     int            `bson:"validation_failures" json:"validation_failures"`
	ValidationChecksPassed int            `bson:"validation_checks_passed" json:"validation_checks_passed"`
	ValidationChecksFailed int            `bson:"validation_checks_failed" json:"validation_checks_failed"`
	LastValidationMessage  string         `bson:"last_validation_message,omitempty" json:"last_validation_message,omitempty"`
	QuestionnaireCompleted bool           `bson:"questionnaire_completed" json:"questionnaire_completed"`
	QuestionnaireUserQuit  bool           `bson:"questionnaire_user_quit" json:"questionnaire_user_quit"`
	QuestionnaireAnswers   map[string]any `bson:"questionnaire_answers,omitempty" json:"questionnaire_answers,omitempty"`
	Logs                   []TaskLogEntry `bson:"logs" json:"logs"`
	Error                  string         `bson:"error,omitempty" json:"error,omitempty"`
}

type TaskLogEntry struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Level     string    `bson:"level" json:"level"`
	Message   string    `bson:"message" json:"message"`
	Source    string    `bson:"source,omitempty" json:"source,omitempty"`
	Action    string    `bson:"action,omitempty" json:"action,omitempty"`
	Target    string    `bson:"target,omitempty" json:"target,omitempty"`
	Result    string    `bson:"result,omitempty" json:"result,omitempty"`
	Attempt   int       `bson:"attempt,omitempty" json:"attempt,omitempty"`
}
