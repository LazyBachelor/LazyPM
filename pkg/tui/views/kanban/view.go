package kanban

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/modal"
	"github.com/LazyBachelor/LazyPM/internal/style"
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
	colWidth := max(totalContentWidth/4, 20)

	// Calculate initial list height (half of content height, minimum 5 rows)
	listHeight := contentHeight / 2
	if listHeight < 5 {
		listHeight = contentHeight
	}

	m.todoList.SetSize(colWidth, listHeight-1)
	m.inProgList.SetSize(colWidth, listHeight-1)
	m.blockedList.SetSize(colWidth, listHeight-1)
	m.doneList.SetSize(colWidth, listHeight-1)

	// Only highlight the focused column's selected row
	currentFocus := m.focusManager.Current()
	m.todoList.SetHighlightSelected(currentFocus == modal.FocusColumn1)
	m.inProgList.SetHighlightSelected(currentFocus == modal.FocusColumn2)
	m.blockedList.SetHighlightSelected(currentFocus == modal.FocusColumn3)
	m.doneList.SetHighlightSelected(currentFocus == modal.FocusColumn4)

	todoLabel := style.LabelStyle.Render("To Do")
	inProgLabel := style.LabelStyle.Render("In Progress")
	blockedLabel := style.LabelStyle.Render("Blocked")
	doneLabel := style.LabelStyle.Render("Done")

	highlight := lipgloss.NewStyle().Foreground(style.Primary).Bold(true)
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
	boardHeight := lipgloss.Height(board)

	detailHeight := max(contentHeight-boardHeight, 5)
	m.issueDetail.SetSize(totalContentWidth, detailHeight)

	content := lipgloss.JoinVertical(lipgloss.Left, board, m.issueDetail.View())

	mainView := lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
	mainViewHeight := lipgloss.Height(mainView)
	if mainViewHeight < m.height {
		spacerHeight := m.height - mainViewHeight
		spacer := lipgloss.NewStyle().Height(spacerHeight).Width(m.width).Render("")
		mainView = lipgloss.JoinVertical(lipgloss.Left, header, content, spacer, footer)
	}

	return tea.NewView(m.modalManager.RenderWithMainView(mainView))
}
