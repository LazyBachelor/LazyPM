package dashboard

import (
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	m.helpBar.SetWidth(m.width)

	header := m.header.View(m.width)
	headerHeight := m.header.Height()

	footer := m.footer()
	footerHeight := lipgloss.Height(footer)

	contentHeight := m.height - headerHeight - footerHeight

	totalContentWidth := m.width - 1
	listWidth := totalContentWidth * styles.ListViewRatio / 100
	detailWidth := totalContentWidth - listWidth

	m.issueList.SetSize(listWidth, contentHeight)
	m.issueDetail.SetSize(detailWidth, contentHeight)

	listView := m.issueList.View()
	detailView := m.issueDetail.View()

	content := lipgloss.JoinHorizontal(lipgloss.Left, listView, detailView)

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

func (m *Model) footer() string {
	feedbackStatus := m.currentFeedback.Message

	if feedbackStatus == "" {
		return m.helpBar.View()
	}

	m.helpBar.SetWidth(m.width - lipgloss.Width(feedbackStatus))
	return lipgloss.JoinHorizontal(lipgloss.Left, m.helpBar.View(), feedbackStatus)
}
