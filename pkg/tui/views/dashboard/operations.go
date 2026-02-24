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

type issuePriorityUpdatedMsg struct {
	IssueID string
	Err     error
}

type issueTypeUpdatedMsg struct {
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

func updateIssueTitleCmd(app *service.App, issueID, newTitle string) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"title": newTitle}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return issueTitleUpdatedMsg{IssueID: issueID, Err: err}
	}
}

func updateIssueDescriptionCmd(app *service.App, issueID, newDescription string) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"description": newDescription}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return issueDescriptionUpdatedMsg{IssueID: issueID, Err: err}
	}
}

func updateIssueStatusCmd(app *service.App, issueID, status string) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"status": status}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return issueStatusUpdatedMsg{IssueID: issueID, Err: err}
	}
}

func updateIssuePriorityCmd(app *service.App, issueID string, priority int) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"priority": priority}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return issuePriorityUpdatedMsg{IssueID: issueID, Err: err}
	}
}

func updateIssueTypeCmd(app *service.App, issueID string, issueType models.IssueType) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"issue_type": string(issueType)}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return issueTypeUpdatedMsg{IssueID: issueID, Err: err}
	}
}

func updateIssuePriorityCmd(app *service.App, issueID string, priority int) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"priority": priority}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return issuePriorityUpdatedMsg{IssueID: issueID, Err: err}
	}
}

func updateIssueTypeCmd(app *service.App, issueID string, issueType models.IssueType) tea.Cmd {
	return func() tea.Msg {
		updates := map[string]interface{}{"issue_type": string(issueType)}
		err := app.Issues.UpdateIssue(context.Background(), issueID, updates, "tui")
		return issueTypeUpdatedMsg{IssueID: issueID, Err: err}
	}
}

func createIssueCmd(app *service.App, title string) tea.Cmd {
	return func() tea.Msg {
		issue := &models.Issue{
			Title:     title,
			Status:    models.StatusOpen,
			IssueType: models.TypeTask,
			Priority:  2,
		}
		err := app.Issues.CreateIssue(context.Background(), issue, "tui")
		return issueCreatedMsg{Issue: issue, Err: err}
	}
}

func deleteIssueCmd(app *service.App, issueID string, currentIndex int) tea.Cmd {
	return func() tea.Msg {
		err := app.Issues.DeleteIssue(context.Background(), issueID)
		return issueDeletedMsg{IssueID: issueID, Err: err, PreviousIndex: currentIndex}
	}
}

func (m *Model) refreshIssueListsAndSelectIssue(issueID string) tea.Cmd {
	/* update handler for issueTitleUpdatedMsg, issueDescriptionUpdatedMsg, and issueStatusUpdatedMsg to avoid using nearly identical code for refreshing the issue lists and updating the detail view
	Fetch all issues, update both lists, set the detail view for the given issue, and return a command to select that issue. Returns nil if fetch fails.
	*/
	issues, err := m.app.Issues.AllIssues(context.Background())
	if err != nil {
		return nil
	}
	setItemsCmd := m.issueList.SetIssues(OpenAndInProgressOnly(issues))
	closedSetCmd := m.closedIssueList.SetIssues(ClosedOnly(issues))
	for _, issue := range issues {
		if issue.ID == issueID {
			m.issueDetail.SetIssue(issue)
			break
		}
	}
	return tea.Sequence(setItemsCmd, closedSetCmd, func() tea.Msg { return selectIssueMsg{IssueID: issueID} })
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
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case issueDescriptionUpdatedMsg:
		m.editingDescription = false
		m.editingDescIssueID = ""
		m.descriptionInput.Blur()
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case issueStatusUpdatedMsg:
		m.choosingStatus = false
		m.statusIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case issuePriorityUpdatedMsg:
		m.choosingPriority = false
		m.priorityIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case issueTypeUpdatedMsg:
		m.choosingType = false
		m.typeIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

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
		issues, err := m.app.Issues.AllIssues(context.Background())
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
		issues, err := m.app.Issues.AllIssues(context.Background())
		if err != nil {
			return m, nil
		}
		openIssues := OpenAndInProgressOnly(issues)
		closedIssues := ClosedOnly(issues)
		setItemsCmd := m.issueList.SetIssues(openIssues)
		closedSetCmd := m.closedIssueList.SetIssues(closedIssues)
		// If there are no issues at all, clear the detail view and return.
		if len(openIssues) == 0 && len(closedIssues) == 0 {
			m.issueDetail.SetIssue(models.Issue{})
			return m, tea.Sequence(setItemsCmd, closedSetCmd)
		}

		// Determine which list to use for the next selection.
		var targetIssues []models.Issue
		if m.focusedWindow == 0 {
			targetIssues = openIssues
			if len(targetIssues) == 0 && len(closedIssues) > 0 {
				// The open list became empty; fall back to closed issues.
				targetIssues = closedIssues
				m.focusedWindow = 1
			}
		} else {
			targetIssues = closedIssues
			if len(targetIssues) == 0 && len(openIssues) > 0 {
				// The closed list became empty; fall back to open/in-progress issues.
				targetIssues = openIssues
				m.focusedWindow = 0
			}
		}

		// Safety: if targetIssues is still empty here, just clear detail and return.
		if len(targetIssues) == 0 {
			m.issueDetail.SetIssue(models.Issue{})
			return m, tea.Sequence(setItemsCmd, closedSetCmd)
		}

		newIndex := msg.PreviousIndex
		if newIndex >= len(targetIssues) {
			newIndex = len(targetIssues) - 1
		}
		selectedIssue := targetIssues[newIndex]
		m.issueDetail.SetIssue(selectedIssue)
		return m, tea.Sequence(setItemsCmd, closedSetCmd, func() tea.Msg {
			return selectIssueMsg{IssueID: selectedIssue.ID}
		})

	case tea.KeyMsg:
		if m.confirmingDelete {
			switch msg.String() {
			case "y", "Y":
				issueID := m.deleteConfirmID
				idx := m.deleteConfirmIndex
				m.confirmingDelete = false
				m.deleteConfirmID = ""
				return m, deleteIssueCmd(m.app, issueID, idx)
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
				return m, updateIssueStatusCmd(m.app, issueID, string(models.StatusOpen))
			case "i":
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, updateIssueStatusCmd(m.app, issueID, string(models.StatusInProgress))
			case "c":
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, updateIssueStatusCmd(m.app, issueID, string(models.StatusClosed))
			case "esc":
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, nil
			}
		}

		if m.choosingPriority {
			switch msg.String() {
			case "0", "1", "2", "3", "4":
				issueID := m.priorityIssueID
				priority := int(msg.String()[0] - '0')
				m.choosingPriority = false
				m.priorityIssueID = ""
				return m, updateIssuePriorityCmd(m.svc, issueID, priority)
			case "esc":
				m.choosingPriority = false
				m.priorityIssueID = ""
				return m, nil
			default:
				return m, nil
			}
		}

		if m.choosingType {
			switch msg.String() {
			case "b":
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, updateIssueTypeCmd(m.svc, issueID, models.TypeBug)
			case "f":
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, updateIssueTypeCmd(m.svc, issueID, models.TypeFeature)
			case "t":
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, updateIssueTypeCmd(m.svc, issueID, models.TypeTask)
			case "e":
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, updateIssueTypeCmd(m.svc, issueID, models.TypeEpic)
			case "c":
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, updateIssueTypeCmd(m.svc, issueID, models.TypeChore)
			case "esc":
				m.choosingType = false
				m.typeIssueID = ""
				return m, nil
			default:
				return m, nil
			}
		}

		if m.creatingIssue {
			if msg.String() == "enter" {
				title := m.createTitleInput.Value()
				if title != "" {
					return m, createIssueCmd(m.app, title)
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
					return m, updateIssueTitleCmd(m.app, m.editingIssueID, newTitle)
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
				return m, updateIssueDescriptionCmd(m.app, issueID, newDesc)
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

		// On main dashboard, ESC does nothing; only q quits; like in lazybeads.
		if msg.String() == "esc" {
			return m, nil
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
