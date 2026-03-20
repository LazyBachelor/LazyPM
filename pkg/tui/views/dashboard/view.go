package dashboard

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

	availableForLists := contentHeight
	halfHeight := max(availableForLists/2, 1)

	totalContentWidth := m.width - 1
	listWidth := totalContentWidth * style.ListViewRatio / 100
	detailWidth := totalContentWidth - listWidth

	m.issueList.SetSize(listWidth, halfHeight)
	m.issueDetail.SetSize(detailWidth, contentHeight)

	m.issueList.SetHighlightSelected(m.focusManager.IsFocused(modal.FocusList))

	listView := m.issueList.View()
	detailView := m.issueDetail.View()

	leftColumn := lipgloss.JoinVertical(lipgloss.Left, listView)
	content := lipgloss.JoinHorizontal(lipgloss.Left, leftColumn, detailView)

	mainView := lipgloss.JoinVertical(lipgloss.Left, header, content, footer)

	return tea.NewView(m.modalManager.RenderWithMainView(mainView))
}
