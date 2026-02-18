package dashboard

import (
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	deletePromptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true)
	deletePrompt := deletePromptStyle.Render("Are you sure you want to delete this issue? (y/n)")

	if m.showDeleteConfirm {
		if m.width == 0 || m.height == 0 {
			return deletePrompt
		}
		m.helpBar.SetWidth(m.width)
		header := m.header.View(m.width)
		headerHeight := m.header.Height()
		bottomView := m.helpBar.View()
		bottomHeight := m.helpBar.Height()
		contentHeight := m.height - headerHeight - bottomHeight
		totalContentWidth := m.width - 1
		listWidth := totalContentWidth * styles.ListViewRatio / 100
		detailWidth := totalContentWidth - listWidth
		m.issueList.SetSize(listWidth, contentHeight)
		m.issueDetail.SetSize(detailWidth, contentHeight)
		listView := m.issueList.View()
		detailView := m.issueDetail.View()
		content := lipgloss.JoinHorizontal(lipgloss.Left, listView, detailView)
		mainView := lipgloss.JoinVertical(lipgloss.Left, header, content, bottomView)
		return lipgloss.JoinVertical(lipgloss.Left, deletePrompt, mainView)
	}

	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	m.helpBar.SetWidth(m.width)

	header := m.header.View(m.width)
	headerHeight := m.header.Height()

	bottomView := m.helpBar.View()
	bottomHeight := m.helpBar.Height()

	contentHeight := m.height - headerHeight - bottomHeight

	totalContentWidth := m.width - 1
	listWidth := totalContentWidth * styles.ListViewRatio / 100
	detailWidth := totalContentWidth - listWidth

	m.issueList.SetSize(listWidth, contentHeight)
	m.issueDetail.SetSize(detailWidth, contentHeight)

	listView := m.issueList.View()
	detailView := m.issueDetail.View()

	content := lipgloss.JoinHorizontal(lipgloss.Left, listView, detailView)

	mainView := lipgloss.JoinVertical(lipgloss.Left, header, content, bottomView)

	if m.deleteErr != nil {
		errMsg := deletePromptStyle.Render("Delete failed: " + m.deleteErr.Error())
		return lipgloss.JoinVertical(lipgloss.Left, errMsg, mainView)
	}

	return mainView
}
