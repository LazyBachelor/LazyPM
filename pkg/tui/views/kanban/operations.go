package kanban

import (
	"context"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
)

// update handler for issueTitleUpdatedMsg, issueDescriptionUpdatedMsg, and issueStatusUpdatedMsg to avoid using nearly identical code for refreshing the issue lists and updating the detail view
// Fetch all msgs, update both lists, set the detail view for the given issue, and return a command to select that issue. Returns nil if fetch fails.
func (m *Model) refreshIssueListsAndSelectIssue(issueID string) tea.Cmd {
	allIssues, err := m.app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
	if err != nil {
		return nil
	}

	todoIssues := components.StatusOnly(allIssues, models.StatusOpen)
	inProgIssues := components.StatusOnly(allIssues, models.StatusInProgress)
	blockedIssues := components.StatusOnly(allIssues, models.StatusBlocked)
	doneIssues := components.StatusOnly(allIssues, models.StatusClosed)

	todoCmd := m.todoList.SetIssues(todoIssues)
	inProgCmd := m.inProgList.SetIssues(inProgIssues)
	blockedCmd := m.blockedList.SetIssues(blockedIssues)
	doneCmd := m.doneList.SetIssues(doneIssues)

	var targetStatus models.Status
	for _, issue := range allIssues {
		if issue.ID == issueID {
			m.issueDetail.SetIssue(*issue)
			targetStatus = issue.Status
			break
		}
	}

	switch targetStatus {
	case models.StatusOpen:
		m.focusedColumn = 0
	case models.StatusInProgress:
		m.focusedColumn = 1
	case models.StatusBlocked:
		m.focusedColumn = 2
	case models.StatusClosed:
		m.focusedColumn = 3
	}

	// Select the moved issue in its new column immediately so the highlight follows it.
	m.todoList.SelectIssueID(issueID)
	m.inProgList.SelectIssueID(issueID)
	m.blockedList.SelectIssueID(issueID)
	m.doneList.SelectIssueID(issueID)

	if m.submitChan != nil {
		select {
		case m.submitChan <- struct{}{}:
			m.logAction("tui submitted validation")
		default:
		}
	}

	return tea.Sequence(todoCmd, inProgCmd, blockedCmd, doneCmd)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgs.TitleUpdatedMsg:
		m.editingTitle = false
		m.editingIssueID = ""
		m.titleInput.Blur()
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case msgs.DescriptionUpdatedMsg:
		m.editingDescription = false
		m.editingDescIssueID = ""
		m.descriptionInput.Blur()
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case msgs.StatusUpdatedMsg:
		m.choosingStatus = false
		m.statusIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case msgs.PriorityUpdatedMsg:
		m.choosingPriority = false
		m.priorityIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case msgs.TypeUpdatedMsg:
		m.choosingType = false
		m.typeIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case msgs.AssigneeUpdatedMsg:
		m.editingAssignee = false
		m.assigneeIssueID = ""
		m.assigneeInput.Blur()
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case msgs.SelectIssueMsg:
		m.todoList.SelectIssueID(msg.IssueID)
		m.inProgList.SelectIssueID(msg.IssueID)
		m.blockedList.SelectIssueID(msg.IssueID)
		m.doneList.SelectIssueID(msg.IssueID)
		return m, nil

	case msgs.CreatedMsg:
		m.creatingIssue = false
		m.createTitleInput.Blur()
		m.createTitleInput.Reset()
		if msg.Err != nil || msg.Issue == nil {
			return m, nil
		}
		allIssues, err := m.app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
		if err != nil {
			return m, nil
		}

		todoIssues := components.StatusOnly(allIssues, models.StatusOpen)
		inProgIssues := components.StatusOnly(allIssues, models.StatusInProgress)
		blockedIssues := components.StatusOnly(allIssues, models.StatusBlocked)
		doneIssues := components.StatusOnly(allIssues, models.StatusClosed)

		todoCmd := m.todoList.SetIssues(todoIssues)
		inProgCmd := m.inProgList.SetIssues(inProgIssues)
		blockedCmd := m.blockedList.SetIssues(blockedIssues)
		doneCmd := m.doneList.SetIssues(doneIssues)

		// Determine the created issue from the refreshed list to ensure all fields (like ID) are populated.
		selectedIssue := msg.Issue
		if selectedIssue.ID == "" {
			for _, issue := range allIssues {
				// Prefer an issue that matches the created issue's title when ID is not yet known.
				if issue.Title == msg.Issue.Title {
					selectedIssue = issue
					break
				}
			}
		}

		m.issueDetail.SetIssue(*selectedIssue)
		return m, tea.Sequence(todoCmd, inProgCmd, blockedCmd, doneCmd, func() tea.Msg { return msgs.SelectIssueMsg{IssueID: selectedIssue.ID} })
	case msgs.DeletedMsg:
		m.confirmingDelete = false
		m.deleteConfirmID = ""
		if msg.Err != nil {
			return m, nil
		}
		allIssues, err := m.app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
		if err != nil {
			return m, nil
		}

		todoIssues := components.StatusOnly(allIssues, models.StatusOpen)
		inProgIssues := components.StatusOnly(allIssues, models.StatusInProgress)
		blockedIssues := components.StatusOnly(allIssues, models.StatusBlocked)
		doneIssues := components.StatusOnly(allIssues, models.StatusClosed)

		todoCmd := m.todoList.SetIssues(todoIssues)
		inProgCmd := m.inProgList.SetIssues(inProgIssues)
		blockedCmd := m.blockedList.SetIssues(blockedIssues)
		doneCmd := m.doneList.SetIssues(doneIssues)

		// If there are no msgs at all, clear the detail view and return.
		if len(todoIssues) == 0 && len(inProgIssues) == 0 && len(blockedIssues) == 0 && len(doneIssues) == 0 {
			m.issueDetail.SetIssue(models.Issue{})
			return m, tea.Sequence(todoCmd, inProgCmd, blockedCmd, doneCmd)
		}

		// Determine which column to use for the next selection based on the current focus.
		var targetIssues []*models.Issue
		switch m.focusedColumn {
		case 0:
			targetIssues = todoIssues
			if len(targetIssues) == 0 {
				if len(inProgIssues) > 0 {
					targetIssues = inProgIssues
					m.focusedColumn = 1
				} else if len(blockedIssues) > 0 {
					targetIssues = blockedIssues
					m.focusedColumn = 2
				} else if len(doneIssues) > 0 {
					targetIssues = doneIssues
					m.focusedColumn = 3
				}
			}
		case 1:
			targetIssues = inProgIssues
			if len(targetIssues) == 0 {
				if len(todoIssues) > 0 {
					targetIssues = todoIssues
					m.focusedColumn = 0
				} else if len(blockedIssues) > 0 {
					targetIssues = blockedIssues
					m.focusedColumn = 2
				} else if len(doneIssues) > 0 {
					targetIssues = doneIssues
					m.focusedColumn = 3
				}
			}
		case 2:
			targetIssues = blockedIssues
			if len(targetIssues) == 0 {
				if len(inProgIssues) > 0 {
					targetIssues = inProgIssues
					m.focusedColumn = 1
				} else if len(todoIssues) > 0 {
					targetIssues = todoIssues
					m.focusedColumn = 0
				} else if len(doneIssues) > 0 {
					targetIssues = doneIssues
					m.focusedColumn = 3
				}
			}
		case 3:
			targetIssues = doneIssues
			if len(targetIssues) == 0 {
				if len(blockedIssues) > 0 {
					targetIssues = blockedIssues
					m.focusedColumn = 2
				} else if len(inProgIssues) > 0 {
					targetIssues = inProgIssues
					m.focusedColumn = 1
				} else if len(todoIssues) > 0 {
					targetIssues = todoIssues
					m.focusedColumn = 0
				}
			}
		}

		// Safety: if targetIssues is still empty here, just clear detail and return.
		if len(targetIssues) == 0 {
			m.issueDetail.SetIssue(models.Issue{})
			return m, tea.Sequence(todoCmd, inProgCmd, blockedCmd, doneCmd)
		}

		newIndex := msg.PreviousIndex
		if newIndex >= len(targetIssues) {
			newIndex = len(targetIssues) - 1
		}
		selectedIssue := targetIssues[newIndex]
		m.issueDetail.SetIssue(*selectedIssue)
		return m, tea.Sequence(todoCmd, inProgCmd, blockedCmd, doneCmd, func() tea.Msg {
			return msgs.SelectIssueMsg{IssueID: selectedIssue.ID}
		})

	case tea.KeyPressMsg:
		if m.confirmingDelete {
			switch msg.String() {
			case "y", "Y":
				issueID := m.deleteConfirmID
				idx := m.deleteConfirmIndex
				m.confirmingDelete = false
				m.deleteConfirmID = ""
				return m, msgs.DeleteIssueCmd(m.app, issueID, idx)
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
				return m, msgs.UpdateIssueStatusCmd(m.app, issueID, string(models.StatusOpen))
			case "i":
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, msgs.UpdateIssueStatusCmd(m.app, issueID, string(models.StatusInProgress))
			case "b":
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, msgs.UpdateIssueStatusCmd(m.app, issueID, string(models.StatusBlocked))
			case "r":
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, msgs.UpdateIssueStatusCmd(m.app, issueID, string(models.StatusReadyToSprint))
			case "c":
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				m.choosingCloseReason = true
				m.closeReasonIssueID = issueID
				return m, nil
			case "esc":
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, nil
			}
		}

		if m.choosingCloseReason {
			var reason string
			switch msg.String() {
			case "d":
				reason = "Done"
			case "u":
				reason = "Duplicate issue"
			case "w":
				reason = "Won't fix"
			case "o":
				reason = "Obsolete"
			case "h":
				m.choosingCloseReason = false
				m.closingOtherReason = true
				m.closeReasonInput.SetValue("")
				m.closeReasonInput.Focus()
				return m, nil
			case "esc":
				m.choosingCloseReason = false
				m.closeReasonIssueID = ""
				return m, nil
			}

			if reason != "" {
				issueID := m.closeReasonIssueID
				m.choosingCloseReason = false
				m.closeReasonIssueID = ""
				return m, msgs.CloseIssueCmd(m.app, issueID, reason)
			}
		}

		if m.closingOtherReason {
			switch msg.String() {
			case "enter", "ctrl+s":
				reason := m.closeReasonInput.Value()
				if reason != "" {
					issueID := m.closeReasonIssueID
					m.closingOtherReason = false
					m.closeReasonIssueID = ""
					m.closeReasonInput.Blur()
					return m, msgs.CloseIssueCmd(m.app, issueID, reason)
				}
			case "esc":
				m.closingOtherReason = false
				m.closeReasonIssueID = ""
				m.closeReasonInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.closeReasonInput, cmd = m.closeReasonInput.Update(msg)
			return m, cmd
		}

		if m.choosingPriority {
			switch msg.String() {
			case "0", "1", "2", "3", "4":
				issueID := m.priorityIssueID
				priority := int(msg.String()[0] - '0')
				m.choosingPriority = false
				m.priorityIssueID = ""
				return m, msgs.UpdateIssuePriorityCmd(m.app, issueID, priority)
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
				return m, msgs.UpdateIssueTypeCmd(m.app, issueID, models.TypeBug)
			case "f":
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, msgs.UpdateIssueTypeCmd(m.app, issueID, models.TypeFeature)
			case "t":
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, msgs.UpdateIssueTypeCmd(m.app, issueID, models.TypeTask)
			case "e":
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, msgs.UpdateIssueTypeCmd(m.app, issueID, models.TypeEpic)
			case "c":
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, msgs.UpdateIssueTypeCmd(m.app, issueID, models.TypeChore)
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
					return m, msgs.CreateIssueCmd(m.app, title)
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

		if m.editingAssignee {
			if msg.String() == "enter" {
				assignee := m.assigneeInput.Value()
				return m, msgs.UpdateIssueAssigneeCmd(m.app, m.assigneeIssueID, assignee)
			}
			if msg.String() == "esc" {
				m.editingAssignee = false
				m.assigneeIssueID = ""
				m.assigneeInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.assigneeInput, cmd = m.assigneeInput.Update(msg)
			return m, cmd
		}

		if m.editingTitle {
			if msg.String() == "enter" {
				newTitle := m.titleInput.Value()
				if newTitle != "" {
					return m, msgs.UpdateIssueTitleCmd(m.app, m.editingIssueID, newTitle)
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
				return m, msgs.UpdateIssueDescriptionCmd(m.app, issueID, newDesc)
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

		cmd := m.handleKeyPressMsg(msg)
		if cmd != nil {
			return m, cmd
		}
	case components.ValidationFeedbackMsg:
		m.currentFeedback = msg.Feedback
		if msg.Feedback.Success {
			m.showComplete = true
			return m, tea.Quit
		}
		return m, components.ListenForValidation(m.feedbackChan)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	fl := m.FocusedIssueList()
	cmd, changed := fl.Update(msg)
	if changed {
		if selected := fl.SelectedItem(); selected.ID != "" {
			m.issueDetail.SetIssue(selected.Issue)
		}
	}
	return m, cmd
}
