package components

import (
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/style"
)

// RenderFooter renders the shared footer with the help bar and optional
// validation feedback message.
func RenderFooter(width int, helpBar *HelpBar, feedback models.ValidationFeedback) string {
	feedbackStatus := feedback.Message

	// Ensure the feedback message does not exceed the total available width.
	if feedback.Message != "" {
		styledFeedback := style.ErrorStyle.Render(feedbackStatus + " [Press '?' for details]")
		feedbackStatus = lipgloss.NewStyle().MaxWidth(width).Render(styledFeedback)

		if helpBar.IsExpanded() && feedbackStatus != "" {
			for _, check := range feedback.Checks {
				var prefix string
				if check.Valid {
					prefix = "✅ "
				} else {
					prefix = "❌ "
				}

				// Ensure each check line does not exceed the available width.
				remainingWidth := max(width-lipgloss.Width(prefix), 0)
				truncatedMsg := lipgloss.NewStyle().MaxWidth(remainingWidth).Render(check.Message)

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
