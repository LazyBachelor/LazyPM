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

type issueDescriptionUpdatedMsg struct {
	IssueID string
	Err     error
}

type issueStatusUpdatedMsg struct {
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

func updateIssueDescriptionCmd(svc *service.Services, issueID, newDescription string) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"description": newDescription}
		err := svc.Beads.UpdateIssue(context.Background(), issueID, updates, "tui")
		return issueDescriptionUpdatedMsg{IssueID: issueID, Err: err}
	}
}

func updateIssueStatusCmd(svc *service.Services, issueID, status string) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"status": status}
		err := svc.Beads.UpdateIssue(context.Background(), issueID, updates, "tui")
		return issueStatusUpdatedMsg{IssueID: issueID, Err: err}
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
		issues, err := m.svc.Beads.AllIssues(context.Background())
		if err != nil {
			return m, nil
		}
		setItemsCmd := m.issueList.SetIssues(OpenAndInProgressOnly(issues))
		closedSetCmd := m.closedIssueList.SetIssues(ClosedOnly(issues))
		for _, issue := range issues {
			if issue.ID == msg.IssueID {
				m.issueDetail.SetIssue(issue)
				break
			}
		}
		return m, tea.Sequence(setItemsCmd, closedSetCmd, func() tea.Msg { return selectIssueMsg{IssueID: msg.IssueID} })

	case issueDescriptionUpdatedMsg:
		m.editingDescription = false
		m.editingDescIssueID = ""
		m.descriptionInput.Blur()
		if msg.Err != nil {
			return m, nil
		}
		issues, err := m.svc.Beads.AllIssues(context.Background())
		if err != nil {
			return m, nil
		}
		setItemsCmd := m.issueList.SetIssues(OpenAndInProgressOnly(issues))
		closedSetCmd := m.closedIssueList.SetIssues(ClosedOnly(issues))
		for _, issue := range issues {
			if issue.ID == msg.IssueID {
				m.issueDetail.SetIssue(issue)
				break
			}
		}
		return m, tea.Sequence(setItemsCmd, closedSetCmd, func() tea.Msg { return selectIssueMsg{IssueID: msg.IssueID} })

	case issueStatusUpdatedMsg:
		m.choosingStatus = false
		m.statusIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		issues, err := m.svc.Beads.AllIssues(context.Background())
		if err != nil {
			return m, nil
		}
		setItemsCmd := m.issueList.SetIssues(OpenAndInProgressOnly(issues))
		closedSetCmd := m.closedIssueList.SetIssues(ClosedOnly(issues))
		for _, issue := range issues {
			if issue.ID == msg.IssueID {
				m.issueDetail.SetIssue(issue)
				break
			}
		}
		return m, tea.Sequence(setItemsCmd, closedSetCmd, func() tea.Msg { return selectIssueMsg{IssueID: msg.IssueID} })

	case selectIssueMsg:
		m.issueList.SelectIssueID(msg.IssueID)
		m.closedIssueList.SelectIssueID(msg.IssueID)
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
		setItemsCmd := m.issueList.SetIssues(OpenAndInProgressOnly(issues))
		closedSetCmd := m.closedIssueList.SetIssues(ClosedOnly(issues))

		// Determine the created issue from the refreshed list to ensure all fields (like ID) are populated.
		selectedIssue := msg.Issue
		if selectedIssue.ID == "" {
			for _, issue := range issues {
				// Prefer an issue that matches the created issue's title when ID is not yet known.
				if issue.Title == msg.Issue.Title {
					selectedIssue = &issue
					break
				}
			}
		}

		m.issueDetail.SetIssue(*selectedIssue)
		return m, tea.Sequence(setItemsCmd, closedSetCmd, func() tea.Msg { return selectIssueMsg{IssueID: selectedIssue.ID} })
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
		setItemsCmd := m.issueList.SetIssues(OpenAndInProgressOnly(issues))
		closedSetCmd := m.closedIssueList.SetIssues(ClosedOnly(issues))
		if len(issues) == 0 {
			m.issueDetail.SetIssue(models.Issue{})
			return m, tea.Sequence(setItemsCmd, closedSetCmd)
		}
		newIndex := msg.PreviousIndex
		if newIndex >= len(issues) {
			newIndex = len(issues) - 1
		}
		selectID := issues[newIndex].ID
		m.issueDetail.SetIssue(issues[newIndex])
		return m, tea.Sequence(setItemsCmd, closedSetCmd, func() tea.Msg { return selectIssueMsg{IssueID: selectID} })

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

		if m.choosingStatus {
			switch msg.String() {
			case "o":
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, updateIssueStatusCmd(m.svc, issueID, "open")
			case "i":
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, updateIssueStatusCmd(m.svc, issueID, "in_progress")
			case "c":
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, updateIssueStatusCmd(m.svc, issueID, "closed")
			case "esc":
				m.choosingStatus = false
				m.statusIssueID = ""
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

		if m.editingDescription {
			if msg.String() == "ctrl+s" {
				issueID := m.editingDescIssueID
				newDesc := m.descriptionInput.Value()
				m.editingDescription = false
				m.editingDescIssueID = ""
				m.descriptionInput.Blur()
				return m, updateIssueDescriptionCmd(m.svc, issueID, newDesc)
			}
			if msg.String() == "esc" {
				m.editingDescription = false
				m.editingDescIssueID = ""
				m.descriptionInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.descriptionInput, cmd = m.descriptionInput.Update(msg)
			return m, cmd
		}

		focusedList := m.FocusedIssueList()
		if focusedList.FilterState() == list.Filtering {
			cmd, _ := focusedList.Update(msg)
			return m, cmd
		}

		cmd := m.handleKeyMsg(msg)
		if cmd != nil {
			return m, cmd
		}
	case ValidationFeedbackMsg:
		m.currentFeedback = msg.Feedback
		if msg.Feedback.Success {
			m.showComplete = true
			return m, tea.Quit
		}
		return m, m.listenForValidation()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	if m.focusedWindow == 0 {
		cmd, changed := m.issueList.Update(msg)
		if changed {
			if selected := m.issueList.SelectedItem(); selected.ID != "" {
				m.issueDetail.SetIssue(selected.Issue)
			}
		}
		return m, cmd
	}
	cmd, changed := m.closedIssueList.Update(msg)
	if changed {
		if selected := m.closedIssueList.SelectedItem(); selected.ID != "" {
			m.issueDetail.SetIssue(selected.Issue)
		}
	}
	return m, cmd
}
