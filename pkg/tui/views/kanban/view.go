package kanban

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

	contentHeight := m.height - headerHeight - footerHeight
	totalContentWidth := m.width - 1
	colWidth := totalContentWidth / 4
	if colWidth < 20 {
		colWidth = 20
	}

	// Leave some space for the detail view below the board.
	boardHeight := contentHeight / 2
	if boardHeight < 5 {
		boardHeight = contentHeight
	}

	m.todoList.SetSize(colWidth, boardHeight-1)
	m.inProgList.SetSize(colWidth, boardHeight-1)
	m.blockedList.SetSize(colWidth, boardHeight-1)
	m.doneList.SetSize(colWidth, boardHeight-1)

	// Only highlight the selected row in the focused column.
	m.todoList.SetHighlightSelected(!m.focusOnDetail && m.focusedColumn == 0)
	m.inProgList.SetHighlightSelected(!m.focusOnDetail && m.focusedColumn == 1)
	m.blockedList.SetHighlightSelected(!m.focusOnDetail && m.focusedColumn == 2)
	m.doneList.SetHighlightSelected(!m.focusOnDetail && m.focusedColumn == 3)

	// Detail view takes full width below the board.
	m.issueDetail.SetSize(totalContentWidth, contentHeight-boardHeight)

	todoLabel := styles.LabelStyle.Render("To Do")
	inProgLabel := styles.LabelStyle.Render("In Progress")
	blockedLabel := styles.LabelStyle.Render("Blocked")
	doneLabel := styles.LabelStyle.Render("Done")

	highlight := lipgloss.NewStyle().Foreground(styles.Primary).Bold(true)
	switch m.focusedColumn {
	case 0:
		todoLabel = highlight.Render("To Do ▶")
	case 1:
		inProgLabel = highlight.Render("In Progress ▶")
	case 2:
		blockedLabel = highlight.Render("Blocked ▶")
	case 3:
		doneLabel = highlight.Render("Done ▶")
	}

	todoCol := lipgloss.JoinVertical(lipgloss.Left, todoLabel, m.todoList.View())
	inProgCol := lipgloss.JoinVertical(lipgloss.Left, inProgLabel, m.inProgList.View())
	blockedCol := lipgloss.JoinVertical(lipgloss.Left, blockedLabel, m.blockedList.View())
	doneCol := lipgloss.JoinVertical(lipgloss.Left, doneLabel, m.doneList.View())

	board := lipgloss.JoinHorizontal(lipgloss.Left, todoCol, inProgCol, blockedCol, doneCol)
	content := lipgloss.JoinVertical(lipgloss.Left, board, m.issueDetail.View())

	// Add spacer to lock footer to bottom of screen when content is shorter than available space
	// This is to avoid having the footer floating above the bottom of the screen
	mainView := lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
	mainViewHeight := lipgloss.Height(mainView)
	if mainViewHeight < m.height {
		spacerHeight := m.height - mainViewHeight
		spacer := lipgloss.NewStyle().Height(spacerHeight).Width(m.width).Render("")
		mainView = lipgloss.JoinVertical(lipgloss.Left, header, content, spacer, footer)
	}

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
		m.editingAssignee,
		m.assigneeInput.View(),
		mainView,
	)

}

func (m *Model) footer() string {
	// Kept for backwards compatibility; delegate to the shared helper.
	return components.RenderFooter(m.width, &m.helpBar, m.currentFeedback)
}
