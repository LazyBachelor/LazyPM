package dashboard

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type IssueDetail struct {
	viewport viewport.Model
	issue    models.Issue
	focused  bool
}

func NewIssueDetail() IssueDetail {
	vp := viewport.New(0, 0)
	return IssueDetail{
		viewport: vp,
	}
}

func (i *IssueDetail) SetIssue(issue models.Issue) {
	i.issue = issue
	i.refreshContent()
}

func (i *IssueDetail) SetSize(width, height int) {
	i.viewport.Height = height
	i.viewport.Width = width
	i.refreshContent()
}

func (i *IssueDetail) SetFocused(focused bool) {
	i.focused = focused
}

func (i *IssueDetail) refreshContent() {

	titleRow := styles.RowStyle.Render(
		styles.TitleStyle.Render(i.issue.Title),
	)

	idRow := styles.RowStyle.Render(
		styles.LabelStyle.Render("ID:") + styles.ValueStyle.Render(i.issue.ID),
	)

	typeRow := styles.RowStyle.Render(
		styles.LabelStyle.Render("Type:") + styles.ValueStyle.Render(string(i.issue.IssueType)),
	)

	statusRow := styles.RowStyle.Render(
		styles.LabelStyle.Render("Status:") + styles.StatusStyle(string(i.issue.Status)).Render(string(i.issue.Status)),
	)

	priorityRow := styles.RowStyle.Render(
		styles.LabelStyle.Render("Priority:") + styles.ValueStyle.Render(fmt.Sprintf("%d", i.issue.Priority)),
	)

	descLabel := styles.LabelStyle.Render("Description:")
	descContent := styles.ValueStyle.Render(i.issue.Description)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleRow,
		idRow,
		typeRow,
		statusRow,
		priorityRow,
		descLabel,
		descContent,
	)

	i.viewport.SetContent(content)
}

func (i IssueDetail) View() string {
	content := i.viewport.View()

	if i.focused {
		return styles.DetailsContainerStyle.
			BorderForeground(styles.PrimaryBorder).
			Width(i.viewport.Width).
			Height(i.viewport.Height).
			MaxHeight(i.viewport.Height).
			Render(content)
	}
	return styles.DetailsContainerStyle.
		Width(i.viewport.Width).
		Height(i.viewport.Height).
		MaxHeight(i.viewport.Height).
		Render(content)
}

func (i *IssueDetail) ScrollUp(lines int) {
	i.viewport.ScrollUp(lines)
}

func (i *IssueDetail) ScrollDown(lines int) {
	i.viewport.ScrollDown(lines)
}
