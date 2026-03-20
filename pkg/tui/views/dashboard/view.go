package dashboard

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
)

func (m *Model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		// if there is no space just print a loading message
		return tea.NewView("Loading...")
	}

	m.helpBar.SetWidth(m.width)

	header := m.header.View(m.width)
	headerHeight := m.header.Height()

	footer := components.RenderFooter(m.width, &m.helpBar, m.currentFeedback)
	footerHeight := lipgloss.Height(footer)

	// To avoid layer overflow or clipping, the label heights are calculated and subtracted from the available height before calculating the list heights to avoid layout overflow or clipping.
	contentHeight := m.height - headerHeight - footerHeight

	//mainLabel := styles.LabelStyle.Render("Display issues")
	//closedLabel := styles.LabelStyle.Render("Closed issues")
	//labelHeight := lipgloss.Height(mainLabel) + lipgloss.Height(closedLabel)
	availableForLists := contentHeight
	halfHeight := max(availableForLists / 2, 1)

	totalContentWidth := m.width - 1
	listWidth := totalContentWidth * styles.ListViewRatio / 100
	detailWidth := totalContentWidth - listWidth

	m.issueList.SetSize(listWidth, halfHeight)
	//m.closedIssueList.SetSize(listWidth, halfHeight)
	m.issueDetail.SetSize(detailWidth, contentHeight)

	// Only highlight the focused list; unfocused list should not show selection highlight.
	m.issueList.SetHighlightSelected(m.focusedWindow == 0 && m.focusedPaneMain == 0)
	m.closedIssueList.SetHighlightSelected(m.focusedWindow == 1 && m.focusedPaneClosed == 0)

	listView := m.issueList.View()
	//closedListView := m.closedIssueList.View()
	detailView := m.issueDetail.View()

	//if m.focusedWindow == 0 {
	//	mainLabel = lipgloss.NewStyle().Foreground(styles.Primary).Bold(true).Render("Display issues ▶")
	//} else {
	//	closedLabel = lipgloss.NewStyle().Foreground(styles.Primary).Bold(true).Render("Closed issues ▶")
	//}

	leftColumn := lipgloss.JoinVertical(lipgloss.Left,
		//mainLabel,
		listView,
		//closedLabel, closedListView,
	)
	content := lipgloss.JoinHorizontal(lipgloss.Left, leftColumn, detailView)

	mainView := lipgloss.JoinVertical(lipgloss.Left, header, content, footer)

	if m.editingAssignee {
		editBoxWidth := min(60, m.width-4)
		m.assigneeInput.SetWidth(editBoxWidth - 2)
		editContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Edit assignee (Enter to save, Esc to cancel):"),
			m.assigneeInput.View(),
		)
		editBox := styles.ContainerStyle.
			Width(editBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(editContent)
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, editBox))
	}

	if m.editingTitle {
		editBoxWidth := min(60, m.width-4)
		m.titleInput.SetWidth(editBoxWidth - 2)
		editContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Edit title (Enter to save, Esc to cancel):"),
			m.titleInput.View(),
		)
		editBox := styles.ContainerStyle.
			Width(editBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(editContent)
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, editBox))
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
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, editBox))
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
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, editBox))
	}

	if m.creatingIssue {
		createBoxWidth := min(60, m.width-4)
		m.createTitleInput.SetWidth(createBoxWidth - 2)
		createContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("New issue (Enter to create, Esc to cancel):"),
			m.createTitleInput.View(),
		)
		createBox := styles.ContainerStyle.
			Width(createBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(createContent)
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, createBox))
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
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, confirmBox))
	}

	if m.choosingStatus {
		statusContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Change status for "+m.statusIssueID+":"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("o = open   i = in_progress   r = ready_to_sprint   c = closing (choose reason)"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("Esc = cancel"),
		)
		statusBoxWidth := min(50, m.width-4)
		statusBox := styles.ContainerStyle.
			Width(statusBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(statusContent)
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, statusBox))
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
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, priorityBox))
	}

	if m.choosingCloseReason {
		reasonContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Choose closing reason for "+m.closeReasonIssueID+":"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("d = Done   u = Duplicate issue   w = Won't fix   o = Obsolete   h = Other"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("Esc = cancel"),
		)
		reasonBoxWidth := min(70, m.width-4)
		reasonBox := styles.ContainerStyle.
			Width(reasonBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(reasonContent)
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, reasonBox))
	}

	if m.closingOtherReason {
		editBoxWidth := min(60, m.width-4)
		m.closeReasonInput.SetWidth(editBoxWidth - 2)
		m.closeReasonInput.SetHeight(4)
		editContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.LabelStyle.Render("Enter closing reason for "+m.closeReasonIssueID+" (Enter or Ctrl+S to save, Esc to cancel):"),
			m.closeReasonInput.View(),
		)
		editBox := styles.ContainerStyle.
			Width(editBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(editContent)
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, editBox))
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
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, typeBox))
	}

	return tea.NewView(mainView)

}
