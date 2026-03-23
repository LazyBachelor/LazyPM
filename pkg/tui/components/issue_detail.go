package components

import (
	"time"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/style"
)

type IssueDetail struct {
	viewport viewport.Model
	issue    models.Issue
	comments []*models.Comment
	focused  bool
}

func NewIssueDetail() IssueDetail {
	vp := viewport.New(viewport.WithWidth(0), viewport.WithHeight(0))
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
	i.viewport = viewport.New(viewport.WithWidth(width), viewport.WithHeight(height))
	i.refreshContent()
}

func (i *IssueDetail) SetFocused(focused bool) {
	i.focused = focused
}

func (i *IssueDetail) refreshContent() {
	contentWidth := max(i.viewport.Width()-2, 1)

	titleRow := style.RowStyle.Render(
		style.TitleStyle.Render(i.issue.Title),
	)

	idRow := style.RowStyle.Render(
		style.LabelStyle.Render("ID:") + style.ValueStyle.Render(i.issue.ID),
	)

	typeRow := style.RowStyle.Render(
		style.LabelStyle.Render("Type:") + style.ValueStyle.Render(string(i.issue.IssueType)),
	)

	statusRow := style.RowStyle.Render(
		style.LabelStyle.Render("Status:") + style.StatusStyle(string(i.issue.Status)).Render(string(i.issue.Status)),
	)

	var closingReasonRow string
	if i.issue.Status == models.StatusClosed {
		var closingReason string
		if i.issue.CloseReason == "" {
			closingReason = "N/A"
		} else {
			closingReason = string(i.issue.CloseReason)
		}
		closingReasonRow = style.RowStyle.Render(
			style.LabelStyle.Render("Close reason: ") + style.ValueStyle.Render(closingReason))
	}

	priorityRow := style.RowStyle.Render(
		style.LabelStyle.Render("Priority:") + style.ValueStyle.Render(PriorityCodeName(i.issue.Priority)),
	)

	assigneeRow := style.RowStyle.Render(
		style.LabelStyle.Render("Assignee:") + style.ValueStyle.Render(i.issue.Assignee),
	)

	descLabel := style.LabelStyle.Render("Description:")
	descStyle := style.ValueStyle.Width(contentWidth)
	descContent := descStyle.Render(i.issue.Description)

	commentsLabel := style.LabelStyle.MarginTop(1).Render("Comments:")

	var parts []string
	parts = append(parts, titleRow, idRow, typeRow, statusRow, closingReasonRow, priorityRow, assigneeRow, descLabel, descContent, commentsLabel)
	parts = append(parts, i.renderComments()...)

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	i.viewport.SetContent(content)
}

func formatCommentTime(t time.Time) string {
	return t.Format("Jan 2, 15:04")
}

func (i IssueDetail) View() string {
	content := i.viewport.View()
	vpWidth := i.viewport.Width()
	vpHeight := i.viewport.Height()

	if i.focused {
		return style.DetailsContainerStyle.
			BorderForeground(style.PrimaryBorder).
			Width(vpWidth).
			Height(vpHeight).
			MaxHeight(vpHeight).
			Render(content)
	}
	return style.DetailsContainerStyle.
		Width(vpWidth).
		Height(vpHeight).
		MaxHeight(vpHeight).
		Render(content)
}

func (i *IssueDetail) ScrollUp(lines int) {
	i.viewport.ScrollUp(lines)
}

func (i *IssueDetail) ScrollDown(lines int) {
	i.viewport.ScrollDown(lines)
}

func (i *IssueDetail) renderComments() []string {
	contentWidth := max(i.viewport.Width()-4, 1)

	var parts []string
	if len(i.comments) == 0 {
		parts = append(parts, style.ValueStyle.Render("No comments yet."))
	} else {
		for _, c := range i.comments {
			authorDate := lipgloss.NewStyle().Foreground(style.Primary).Render(c.Author) + " " +
				lipgloss.NewStyle().Foreground(style.FaintText).Render(formatCommentTime(c.CreatedAt))
			commentTextStyle := style.ValueStyle.Width(contentWidth)
			commentRow := lipgloss.JoinVertical(lipgloss.Left,
				authorDate,
				commentTextStyle.MarginLeft(1).Render(c.Text),
			)
			parts = append(parts, commentRow)
		}
	}
	return parts
}
