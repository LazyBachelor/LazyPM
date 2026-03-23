package kanban

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/style"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/modal"
)

func (m *Model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Loading...")
	}

	m.helpBar.SetWidth(m.width)
	m.modalManager.SetSize(m.width, m.height)

	var sprintName string
	if m.currentSprintNum > 0 {
		sprintName = fmt.Sprintf("Sprint %d", m.currentSprintNum)
	} else {
		sprintName = "No Sprint"
	}

	m.header = components.NewHeader(fmt.Sprintf("Kanban Board - %s", sprintName))
	header := m.header.View(m.width)
	headerHeight := m.header.Height()

	footer := components.RenderFooter(m.width, &m.helpBar, m.currentFeedback)
	footerHeight := lipgloss.Height(footer)

	contentHeight := m.height - headerHeight - footerHeight
	totalContentWidth := m.width - 1
	colWidth := max(totalContentWidth/5, 16)

	listHeight := contentHeight / 2
	if listHeight < 5 {
		listHeight = contentHeight
	}

	m.backlogList.SetSize(colWidth, listHeight-1)
	m.todoList.SetSize(colWidth, listHeight-1)
	m.inProgList.SetSize(colWidth, listHeight-1)
	m.blockedList.SetSize(colWidth, listHeight-1)
	m.doneList.SetSize(colWidth, listHeight-1)

	currentFocus := m.focusManager.Current()
	m.backlogList.SetHighlightSelected(currentFocus == modal.FocusColumn0)
	m.todoList.SetHighlightSelected(currentFocus == modal.FocusColumn1)
	m.inProgList.SetHighlightSelected(currentFocus == modal.FocusColumn2)
	m.blockedList.SetHighlightSelected(currentFocus == modal.FocusColumn3)
	m.doneList.SetHighlightSelected(currentFocus == modal.FocusColumn4)

	backlogLabel := style.LabelStyle.Render("Backlog")
	todoLabel := style.LabelStyle.Render(fmt.Sprintf("%s - To Do", sprintName))
	inProgLabel := style.LabelStyle.Render(fmt.Sprintf("%s - In Progress", sprintName))
	blockedLabel := style.LabelStyle.Render(fmt.Sprintf("%s - Blocked", sprintName))
	doneLabel := style.LabelStyle.Render(fmt.Sprintf("%s - Done", sprintName))

	highlight := lipgloss.NewStyle().Foreground(style.Primary).Bold(true)
	switch currentFocus {
	case modal.FocusColumn0:
		backlogLabel = highlight.Render("Backlog ▶")
	case modal.FocusColumn1:
		todoLabel = highlight.Render(fmt.Sprintf("%s - To Do ▶", sprintName))
	case modal.FocusColumn2:
		inProgLabel = highlight.Render(fmt.Sprintf("%s - In Progress ▶", sprintName))
	case modal.FocusColumn3:
		blockedLabel = highlight.Render(fmt.Sprintf("%s - Blocked ▶", sprintName))
	case modal.FocusColumn4:
		doneLabel = highlight.Render(fmt.Sprintf("%s - Done ▶", sprintName))
	}

	backlogCol := lipgloss.JoinVertical(lipgloss.Left, backlogLabel, m.backlogList.View())
	todoCol := lipgloss.JoinVertical(lipgloss.Left, todoLabel, m.todoList.View())
	inProgCol := lipgloss.JoinVertical(lipgloss.Left, inProgLabel, m.inProgList.View())
	blockedCol := lipgloss.JoinVertical(lipgloss.Left, blockedLabel, m.blockedList.View())
	doneCol := lipgloss.JoinVertical(lipgloss.Left, doneLabel, m.doneList.View())

	board := lipgloss.JoinHorizontal(lipgloss.Left, backlogCol, todoCol, inProgCol, blockedCol, doneCol)
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
