package dashboard

import (
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		// if there is no space just print a loading message
		return "Loading..."
	}

	m.helpBar.SetWidth(m.width)

	header := m.header.View(m.width)
	headerHeight := m.header.Height()

	footer := components.RenderFooter(m.width, &m.helpBar, m.currentFeedback)
	footerHeight := lipgloss.Height(footer)

	// To avoid layer overflow or clipping, the label heights are calculated and subtracted from the available height before calculating the list heights to avoid layout overflow or clipping.
	contentHeight := m.height - headerHeight - footerHeight

	mainLabel := styles.LabelStyle.Render("Display issues")
	closedLabel := styles.LabelStyle.Render("Closed issues")
	labelHeight := lipgloss.Height(mainLabel) + lipgloss.Height(closedLabel)
	availableForLists := contentHeight - labelHeight
	halfHeight := availableForLists / 2
	if halfHeight < 1 {
		halfHeight = 1
	}

	totalContentWidth := m.width - 1
	listWidth := totalContentWidth * styles.ListViewRatio / 100
	detailWidth := totalContentWidth - listWidth

	m.issueList.SetSize(listWidth, halfHeight)
	m.closedIssueList.SetSize(listWidth, halfHeight)
	m.issueDetail.SetSize(detailWidth, contentHeight)

	// Only highlight the focused list; unfocused list should not show selection highlight.
	m.issueList.SetHighlightSelected(m.focusedWindow == 0 && m.focusedPaneMain == 0)
	m.closedIssueList.SetHighlightSelected(m.focusedWindow == 1 && m.focusedPaneClosed == 0)

	listView := m.issueList.View()
	closedListView := m.closedIssueList.View()
	detailView := m.issueDetail.View()

	if m.focusedWindow == 0 {
		mainLabel = lipgloss.NewStyle().Foreground(styles.Primary).Bold(true).Render("Display issues ▶")
	} else {
		closedLabel = lipgloss.NewStyle().Foreground(styles.Primary).Bold(true).Render("Closed issues ▶")
	}
	leftColumn := lipgloss.JoinVertical(lipgloss.Left,
		mainLabel, listView,
		closedLabel, closedListView,
	)
	content := lipgloss.JoinHorizontal(lipgloss.Left, leftColumn, detailView)

	mainView := lipgloss.JoinVertical(lipgloss.Left, header, content, footer)

	if m.editingAssignee {
		editBoxWidth := min(60, m.width-4)
		m.assigneeInput.Width = editBoxWidth - 2
		editContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Edit assignee (Enter to save, Esc to cancel):"),
			m.assigneeInput.View(),
		)
		editBox := styles.ContainerStyle.
			Width(editBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(editContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, editBox)
	}

	if m.editingTitle {
		editBoxWidth := min(60, m.width-4)
		m.titleInput.Width = editBoxWidth - 2
		editContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Edit title (Enter to save, Esc to cancel):"),
			m.titleInput.View(),
		)
		editBox := styles.ContainerStyle.
			Width(editBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(editContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, editBox)
	}

	if m.addingComment {
		editBoxWidth := min(60, m.width-4)
		m.commentInput.SetWidth(editBoxWidth - 2)
		m.commentInput.SetHeight(8)
		editContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Add comment for "+m.commentIssueID+" (Ctrl+S or Enter to save, Esc to cancel):"),
			m.commentInput.View(),
		)
		editBox := styles.ContainerStyle.
			Width(editBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(editContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, editBox)
	}

	if m.editingDescription {
		editBoxWidth := min(60, m.width-4)
		m.descriptionInput.SetWidth(editBoxWidth - 2)
		m.descriptionInput.SetHeight(10)
		editContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Edit description (Ctrl+S to save, Esc to cancel):"),
			m.descriptionInput.View(),
		)
		editBox := styles.ContainerStyle.
			Width(editBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(editContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, editBox)
	}

	if m.creatingIssue {
		createBoxWidth := min(60, m.width-4)
		m.createTitleInput.Width = createBoxWidth - 2
		createContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("New issue (Enter to create, Esc to cancel):"),
			m.createTitleInput.View(),
		)
		createBox := styles.ContainerStyle.
			Width(createBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(createContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, createBox)
	}

	if m.confirmingDelete {
		confirmContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Delete issue "+m.deleteConfirmID+"?"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("Press y to delete, n or Esc to cancel"),
		)
		confirmBoxWidth := min(50, m.width-4)
		confirmBox := styles.ContainerStyle.
			Width(confirmBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(confirmContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, confirmBox)
	}

	if m.choosingStatus {
		statusContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Change status for "+m.statusIssueID+":"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("o = open   i = in_progress   b = blocked   r = ready_to_sprint   c = closed"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("Esc = cancel"),
		)
		statusBoxWidth := min(50, m.width-4)
		statusBox := styles.ContainerStyle.
			Width(statusBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(statusContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, statusBox)
	}

	if m.choosingPriority {
		priorityContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Change priority for "+m.priorityIssueID+":"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("0 = irrelevant 1 = low  2 = normal  3 = high  4 = critical"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("Esc = cancel"),
		)
		priorityBoxWidth := min(60, m.width-4)
		priorityBox := styles.ContainerStyle.
			Width(priorityBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(priorityContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, priorityBox)
	}

	if m.choosingType {
		typeContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Change type for "+m.typeIssueID+":"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("b = bug   f = feature   t = task   e = epic   c = chore"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("Esc = cancel"),
		)
		typeBoxWidth := min(65, m.width-4)
		typeBox := styles.ContainerStyle.
			Width(typeBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(typeContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, typeBox)
	}

	return mainView

}
