package issues

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// Msg types used by both dashboard and kanban TUI views.
type (
	TitleUpdatedMsg    struct{ IssueID string; Err error }
	DescriptionUpdatedMsg struct{ IssueID string; Err error }
	StatusUpdatedMsg   struct{ IssueID string; Err error }
	PriorityUpdatedMsg struct{ IssueID string; Err error }
	TypeUpdatedMsg     struct{ IssueID string; Err error }
	SelectIssueMsg     struct{ IssueID string }
	CreatedMsg         struct{ Issue *models.Issue; Err error }
	DeletedMsg         struct{ IssueID string; Err error; PreviousIndex int }
)

// UpdateIssueTitleCmd returns a command that updates an issue's title.
func UpdateIssueTitleCmd(app *app.App, issueID, newTitle string) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"title": newTitle}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return TitleUpdatedMsg{IssueID: issueID, Err: err}
	}
}

// UpdateIssueDescriptionCmd returns a command that updates an issue's description.
func UpdateIssueDescriptionCmd(app *app.App, issueID, newDescription string) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"description": newDescription}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return DescriptionUpdatedMsg{IssueID: issueID, Err: err}
	}
}

// UpdateIssueStatusCmd returns a command that updates an issue's status.
func UpdateIssueStatusCmd(app *app.App, issueID, status string) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"status": status}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return StatusUpdatedMsg{IssueID: issueID, Err: err}
	}
}

// UpdateIssuePriorityCmd returns a command that updates an issue's priority.
func UpdateIssuePriorityCmd(app *app.App, issueID string, priority int) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"priority": priority}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return PriorityUpdatedMsg{IssueID: issueID, Err: err}
	}
}

// UpdateIssueTypeCmd returns a command that updates an issue's type.
func UpdateIssueTypeCmd(app *app.App, issueID string, issueType models.IssueType) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"issue_type": string(issueType)}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return TypeUpdatedMsg{IssueID: issueID, Err: err}
	}
}

// CreateIssueCmd returns a command that creates a new issue.
func CreateIssueCmd(app *app.App, title string) tea.Cmd {
	return func() tea.Msg {
		issue := &models.Issue{
			Title:     title,
			Status:    models.StatusOpen,
			IssueType: models.TypeTask,
			Priority:  2,
		}
		err := app.Issues.CreateIssue(context.Background(), issue, "tui")
		return CreatedMsg{Issue: issue, Err: err}
	}
}

// DeleteIssueCmd returns a command that deletes an issue.
func DeleteIssueCmd(app *app.App, issueID string, currentIndex int) tea.Cmd {
	return func() tea.Msg {
		err := app.Issues.DeleteIssue(context.Background(), issueID)
		return DeletedMsg{IssueID: issueID, Err: err, PreviousIndex: currentIndex}
	}
}
