package task

import (
	"github.com/LazyBachelor/LazyPM/internal/models"
)

func buildTaskStatsSummary(runs []models.TaskRunMetrics) models.TaskStatsSummary {
	var summary models.TaskStatsSummary

	if len(runs) == 0 {
		return summary
	}

	summary.TotalRuns = len(runs)

	for i, run := range runs {

		// First run handling
		if i == 0 || run.StartedAt.Before(summary.FirstRunStartedAt) {
			summary.FirstRunStartedAt = run.StartedAt
		}

		if run.StartedAt.After(summary.LastRunStartedAt) {
			summary.LastRunStartedAt = run.StartedAt
		}

		if run.EndedAt.After(summary.LastRunEndedAt) {
			summary.LastRunEndedAt = run.EndedAt
			summary.LastInterfaceType = run.InterfaceType
		}

		// Duration
		summary.TotalDurationMs += run.DurationMs

		// Completion
		if run.Completed {
			summary.CompletedRuns++
		} else {
			summary.IncompleteRuns++
		}

		// Validation
		summary.ValidationAttempts += run.ValidationAttempts
		summary.ValidationSuccesses += run.ValidationSuccesses
		summary.ValidationFailures += run.ValidationFailures
		summary.ValidationChecksPassed += run.ValidationChecksPassed
		summary.ValidationChecksFailed += run.ValidationChecksFailed

		// Questionnaire
		if run.QuestionnaireCompleted {
			summary.QuestionnairesCompleted++
		}
		if run.QuestionnaireUserQuit {
			summary.QuestionnairesAbandoned++
		}

		// User actions
		for _, log := range run.Logs {
			if log.Level == "user_action" {
				summary.TotalUserActions++
			}
		}
	}

	if summary.TotalRuns > 0 {
		summary.AverageDurationMs =
			summary.TotalDurationMs / int64(summary.TotalRuns)
	}

	return summary
}
