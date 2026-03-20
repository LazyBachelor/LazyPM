package dashboard

import (
	"context"
	"os"
	"os/user"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
)

func defaultCommentAuthor() string {
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	if s := os.Getenv("USER"); s != "" {
		return s
	}
	if s := os.Getenv("USERNAME"); s != "" {
		return s
	}
	return "user"
}

func (m *Model) refreshIssueListsAndSelectIssue(issueID string) tea.Cmd {
	/* update handler for issueTitleUpdatedMsg, issueDescriptionUpdatedMsg, and issueStatusUpdatedMsg to avoid using nearly identical code for refreshing the issue lists and updating the detail view
	Fetch all issues, update both lists, set the detail view for the given issue, and return a command to select that issue. Returns nil if fetch fails.
	*/
	allIssues, err := m.app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
	if err != nil {
		return nil
	}
	setItemsCmd := m.issueList.SetIssues(components.SortedIssues(allIssues))
	for _, issue := range allIssues {
		if issue.ID == issueID {
			m.setDetailIssueWithComments(*issue)
			break
		}
	}

	return tea.Sequence(setItemsCmd, func() tea.Msg { return msgs.SelectIssueMsg{IssueID: issueID} })
}

// refreshAndSubmit refreshes the issue lists and submits validation.
// This is a wrapper that should be used after any successful user action.
func (m *Model) refreshAndSubmit(issueID string) tea.Cmd {
	refreshCmd := m.refreshIssueListsAndSelectIssue(issueID)
	m.submitValidation()
	return refreshCmd
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgs.TitleUpdatedMsg:
		m.editingTitle = false
		m.editingIssueID = ""
		m.titleInput.Blur()
		if msg.Err != nil {
			m.logAction("tui failed to update issue title")
			return m, nil
		}
		m.logAction("tui updated issue title")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.DescriptionUpdatedMsg:
		m.editingDescription = false
		m.editingDescIssueID = ""
		m.descriptionInput.Blur()
		if msg.Err != nil {
			m.logAction("tui failed to update issue description")
			return m, nil
		}
		m.logAction("tui updated issue description")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.StatusUpdatedMsg:
		m.choosingStatus = false
		m.statusIssueID = ""
		if msg.Err != nil {
			m.logAction("tui failed to update issue status")
			return m, nil
		}
		m.logAction("tui updated issue status")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.PriorityUpdatedMsg:
		m.choosingPriority = false
		m.priorityIssueID = ""
		if msg.Err != nil {
			m.logAction("tui failed to update issue priority")
			return m, nil
		}
		m.logAction("tui updated issue priority")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.TypeUpdatedMsg:
		m.choosingType = false
		m.typeIssueID = ""
		if msg.Err != nil {
			m.logAction("tui failed to update issue type")
			return m, nil
		}
		m.logAction("tui updated issue type")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.AssigneeUpdatedMsg:
		m.editingAssignee = false
		m.assigneeIssueID = ""
		m.assigneeInput.Blur()
		if msg.Err != nil {
			m.logAction("tui failed to update issue assignee")
			return m, nil
		}
		m.logAction("tui updated issue assignee")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.SelectIssueMsg:
		m.issueList.SelectIssueID(msg.IssueID)
		m.closedIssueList.SelectIssueID(msg.IssueID)
		return m, nil

	case msgs.CreatedMsg:
		m.creatingIssue = false
		m.createTitleInput.Blur()
		m.createTitleInput.Reset()
		if msg.Err != nil || msg.Issue == nil {
			m.logAction("tui failed to create issue")
			return m, nil
		}
		allIssues, err := m.app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
		if err != nil {
			return m, nil
		}
		setItemsCmd := m.issueList.SetIssues(components.SortedIssues(allIssues))

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

		m.setDetailIssueWithComments(*selectedIssue)
		m.logAction("tui created issue")
		m.submitValidation()
		return m, tea.Sequence(setItemsCmd, func() tea.Msg { return msgs.SelectIssueMsg{IssueID: selectedIssue.ID} })
	case msgs.IssueCommentAddedMsg:
		m.addingComment = false
		m.commentIssueID = ""
		m.commentInput.Blur()
		m.commentInput.Reset()
		if msg.Err != nil {
			return m, nil
		}
		m.submitValidation()
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)
	case msgs.DeletedMsg:
		m.confirmingDelete = false
		m.deleteConfirmID = ""
		if msg.Err != nil {
			m.logAction("tui failed to delete issue")
			return m, nil
		}
		m.logAction("tui deleted issue")
		m.submitValidation()
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case tea.KeyPressMsg:
		if m.confirmingDelete {
			switch msg.String() {
			case "y", "Y":
				m.logAction("tui confirmed issue deletion")
				issueID := m.deleteConfirmID
				idx := m.deleteConfirmIndex
				m.confirmingDelete = false
				m.deleteConfirmID = ""
				return m, msgs.DeleteIssueCmd(m.app, issueID, idx)
			case "n", "N", "esc":
				m.logAction("tui canceled issue deletion")
				m.confirmingDelete = false
				m.deleteConfirmID = ""
				return m, nil
			}
		}

		if m.choosingStatus {
			switch msg.String() {
			case "o":
				m.logAction("tui selected issue status open")
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, msgs.UpdateIssueStatusCmd(m.app, issueID, string(models.StatusOpen))
			case "i":
				m.logAction("tui selected issue status in_progress")
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, msgs.UpdateIssueStatusCmd(m.app, issueID, string(models.StatusInProgress))
			case "b":
				m.logAction("tui selected issue status blocked")
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, msgs.UpdateIssueStatusCmd(m.app, issueID, string(models.StatusBlocked))
			case "r":
				m.logAction("tui selected issue status ready_to_sprint")
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, msgs.UpdateIssueStatusCmd(m.app, issueID, string(models.StatusReadyToSprint))
			case "c":
				m.logAction("tui selected issue status closing")
				issueID := m.statusIssueID
				m.choosingStatus = false
				m.statusIssueID = ""
				m.choosingCloseReason = true
				m.closeReasonIssueID = issueID
				return m, nil
			case "esc":
				m.logAction("tui canceled status picker")
				m.choosingStatus = false
				m.statusIssueID = ""
				return m, nil
			}
		}

		if m.choosingCloseReason {
			var reason string
			switch msg.String() {
			case "d":
				m.logAction("tui selected close reason done")
				reason = "Done"
			case "u":
				m.logAction("tui selected close reason duplicate issue")
				reason = "Duplicate issue"
			case "w":
				m.logAction("tui selected close reason won't fix")
				reason = "Won't fix"
			case "o":
				m.logAction("tui selected close reason obsolete")
				reason = "Obsolete"
			case "h":
				m.logAction("tui selected close reason other")
				m.choosingCloseReason = false
				m.closingOtherReason = true
				m.closeReasonInput.SetValue("")
				m.closeReasonInput.Focus()
				return m, nil
			case "esc":
				m.logAction("tui canceled close reason picker")
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
					m.logAction("tui submitted custom close reason")
					issueID := m.closeReasonIssueID
					m.closingOtherReason = false
					m.closeReasonIssueID = ""
					m.closeReasonInput.Blur()
					return m, msgs.CloseIssueCmd(m.app, issueID, reason)
				}
			case "esc":
				m.logAction("tui canceled custom close reason")
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
				m.logAction("tui selected issue priority")
				issueID := m.priorityIssueID
				priority := int(msg.String()[0] - '0')
				m.choosingPriority = false
				m.priorityIssueID = ""
				return m, msgs.UpdateIssuePriorityCmd(m.app, issueID, priority)
			case "esc":
				m.logAction("tui canceled priority picker")
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
				m.logAction("tui selected issue type bug")
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, msgs.UpdateIssueTypeCmd(m.app, issueID, models.TypeBug)
			case "f":
				m.logAction("tui selected issue type feature")
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, msgs.UpdateIssueTypeCmd(m.app, issueID, models.TypeFeature)
			case "t":
				m.logAction("tui selected issue type task")
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, msgs.UpdateIssueTypeCmd(m.app, issueID, models.TypeTask)
			case "e":
				m.logAction("tui selected issue type epic")
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, msgs.UpdateIssueTypeCmd(m.app, issueID, models.TypeEpic)
			case "c":
				m.logAction("tui selected issue type chore")
				issueID := m.typeIssueID
				m.choosingType = false
				m.typeIssueID = ""
				return m, msgs.UpdateIssueTypeCmd(m.app, issueID, models.TypeChore)
			case "esc":
				m.logAction("tui canceled type picker")
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
					m.logAction("tui submitted new issue")
					return m, msgs.CreateIssueCmd(m.app, title)
				}
			}
			if msg.String() == "esc" {
				m.logAction("tui canceled issue creation")
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
				m.logAction("tui submitted assignee edit")
				return m, msgs.UpdateIssueAssigneeCmd(m.app, m.assigneeIssueID, assignee)
			}
			if msg.String() == "esc" {
				m.logAction("tui canceled assignee edit")
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
					m.logAction("tui submitted issue title edit")
					return m, msgs.UpdateIssueTitleCmd(m.app, m.editingIssueID, newTitle)
				}
			}
			if msg.String() == "esc" {
				m.logAction("tui canceled issue title edit")
				m.editingTitle = false
				m.editingIssueID = ""
				m.titleInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.titleInput, cmd = m.titleInput.Update(msg)
			return m, cmd
		}

		if m.addingComment {
			if msg.String() == "ctrl+s" || msg.String() == "enter" {
				text := m.commentInput.Value()
				if text != "" {
					issueID := m.commentIssueID
					m.addingComment = false
					m.commentIssueID = ""
					m.commentInput.Blur()
					m.commentInput.Reset()
					return m, msgs.AddIssueCommentCmd(m.app, issueID, defaultCommentAuthor(), text)
				}
			}
			if msg.String() == "esc" {
				m.addingComment = false
				m.commentIssueID = ""
				m.commentInput.Blur()
				m.commentInput.Reset()
				return m, nil
			}
			var cmd tea.Cmd
			m.commentInput, cmd = m.commentInput.Update(msg)
			return m, cmd
		}

		if m.editingDescription {
			if msg.String() == "ctrl+s" {
				m.logAction("tui submitted issue description edit")
				issueID := m.editingDescIssueID
				newDesc := m.descriptionInput.Value()
				m.editingDescription = false
				m.editingDescIssueID = ""
				m.descriptionInput.Blur()
				return m, msgs.UpdateIssueDescriptionCmd(m.app, issueID, newDesc)
			}
			if msg.String() == "esc" {
				m.logAction("tui canceled issue description edit")
				m.editingDescription = false
				m.editingDescIssueID = ""
				m.descriptionInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.descriptionInput, cmd = m.descriptionInput.Update(msg)
			return m, cmd
		}

		if m.issueList.FilterState() == list.Filtering {
			cmd, _ := m.issueList.Update(msg)
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

	if m.focusedWindow == 0 {
		cmd, changed := m.issueList.Update(msg)
		if changed {
			if selected := m.issueList.SelectedItem(); selected.ID != "" {
				m.setDetailIssueWithComments(selected.Issue)
			}
		}
		return m, cmd
	}
	cmd, changed := m.closedIssueList.Update(msg)
	if changed {
		if selected := m.closedIssueList.SelectedItem(); selected.ID != "" {
			m.setDetailIssueWithComments(selected.Issue)
		}
	}
	return m, cmd
}
