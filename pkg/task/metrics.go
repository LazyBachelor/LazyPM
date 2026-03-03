package task

import (
	"fmt"
	"log/slog"
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

func (c *taskRunCollector) appendLog(entry models.TaskLogEntry) {
	entry.Timestamp = time.Now()

	c.run.Logs = append(c.run.Logs, entry)

	if c.logger == nil {
		return
	}

	attrs := []any{
		"task", c.run.TaskName,
		"interface", c.run.InterfaceType,
		"action", entry.Action,
		"result", entry.Result,
	}

	switch entry.Level {
	case "error":
		c.logger.Error(entry.Message, attrs...)
	case "warn":
		c.logger.Warn(entry.Message, attrs...)
	default:
		c.logger.Info(entry.Message, attrs...)
	}
}

func (c *taskRunCollector) log(level, message string) {
	result := "ok"
	switch level {
	case "error":
		result = "failed"
	case "warn":
		result = "warning"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.appendLog(models.TaskLogEntry{
		Level:   level,
		Message: message,
		Source:  "system",
		Action:  normalizeAction(message),
		Result:  result,
	})
}

func (c *taskRunCollector) recordUserAction(raw string) {
	source, actionText, target, result := normalizeUserAction(raw)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.appendLog(models.TaskLogEntry{
		Level:   "user_action",
		Message: raw,
		Source:  source,
		Action:  normalizeAction(actionText),
		Target:  target,
		Result:  result,
	})
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

	c.appendLog(models.TaskLogEntry{
		Level:   "questionnaire",
		Message: "questionnaire finished",
		Source:  "system",
		Action:  "questionnaire_finish",
		Result:  result,
	})

	if len(answers) == 0 {
		return
	}

	keys := make([]string, 0, len(answers))
	for k := range answers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := answers[k]
		if v == nil {
			continue
		}

		valueText := fmt.Sprintf("%v", v)

		c.appendLog(models.TaskLogEntry{
			Level:   "questionnaire",
			Message: fmt.Sprintf("questionnaire answer: %s=%s", k, valueText),
			Source:  "system",
			Action:  "questionnaire_answer",
			Target:  k,
			Result:  valueText,
		})
	}
}

func (c *taskRunCollector) recordValidation(feedback ValidationFeedback) {
	c.mu.Lock()

	c.run.ValidationAttempts++
	attempt := c.run.ValidationAttempts

	result := "failed"
	if feedback.Success {
		result = "passed"
		c.run.ValidationSuccesses++
	} else {
		c.run.ValidationFailures++
	}

	c.appendLog(models.TaskLogEntry{
		Level:   "validation",
		Message: fmt.Sprintf("validation attempt %d", attempt),
		Source:  "system",
		Action:  "validate_attempt",
		Result:  result,
		Attempt: attempt,
	})

	for _, check := range feedback.Checks {
		checkResult := "failed"
		if check.Valid {
			checkResult = "passed"
			c.run.ValidationChecksPassed++
		} else {
			c.run.ValidationChecksFailed++
		}

		c.appendLog(models.TaskLogEntry{
			Level:   "validation",
			Message: fmt.Sprintf("validation check: %s", check.Message),
			Source:  "system",
			Action:  "validate_check",
			Target:  check.Message,
			Result:  checkResult,
			Attempt: attempt,
		})
	}

	c.run.LastValidationMessage = feedback.Message
	c.mu.Unlock()
}

func (c *taskRunCollector) setCompleted(completed bool) {
	c.mu.Lock()
	c.run.Completed = completed
	c.mu.Unlock()
}

func (c *taskRunCollector) setError(err error) {
	if err == nil {
		return
	}
	c.mu.Lock()
	c.run.Error = err.Error()
	c.mu.Unlock()
}

func (c *taskRunCollector) finalize() models.TaskRunMetrics {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.run.EndedAt.IsZero() {
		c.run.EndedAt = time.Now()
	}

	c.run.DurationMs =
		c.run.EndedAt.Sub(c.run.StartedAt).Milliseconds()

	final := c.run
	final.Logs = append([]models.TaskLogEntry(nil), c.run.Logs...)
	return final
}

func normalizeUserAction(raw string) (source, actionText, target, result string) {

	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "unknown", "unknown_action", "", "unknown"
	}

	if event, ok := models.DecodeActionEvent(trimmed); ok {
		source = strings.TrimSpace(event.Source)
		if source == "" {
			source = "unknown"
		}

		actionText = stripSourcePrefix(
			strings.TrimSpace(event.Action),
			source,
		)

		target = strings.TrimSpace(event.Target)
		result = strings.TrimSpace(event.Result)

		if result == "" {
			result = inferActionResult(actionText)
		}
		return
	}

	lower := strings.ToLower(trimmed)

	if strings.HasPrefix(lower, "web request:") {
		rest := strings.TrimSpace(trimmed[len("web request:"):])
		return "web", "request", rest, "ok"
	}

	if strings.HasPrefix(lower, "repl command:") {
		cmd := strings.TrimSpace(trimmed[len("repl command:"):])
		return "repl", "run_command", cmd, "ok"
	}

	return "unknown", trimmed, "", inferActionResult(trimmed)
}

func stripSourcePrefix(actionText, source string) string {
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

	switch {
	case strings.Contains(lower, "failed"):
		return "failed"
	case strings.Contains(lower, "canceled"):
		return "canceled"
	case strings.Contains(lower, "started"):
		return "started"
	case strings.Contains(lower, "submitted"):
		return "submitted"
	case strings.Contains(lower, "requested"):
		return "requested"
	default:
		return "ok"
	}
}

func normalizeAction(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	if lower == "" {
		return "unknown_action"
	}

	clean := strings.Trim(
		nonWord.ReplaceAllString(lower, "_"),
		"_",
	)

	if clean == "" {
		return "unknown_action"
	}

	return clean
}
