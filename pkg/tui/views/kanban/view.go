package kanban

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/modal"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
)

func (m *Model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Loading...")
	}

	m.helpBar.SetWidth(m.width)
	m.modalManager.SetSize(m.width, m.height)

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

	boardHeight := contentHeight / 2
	if boardHeight < 5 {
		boardHeight = contentHeight
	}

	m.todoList.SetSize(colWidth, boardHeight-1)
	m.inProgList.SetSize(colWidth, boardHeight-1)
	m.blockedList.SetSize(colWidth, boardHeight-1)
	m.doneList.SetSize(colWidth, boardHeight-1)

	// Only highlight the focused column's selected row
	currentFocus := m.focusManager.Current()
	m.todoList.SetHighlightSelected(currentFocus == modal.FocusColumn1)
	m.inProgList.SetHighlightSelected(currentFocus == modal.FocusColumn2)
	m.blockedList.SetHighlightSelected(currentFocus == modal.FocusColumn3)
	m.doneList.SetHighlightSelected(currentFocus == modal.FocusColumn4)

	m.issueDetail.SetSize(totalContentWidth, contentHeight-boardHeight)

	todoLabel := styles.LabelStyle.Render("To Do")
	inProgLabel := styles.LabelStyle.Render("In Progress")
	blockedLabel := styles.LabelStyle.Render("Blocked")
	doneLabel := styles.LabelStyle.Render("Done")

	highlight := lipgloss.NewStyle().Foreground(styles.Primary).Bold(true)
	switch currentFocus {
	case modal.FocusColumn1:
		todoLabel = highlight.Render("To Do ▶")
	case modal.FocusColumn2:
		inProgLabel = highlight.Render("In Progress ▶")
	case modal.FocusColumn3:
		blockedLabel = highlight.Render("Blocked ▶")
	case modal.FocusColumn4:
		doneLabel = highlight.Render("Done ▶")
	}

	todoCol := lipgloss.JoinVertical(lipgloss.Left, todoLabel, m.todoList.View())
	inProgCol := lipgloss.JoinVertical(lipgloss.Left, inProgLabel, m.inProgList.View())
	blockedCol := lipgloss.JoinVertical(lipgloss.Left, blockedLabel, m.blockedList.View())
	doneCol := lipgloss.JoinVertical(lipgloss.Left, doneLabel, m.doneList.View())

	board := lipgloss.JoinHorizontal(lipgloss.Left, todoCol, inProgCol, blockedCol, doneCol)
	content := lipgloss.JoinVertical(lipgloss.Left, board, m.issueDetail.View())

	mainView := lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
	mainViewHeight := lipgloss.Height(mainView)
	if mainViewHeight < m.height {
		spacerHeight := m.height - mainViewHeight
		spacer := lipgloss.NewStyle().Height(spacerHeight).Width(m.width).Render("")
		mainView = lipgloss.JoinVertical(lipgloss.Left, header, content, spacer, footer)
	}

	// Use modal manager to render active modal or main view
	return tea.NewView(m.modalManager.RenderWithMainView(mainView))
}

func (m *Model) footer() string {
	return components.RenderFooter(m.width, &m.helpBar, m.currentFeedback)
}
