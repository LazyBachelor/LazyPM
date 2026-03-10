package components

import (
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// components contains reusable TUI modal renderers for issue actions.

func RenderEditTitle(width, height int, inputView string) string {
	editBoxWidth := min(60, width-4)
	editContent := lipgloss.JoinVertical(lipgloss.Left,
		styles.LabelStyle.Render("Edit title (Enter to save, Esc to cancel):"),
		inputView,
	)
	editBox := styles.ContainerStyle.
		Width(editBoxWidth).
		BorderForeground(styles.PrimaryBorder).
		Render(editContent)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, editBox)
}

func RenderEditDescription(width, height int, inputView string) string {
	editBoxWidth := min(60, width-4)
	editContent := lipgloss.JoinVertical(lipgloss.Left,
		styles.LabelStyle.Render("Edit description (Ctrl+S to save, Esc to cancel):"),
		inputView,
	)
	editBox := styles.ContainerStyle.
		Width(editBoxWidth).
		BorderForeground(styles.PrimaryBorder).
		Render(editContent)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, editBox)
}

func RenderCreateIssue(width, height int, inputView string) string {
	createBoxWidth := min(60, width-4)
	createContent := lipgloss.JoinVertical(lipgloss.Left,
		styles.LabelStyle.Render("New issue (Enter to create, Esc to cancel):"),
		inputView,
	)
	createBox := styles.ContainerStyle.
		Width(createBoxWidth).
		BorderForeground(styles.PrimaryBorder).
		Render(createContent)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, createBox)
}

func RenderConfirmDelete(width, height int, issueID string) string {
	confirmContent := lipgloss.JoinVertical(lipgloss.Left,
		styles.LabelStyle.Render("Delete issue "+issueID+"?"),
		lipgloss.NewStyle().Foreground(styles.FaintText).Render("Press y to delete, n or Esc to cancel"),
	)
	confirmBoxWidth := min(50, width-4)
	confirmBox := styles.ContainerStyle.
		Width(confirmBoxWidth).
		BorderForeground(styles.PrimaryBorder).
		Render(confirmContent)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, confirmBox)
}

func RenderChooseStatus(width, height int, issueID string) string {
	statusContent := lipgloss.JoinVertical(lipgloss.Left,
		styles.LabelStyle.Render("Change status for "+issueID+":"),
		lipgloss.NewStyle().Foreground(styles.FaintText).Render("o = open   i = in_progress   c = closed"),
		lipgloss.NewStyle().Foreground(styles.FaintText).Render("Esc = cancel"),
	)
	statusBoxWidth := min(50, width-4)
	statusBox := styles.ContainerStyle.
		Width(statusBoxWidth).
		BorderForeground(styles.PrimaryBorder).
		Render(statusContent)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, statusBox)
}

func RenderChoosePriority(width, height int, issueID string) string {
	priorityContent := lipgloss.JoinVertical(lipgloss.Left,
		styles.LabelStyle.Render("Change priority for "+issueID+":"),
		lipgloss.NewStyle().Foreground(styles.FaintText).Render("0 = irrelevant 1 = low  2 = normal  3 = high  4 = critical"),
		lipgloss.NewStyle().Foreground(styles.FaintText).Render("Esc = cancel"),
	)
	priorityBoxWidth := min(60, width-4)
	priorityBox := styles.ContainerStyle.
		Width(priorityBoxWidth).
		BorderForeground(styles.PrimaryBorder).
		Render(priorityContent)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, priorityBox)
}

func RenderChooseType(width, height int, issueID string) string {
	typeContent := lipgloss.JoinVertical(lipgloss.Left,
		styles.LabelStyle.Render("Change type for "+issueID+":"),
		lipgloss.NewStyle().Foreground(styles.FaintText).Render("b = bug   f = feature   t = task   e = epic   c = chore"),
		lipgloss.NewStyle().Foreground(styles.FaintText).Render("Esc = cancel"),
	)
	typeBoxWidth := min(65, width-4)
	typeBox := styles.ContainerStyle.
		Width(typeBoxWidth).
		BorderForeground(styles.PrimaryBorder).
		Render(typeContent)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, typeBox)
}

// RenderModals wraps the common modal overlay logic used by different views.
// It returns either one of the modal overlays (edit title/description, create,
// confirm delete, choose status/priority/type) or the provided main view if
// no modal is active.
func RenderModals(
	width, height int,
	editingTitle bool, titleInputView string,
	editingDescription bool, descriptionInputView string,
	creatingIssue bool, createTitleInputView string,
	confirmingDelete bool, deleteIssueID string,
	choosingStatus bool, statusIssueID string,
	choosingPriority bool, priorityIssueID string,
	choosingType bool, typeIssueID string,
	mainView string,
) string {
	if editingTitle {
		return RenderEditTitle(width, height, titleInputView)
	}

	if editingDescription {
		return RenderEditDescription(width, height, descriptionInputView)
	}

	if creatingIssue {
		return RenderCreateIssue(width, height, createTitleInputView)
	}

	if confirmingDelete {
		return RenderConfirmDelete(width, height, deleteIssueID)
	}

	if choosingStatus {
		return RenderChooseStatus(width, height, statusIssueID)
	}

	if choosingPriority {
		return RenderChoosePriority(width, height, priorityIssueID)
	}

	if choosingType {
		return RenderChooseType(width, height, typeIssueID)
	}

	return mainView
}

// RenderFooter renders the shared footer with the help bar and optional
// validation feedback message.

// truncateToWidth trims the given text so that its rendered width does not
// exceed maxWidth. If truncation occurs and there is room, an ellipsis is
// appended to indicate that the text was shortened.
func truncateToWidth(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	if lipgloss.Width(text) <= maxWidth {
		return text
	}

	const ellipsis = "…"
	ellipsisWidth := lipgloss.Width(ellipsis)
	if ellipsisWidth > maxWidth {
		// Not enough space even for an ellipsis; return empty.
		return ""
	}

	runes := []rune(text)
	current := ""
	for _, r := range runes {
		next := current + string(r)
		if lipgloss.Width(next)+ellipsisWidth > maxWidth {
			break
		}
		current = next
	}

	return current + ellipsis
}

func RenderFooter(width int, helpBar *HelpBar, feedback models.ValidationFeedback) string {
	feedbackStatus := feedback.Message

	// Ensure the feedback message does not exceed the total available width.
	feedbackStatus = truncateToWidth(feedbackStatus, width)

	if feedbackStatus == "" {
		return helpBar.View()
	}

	helpWidth := width - lipgloss.Width(feedbackStatus)
	if helpWidth < 0 {
		helpWidth = 0
	}

	helpBar.SetWidth(helpWidth)
	return lipgloss.JoinHorizontal(lipgloss.Left, helpBar.View(), feedbackStatus)
}

