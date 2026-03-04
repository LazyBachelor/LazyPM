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

	return components.RenderModals(
		m.width,
		m.height,
		m.editingTitle,
		m.titleInput.View(),
		m.editingDescription,
		m.descriptionInput.View(),
		m.creatingIssue,
		m.createTitleInput.View(),
		m.confirmingDelete,
		m.deleteConfirmID,
		m.choosingStatus,
		m.statusIssueID,
		m.choosingPriority,
		m.priorityIssueID,
		m.choosingType,
		m.typeIssueID,
		mainView,
	)

}

func (m *Model) footer() string {
	// Kept for backwards compatibility; delegate to the shared helper.
	return components.RenderFooter(m.width, &m.helpBar, m.currentFeedback)
}
