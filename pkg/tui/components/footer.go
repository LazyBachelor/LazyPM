package components

import (
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/style"
	"github.com/LazyBachelor/LazyPM/internal/utils/truncate"
)

// RenderFooter renders the shared footer with the help bar and optional
// validation feedback message.
func RenderFooter(width int, helpBar *HelpBar, feedback models.ValidationFeedback) string {
	feedbackStatus := feedback.Message

	// Ensure the feedback message does not exceed the total available width.
	if feedback.Message != "" {
		// Allocate at least 30% of width for feedback, but not more than 60%
		feedbackWidth := max(width*3/10, min(width/2, 50))
		styledFeedback := style.ErrorStyle.Render(feedbackStatus + " [Press '?' for details]")
		feedbackStatus = truncate.TruncateToWidth(styledFeedback, feedbackWidth)

		if helpBar.IsExpanded() && feedbackStatus != "" {
			for _, check := range feedback.Checks {
				var prefix string
				if check.Valid {
					prefix = "✅ "
				} else {
					prefix = "❌ "
				}

				remainingWidth := max(width/2, 10)
				styledCheckMsg := style.TextStyle.Render(check.Message)
				truncatedMsg := truncate.TruncateToWidth(styledCheckMsg, remainingWidth)

				feedbackStatus += "\n" + prefix + truncatedMsg
			}
		}
	}

	if feedbackStatus == "" {
		return helpBar.View()
	}

	helpWidth := max(width-lipgloss.Width(feedbackStatus), 0)

	helpBar.SetWidth(helpWidth)
	return lipgloss.JoinHorizontal(lipgloss.Left, helpBar.View(), feedbackStatus)
}
