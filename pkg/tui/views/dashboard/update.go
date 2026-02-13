package dashboard

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type issueTitleUpdatedMsg struct {
	IssueID string
	Err     error
}

type selectIssueMsg struct {
	IssueID string
}

type issueCreatedMsg struct {
	Issue *models.Issue
	Err   error
}

type issueDeletedMsg struct {
	IssueID       string
	Err           error
	PreviousIndex int
}

func updateIssueTitleCmd(svc *service.Services, issueID, newTitle string) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"title": newTitle}
		err := svc.Beads.UpdateIssue(context.Background(), issueID, updates, "tui")
		return issueTitleUpdatedMsg{IssueID: issueID, Err: err}
	}
}

func createIssueCmd(svc *service.Services, title string) tea.Cmd {
	return func() tea.Msg {
		issue := &models.Issue{
			Title:     title,
			Status:    models.StatusOpen,
			IssueType: models.TypeTask,
		}
		err := svc.Beads.CreateIssue(context.Background(), issue, "tui")
		return issueCreatedMsg{Issue: issue, Err: err}
	}
}

func deleteIssueCmd(svc *service.Services, issueID string, currentIndex int) tea.Cmd {
	return func() tea.Msg {
		err := svc.Beads.DeleteIssue(context.Background(), issueID)
		return issueDeletedMsg{IssueID: issueID, Err: err, PreviousIndex: currentIndex}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case issueTitleUpdatedMsg:
		m.editingTitle = false
		m.editingIssueID = ""
		m.titleInput.Blur()
		if msg.Err != nil {
			return m, nil
		}
		// reload issues and keep selection on the updated issue
		issues, err := m.svc.Beads.AllIssues(context.Background())
		if err != nil {
			return m, nil
		}
		setItemsCmd := m.issueList.SetIssues(issues)
		for _, issue := range issues {
			if issue.ID == msg.IssueID {
				m.issueDetail.SetIssue(issue)
				break
			}
		}
		// after list is updated, select the edited issue again
		return m, tea.Sequence(setItemsCmd, func() tea.Msg { return selectIssueMsg{IssueID: msg.IssueID} })

	case selectIssueMsg:
		m.issueList.SelectIssueID(msg.IssueID)
		return m, nil

	case issueCreatedMsg:
		m.creatingIssue = false
		m.createTitleInput.Blur()
		m.createTitleInput.Reset()
		if msg.Err != nil || msg.Issue == nil {
			return m, nil
		}
		issues, err := m.svc.Beads.AllIssues(context.Background())
		if err != nil {
			return m, nil
		}
		setItemsCmd := m.issueList.SetIssues(issues)
		m.issueDetail.SetIssue(*msg.Issue)
		return m, tea.Sequence(setItemsCmd, func() tea.Msg { return selectIssueMsg{IssueID: msg.Issue.ID} })

	case issueDeletedMsg:
		m.confirmingDelete = false
		m.deleteConfirmID = ""
		if msg.Err != nil {
			return m, nil
		}
		issues, err := m.svc.Beads.AllIssues(context.Background())
		if err != nil {
			return m, nil
		}
		setItemsCmd := m.issueList.SetIssues(issues)
		if len(issues) == 0 {
			m.issueDetail.SetIssue(models.Issue{})
			return m, setItemsCmd
		}
		newIndex := msg.PreviousIndex
		if newIndex >= len(issues) {
			newIndex = len(issues) - 1
		}
		selectID := issues[newIndex].ID
		m.issueDetail.SetIssue(issues[newIndex])
		return m, tea.Sequence(setItemsCmd, func() tea.Msg { return selectIssueMsg{IssueID: selectID} })

	case tea.KeyMsg:
		if m.confirmingDelete {
			switch msg.String() {
			case "y", "Y":
				issueID := m.deleteConfirmID
				idx := m.deleteConfirmIndex
				m.confirmingDelete = false
				m.deleteConfirmID = ""
				return m, deleteIssueCmd(m.svc, issueID, idx)
			case "n", "N", "esc":
				m.confirmingDelete = false
				m.deleteConfirmID = ""
				return m, nil
			}
		}

		if m.creatingIssue {
			if msg.String() == "enter" {
				title := m.createTitleInput.Value()
				if title != "" {
					return m, createIssueCmd(m.svc, title)
				}
			}
			if msg.String() == "esc" {
				m.creatingIssue = false
				m.createTitleInput.Blur()
				m.createTitleInput.Reset()
				return m, nil
			}
			var cmd tea.Cmd
			m.createTitleInput, cmd = m.createTitleInput.Update(msg)
			return m, cmd
		}

		if m.editingTitle {
			if msg.String() == "enter" {
				newTitle := m.titleInput.Value()
				// dont update if its empty or same as the original
				if newTitle != "" {
					return m, updateIssueTitleCmd(m.svc, m.editingIssueID, newTitle)
				}
			}
			if msg.String() == "esc" {
				m.editingTitle = false
				m.editingIssueID = ""
				m.titleInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.titleInput, cmd = m.titleInput.Update(msg)
			return m, cmd
		}

		if m.issueList.FilterState() == list.Filtering {
			cmd, _ := m.issueList.Update(msg)
			return m, cmd
		}

		cmd := m.handleKeyMsg(msg)
		if cmd != nil {
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	cmd, changed := m.issueList.Update(msg)
	if changed {
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			m.issueDetail.SetIssue(selected.Issue)
		}
	}

	return m, cmd
}
