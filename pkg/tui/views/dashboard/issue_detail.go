package dashboard

import (
	"time"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type IssueDetail struct {
	viewport viewport.Model
	issue    models.Issue
	comments []*models.Comment
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

// SetComments updates the list of comments displayed for the current issue.
func (i *IssueDetail) SetComments(comments []*models.Comment) {
	i.comments = comments
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
		styles.LabelStyle.Render("Priority:") + styles.ValueStyle.Render(priorityCodeName(i.issue.Priority)),
	)

	assignee := i.issue.Assignee
	if assignee == "" {
		assignee = "—"
	}
	assigneeRow := styles.RowStyle.Render(
		styles.LabelStyle.Render("Assignee:") + styles.ValueStyle.Render(assignee),
	)

	descLabel := styles.LabelStyle.Render("Description:")
	descContent := styles.ValueStyle.Render(i.issue.Description)

	var parts []string
	parts = append(parts, titleRow, idRow, typeRow, statusRow, priorityRow, assigneeRow, descLabel, descContent)

	// Comments section
	commentsLabel := styles.LabelStyle.Render("Comments:")
	parts = append(parts, commentsLabel)
	if len(i.comments) == 0 {
		parts = append(parts, lipgloss.NewStyle().Foreground(styles.FaintText).Render("  No comments yet."))
	} else {
		for _, c := range i.comments {
			authorDate := lipgloss.NewStyle().Foreground(styles.Primary).Render(c.Author) + " " +
				lipgloss.NewStyle().Foreground(styles.FaintText).Render(formatCommentTime(c.CreatedAt))
			commentRow := lipgloss.JoinVertical(lipgloss.Left,
				authorDate,
				styles.ValueStyle.Render(c.Text),
			)
			parts = append(parts, commentRow)
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	i.viewport.SetContent(content)
}

func formatCommentTime(t time.Time) string {
	return t.Format("Jan 2, 15:04")
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
