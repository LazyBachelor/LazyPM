package task

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
)

type taskRunCollector struct {
	mu     sync.Mutex
	run    models.TaskRunMetrics
	logger *slog.Logger
}

var nonWord = regexp.MustCompile(`[^a-z0-9]+`)

func newTaskRunCollector(taskName string, interfaceType InterfaceType, logger *slog.Logger) *taskRunCollector {
	return &taskRunCollector{
		run: models.TaskRunMetrics{
			TaskName:      taskName,
			InterfaceType: interfaceType,
			StartedAt:     time.Now(),
			Logs:          make([]models.TaskLogEntry, 0, 8),
		},
		logger: logger,
	}
}

func (c *taskRunCollector) log(level string, message string) {
	action := normalizeAction(message)
	result := ""
	switch level {
	case "error":
		result = "failed"
	case "warn":
		result = "warning"
	default:
		result = "ok"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.run.Logs = append(c.run.Logs, models.TaskLogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Source:    "system",
		Action:    action,
		Result:    result,
	})

	attrs := []any{"action", action, "result", result, "task", c.run.TaskName, "interface", c.run.InterfaceType}
	c.logWithLogger(level, message, attrs...)
}

func (c *taskRunCollector) recordUserAction(raw string) {
	source, actionText, target, result := normalizeUserAction(raw)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.run.Logs = append(c.run.Logs, models.TaskLogEntry{
		Timestamp: time.Now(),
		Level:     "user_action",
		Message:   raw,
		Source:    source,
		Action:    normalizeAction(actionText),
		Target:    target,
		Result:    result,
	})

	c.logWithLogger("info", "user action recorded",
		"task", c.run.TaskName,
		"interface", c.run.InterfaceType,
		"source", source,
		"action", normalizeAction(actionText),
		"target", target,
		"result", result,
	)
}

func (c *taskRunCollector) setCompleted(completed bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.run.Completed = completed
}

func (c *taskRunCollector) setError(err error) {
	if err == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.run.Error = err.Error()
}

func (c *taskRunCollector) recordQuestionnaire(completed bool, userQuit bool, answers map[string]any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.run.QuestionnaireCompleted = completed
	c.run.QuestionnaireUserQuit = userQuit
	if len(answers) > 0 {
		c.run.QuestionnaireAnswers = answers
	}

	result := "completed"
	if userQuit {
		result = "user_quit"
	} else if !completed {
		result = "incomplete"
	}

	c.run.Logs = append(c.run.Logs, models.TaskLogEntry{
		Timestamp: time.Now(),
		Level:     "questionnaire",
		Message:   "questionnaire finished",
		Source:    "system",
		Action:    "questionnaire_finish",
		Result:    result,
	})

	if len(answers) > 0 {
		keys := make([]string, 0, len(answers))
		for key := range answers {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := answers[key]
			if value == nil {
				continue
			}

			valueText := fmt.Sprintf("%v", value)
			c.run.Logs = append(c.run.Logs, models.TaskLogEntry{
				Timestamp: time.Now(),
				Level:     "questionnaire",
				Message:   fmt.Sprintf("questionnaire answer: %s=%s", key, valueText),
				Source:    "system",
				Action:    "questionnaire_answer",
				Target:    key,
				Result:    valueText,
			})
		}
	}

	c.logWithLogger("info", "questionnaire recorded",
		"task", c.run.TaskName,
		"completed", completed,
		"user_quit", userQuit,
		"answers_count", len(answers),
	)
}

func (c *taskRunCollector) recordValidation(feedback ValidationFeedback) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.run.ValidationAttempts++
	attempt := c.run.ValidationAttempts
	if feedback.Success {
		c.run.ValidationSuccesses++
	} else {
		c.run.ValidationFailures++
	}

	c.run.Logs = append(c.run.Logs, models.TaskLogEntry{
		Timestamp: time.Now(),
		Level:     "validation",
		Message:   fmt.Sprintf("validation attempt %d success=%t", attempt, feedback.Success),
		Source:    "system",
		Action:    "validate_attempt",
		Result:    validationResult(feedback.Success),
		Attempt:   attempt,
	})
	c.logWithLogger("info", "validation attempt",
		"task", c.run.TaskName,
		"interface", c.run.InterfaceType,
		"attempt", attempt,
		"result", validationResult(feedback.Success),
	)

	for _, check := range feedback.Checks {
		if check.Valid {
			c.run.ValidationChecksPassed++
			c.run.Logs = append(c.run.Logs, models.TaskLogEntry{
				Timestamp: time.Now(),
				Level:     "validation",
				Message:   fmt.Sprintf("validation check passed: %s", check.Message),
				Source:    "system",
				Action:    "validate_check",
				Target:    check.Message,
				Result:    "passed",
				Attempt:   attempt,
			})
			c.logWithLogger("info", "validation check",
				"task", c.run.TaskName,
				"attempt", attempt,
				"result", "passed",
				"target", check.Message,
			)
			continue
		}
		c.run.ValidationChecksFailed++
		c.run.Logs = append(c.run.Logs, models.TaskLogEntry{
			Timestamp: time.Now(),
			Level:     "validation",
			Message:   fmt.Sprintf("validation check failed: %s", check.Message),
			Source:    "system",
			Action:    "validate_check",
			Target:    check.Message,
			Result:    "failed",
			Attempt:   attempt,
		})
		c.logWithLogger("warn", "validation check",
			"task", c.run.TaskName,
			"attempt", attempt,
			"result", "failed",
			"target", check.Message,
		)
	}

	c.run.LastValidationMessage = feedback.Message
}

func validationResult(success bool) string {
	if success {
		return "passed"
	}
	return "failed"
}

func normalizeUserAction(raw string) (source string, actionText string, target string, result string) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "unknown", "unknown_action", "", "unknown"
	}

	if event, ok := models.DecodeActionEvent(trimmed); ok {
		source = strings.TrimSpace(event.Source)
		if source == "" {
			source = "unknown"
		}
		actionText = strings.TrimSpace(event.Action)
		if actionText == "" {
			actionText = "unknown_action"
		}
		actionText = stripSourcePrefix(actionText, source)
		target = strings.TrimSpace(event.Target)
		result = strings.TrimSpace(event.Result)
		if result == "" {
			result = inferActionResult(actionText)
		}
		return source, actionText, target, result
	}

	result = "ok"
	lower := strings.ToLower(trimmed)

	if strings.HasPrefix(lower, "web request:") {
		rest := strings.TrimSpace(trimmed[len("web request:"):])
		parts := strings.Fields(rest)
		if len(parts) >= 2 {
			return "web", "request", strings.ToUpper(parts[0]) + " " + parts[1], "ok"
		}
		return "web", "request", rest, "ok"
	}

	if strings.HasPrefix(lower, "repl command:") {
		command := strings.TrimSpace(trimmed[len("repl command:"):])
		return "repl", "run_command", command, "ok"
	}

	words := strings.Fields(trimmed)
	if len(words) > 1 {
		head := strings.ToLower(words[0])
		if head == "tui" || head == "web" || head == "repl" {
			source = head
			actionText = strings.Join(words[1:], " ")
		} else {
			source = "unknown"
			actionText = trimmed
		}
	} else {
		source = "unknown"
		actionText = trimmed
	}

	result = inferActionResult(actionText)

	return source, actionText, "", result
}

func stripSourcePrefix(actionText string, source string) string {
	actionText = strings.TrimSpace(actionText)
	if actionText == "" || source == "" {
		return actionText
	}

	lowerSource := strings.ToLower(source)
	lowerAction := strings.ToLower(actionText)

	if strings.HasPrefix(lowerAction, lowerSource+" ") {
		return strings.TrimSpace(actionText[len(source):])
	}
	if strings.HasPrefix(lowerAction, lowerSource+"_") {
		return strings.TrimSpace(actionText[len(source)+1:])
	}

	return actionText
}

func inferActionResult(actionText string) string {
	lower := strings.ToLower(actionText)
	if strings.Contains(lower, "failed") {
		return "failed"
	}
	if strings.Contains(lower, "canceled") {
		return "canceled"
	}
	if strings.Contains(lower, "started") {
		return "started"
	}
	if strings.Contains(lower, "submitted") {
		return "submitted"
	}
	if strings.Contains(lower, "requested") {
		return "requested"
	}
	return "ok"
}

func normalizeAction(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	if lower == "" {
		return "unknown_action"
	}
	clean := nonWord.ReplaceAllString(lower, "_")
	clean = strings.Trim(clean, "_")
	if clean == "" {
		return "unknown_action"
	}
	return clean
}

func (c *taskRunCollector) finalize() models.TaskRunMetrics {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.run.EndedAt.IsZero() {
		c.run.EndedAt = time.Now()
	}
	c.run.DurationMs = c.run.EndedAt.Sub(c.run.StartedAt).Milliseconds()

	finalLogs := make([]models.TaskLogEntry, len(c.run.Logs))
	copy(finalLogs, c.run.Logs)

	final := c.run
	final.Logs = finalLogs

	return final
}

func appendTaskMetrics(path string, taskName string, run models.TaskRunMetrics, logger *slog.Logger) error {
	if path == "" {
		return nil
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create metrics directory: %w", err)
		}
	}

	metrics := models.TaskMetricsFile{
		TaskName: taskName,
		Runs:     []models.TaskRunMetrics{},
	}

	if bytes, err := os.ReadFile(path); err == nil {
		if len(bytes) > 0 {
			if err := json.Unmarshal(bytes, &metrics); err != nil {
				return fmt.Errorf("failed to parse metrics file %q: %w", path, err)
			}
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read metrics file %q: %w", path, err)
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
		return fmt.Errorf("failed to encode task metrics: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write metrics file %q: %w", path, err)
	}

	if logger != nil {
		logger.Info("task metrics persisted",
			"path", path,
			"task", taskName,
			"run_id", run.RunID,
			"total_runs", metrics.Summary.TotalRuns,
			"completed_runs", metrics.Summary.CompletedRuns,
		)
	}

	return nil
}

func (c *taskRunCollector) logWithLogger(level string, message string, attrs ...any) {
	if c.logger == nil {
		return
	}

	switch level {
	case "error":
		c.logger.Error(message, attrs...)
	case "warn":
		c.logger.Warn(message, attrs...)
	default:
		c.logger.Info(message, attrs...)
	}
}

func buildTaskStatsSummary(runs []models.TaskRunMetrics) models.TaskStatsSummary {
	summary := models.TaskStatsSummary{}
	if len(runs) == 0 {
		return summary
	}

	summary.TotalRuns = len(runs)
	summary.FirstRunStartedAt = runs[0].StartedAt

	for _, run := range runs {
		summary.TotalDurationMs += run.DurationMs
		summary.ValidationAttempts += run.ValidationAttempts
		summary.ValidationSuccesses += run.ValidationSuccesses
		summary.ValidationFailures += run.ValidationFailures
		summary.ValidationChecksPassed += run.ValidationChecksPassed
		summary.ValidationChecksFailed += run.ValidationChecksFailed

		if run.Completed {
			summary.CompletedRuns++
		} else {
			summary.IncompleteRuns++
		}

		if run.QuestionnaireCompleted {
			summary.QuestionnairesCompleted++
		}
		if run.QuestionnaireUserQuit {
			summary.QuestionnairesAbandoned++
		}

		if run.StartedAt.Before(summary.FirstRunStartedAt) {
			summary.FirstRunStartedAt = run.StartedAt
		}
		if run.StartedAt.After(summary.LastRunStartedAt) {
			summary.LastRunStartedAt = run.StartedAt
		}
		if run.EndedAt.After(summary.LastRunEndedAt) {
			summary.LastRunEndedAt = run.EndedAt
			summary.LastInterfaceType = run.InterfaceType
		}

		for _, log := range run.Logs {
			if log.Level == "user_action" {
				summary.TotalUserActions++
			}
		}
	}

	summary.AverageDurationMs = summary.TotalDurationMs / int64(summary.TotalRuns)
	return summary
}
