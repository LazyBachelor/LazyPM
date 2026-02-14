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

	// return lipgloss.JoinVertical(lipgloss.Left, header, content, bottomView)
	mainView := lipgloss.JoinVertical(lipgloss.Left, header, content, bottomView)

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
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("o = open   i = in_progress   c = closed"),
			lipgloss.NewStyle().Foreground(styles.FaintText).Render("Esc = cancel"),
		)
		statusBoxWidth := min(50, m.width-4)
		statusBox := styles.ContainerStyle.
			Width(statusBoxWidth).
			BorderForeground(styles.PrimaryBorder).
			Render(statusContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, statusBox)
	}

	return mainView
}
