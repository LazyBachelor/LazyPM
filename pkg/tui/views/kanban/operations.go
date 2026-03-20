package kanban

import (
	"context"
	"os"
	"os/user"
	"strconv"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/modal"
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
			m.setDetailIssueWithComments(*issue)
			targetStatus = issue.Status
			break
		}
	}

	switch targetStatus {
	case models.StatusOpen:
		m.focusManager.SetCurrent(modal.FocusColumn1)
	case models.StatusInProgress:
		m.focusManager.SetCurrent(modal.FocusColumn2)
	case models.StatusBlocked:
		m.focusManager.SetCurrent(modal.FocusColumn3)
	case models.StatusClosed:
		m.focusManager.SetCurrent(modal.FocusColumn4)
	}

	return tea.Sequence(todoCmd, inProgCmd, blockedCmd, doneCmd, func() tea.Msg {
		return msgs.SelectIssueMsg{IssueID: issueID}
	})
}

func (m *Model) refreshAndSubmit(issueID string) tea.Cmd {
	refreshCmd := m.refreshIssueListsAndSelectIssue(issueID)
	m.submitValidation()
	return refreshCmd
}

// Modal action handlers

func (m *Model) startEditTitle(selected ListIssue) tea.Cmd {
	m.currentIssueID = selected.ID
	titleModal := m.modalManager.GetTextInputModal(modal.ModalEditTitle)
	if titleModal != nil {
		titleModal.SetValue(selected.Issue.Title)
		titleModal.CursorEnd()
		return m.modalManager.ShowModal(modal.ModalEditTitle)
	}
	return nil
}

func (m *Model) startEditDescription(selected ListIssue) tea.Cmd {
	m.currentIssueID = selected.ID
	descModal := m.modalManager.GetTextAreaModal(modal.ModalEditDescription)
	if descModal != nil {
		descModal.SetValue(selected.Issue.Description)
		return m.modalManager.ShowModal(modal.ModalEditDescription)
	}
	return nil
}

func (m *Model) startCreateIssue() tea.Cmd {
	createModal := m.modalManager.GetTextInputModal(modal.ModalCreateIssue)
	if createModal != nil {
		createModal.Reset()
		return m.modalManager.ShowModal(modal.ModalCreateIssue)
	}
	return nil
}

func (m *Model) startConfirmDelete(issueID string, index int) tea.Cmd {
	m.currentIssueID = issueID
	m.deleteIndex = index
	deleteModal := m.modalManager.GetConfirmModal(modal.ModalConfirmDelete)
	if deleteModal != nil {
		return m.modalManager.ShowModal(modal.ModalConfirmDelete)
	}
	return nil
}

func (m *Model) startChooseStatus(selected ListIssue) tea.Cmd {
	m.currentIssueID = selected.ID
	return m.modalManager.ShowModal(modal.ModalSelectStatus)
}

func (m *Model) startChoosePriority(selected ListIssue) tea.Cmd {
	m.currentIssueID = selected.ID
	return m.modalManager.ShowModal(modal.ModalSelectPriority)
}

func (m *Model) startChooseType(selected ListIssue) tea.Cmd {
	m.currentIssueID = selected.ID
	return m.modalManager.ShowModal(modal.ModalSelectType)
}

func (m *Model) startEditAssignee(selected ListIssue) tea.Cmd {
	m.currentIssueID = selected.ID
	assigneeModal := m.modalManager.GetTextInputModal(modal.ModalEditAssignee)
	if assigneeModal != nil {
		assigneeModal.SetValue(selected.Assignee)
		assigneeModal.CursorEnd()
		return m.modalManager.ShowModal(modal.ModalEditAssignee)
	}
	return nil
}

func (m *Model) startAddComment(selected ListIssue) tea.Cmd {
	m.currentIssueID = selected.ID
	commentModal := m.modalManager.GetTextAreaModal(modal.ModalAddComment)
	if commentModal != nil {
		commentModal.Reset()
		return m.modalManager.ShowModal(modal.ModalAddComment)
	}
	return nil
}

// handleModalCompleted handles all modal completion messages
func (m *Model) handleModalCompleted(msg modal.ModalCompletedMsg) tea.Cmd {
	switch msg.ModalID {
	case modal.ModalEditTitle:
		if r, ok := msg.Value.(modal.TextInputResult); ok {
			cmd := msgs.UpdateIssueTitleCmd(m.app, m.currentIssueID, r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalCreateIssue:
		if r, ok := msg.Value.(modal.TextInputResult); ok && r.Value != "" {
			cmd := msgs.CreateIssueCmd(m.app, r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalEditAssignee:
		if r, ok := msg.Value.(modal.TextInputResult); ok {
			cmd := msgs.UpdateIssueAssigneeCmd(m.app, m.currentIssueID, r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalEditDescription:
		if r, ok := msg.Value.(modal.TextAreaResult); ok {
			cmd := msgs.UpdateIssueDescriptionCmd(m.app, m.currentIssueID, r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalConfirmDelete:
		if r, ok := msg.Value.(modal.ConfirmResult); ok && r.Confirmed {
			idx := m.deleteIndex
			issueID := m.currentIssueID
			m.deleteIndex = -1
			m.currentIssueID = ""
			cmd := msgs.DeleteIssueCmd(m.app, issueID, idx)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalSelectStatus:
		if r, ok := msg.Value.(modal.SelectResult); ok {
			if r.SelectedValue == "closing" {
				return m.modalManager.ShowModal(modal.ModalSelectCloseReason)
			}
			cmd := msgs.UpdateIssueStatusCmd(m.app, m.currentIssueID, r.SelectedValue)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalSelectCloseReason:
		if r, ok := msg.Value.(modal.SelectResult); ok {
			if r.SelectedValue == "other" {
				return m.modalManager.ShowModal(modal.ModalCloseReason)
			}
			cmd := msgs.CloseIssueCmd(m.app, m.currentIssueID, r.SelectedValue)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalSelectPriority:
		if r, ok := msg.Value.(modal.SelectResult); ok {
			priority, err := strconv.Atoi(r.SelectedValue)
			if err != nil {
				return nil
			}
			cmd := msgs.UpdateIssuePriorityCmd(m.app, m.currentIssueID, priority)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalSelectType:
		if r, ok := msg.Value.(modal.SelectResult); ok {
			issueType := models.IssueType(r.SelectedValue)
			cmd := msgs.UpdateIssueTypeCmd(m.app, m.currentIssueID, issueType)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalCloseReason:
		if r, ok := msg.Value.(modal.TextAreaResult); ok && r.Value != "" {
			cmd := msgs.CloseIssueCmd(m.app, m.currentIssueID, r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalAddComment:
		if r, ok := msg.Value.(modal.TextAreaResult); ok && r.Value != "" {
			cmd := msgs.AddIssueCommentCmd(m.app, m.currentIssueID, defaultCommentAuthor(), r.Value)
			return func() tea.Msg { return cmd() }
		}
	}
	return nil
}

// handleModalCancelled handles all modal cancellation messages
func (m *Model) handleModalCancelled(msg modal.ModalCancelledMsg) {
	switch msg.ModalID {
	case modal.ModalEditTitle:
		m.currentIssueID = ""
	case modal.ModalCreateIssue:
		// No cleanup needed
	case modal.ModalEditAssignee:
		m.currentIssueID = ""
	case modal.ModalEditDescription:
		m.currentIssueID = ""
	case modal.ModalConfirmDelete:
		m.deleteIndex = -1
		m.currentIssueID = ""
	case modal.ModalSelectStatus:
		m.currentIssueID = ""
	case modal.ModalSelectCloseReason:
		m.currentIssueID = ""
	case modal.ModalSelectPriority:
		m.currentIssueID = ""
	case modal.ModalSelectType:
		m.currentIssueID = ""
	case modal.ModalCloseReason:
		m.currentIssueID = ""
	case modal.ModalAddComment:
		m.currentIssueID = ""
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if cmd, handled := m.modalManager.Update(msg); handled {
		return m, cmd
	}

	switch msg := msg.(type) {
	case modal.ModalCompletedMsg:
		return m, m.handleModalCompleted(msg)

	case modal.ModalCancelledMsg:
		m.handleModalCancelled(msg)
		return m, nil
	case msgs.TitleUpdatedMsg:
		m.modalManager.GetTextInputModal(modal.ModalEditTitle).Reset()
		m.currentIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.DescriptionUpdatedMsg:
		m.modalManager.GetTextAreaModal(modal.ModalEditDescription).Reset()
		m.currentIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.StatusUpdatedMsg:
		m.currentIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.PriorityUpdatedMsg:
		m.currentIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.TypeUpdatedMsg:
		m.currentIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.AssigneeUpdatedMsg:
		m.modalManager.GetTextInputModal(modal.ModalEditAssignee).Reset()
		m.currentIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.IssueCommentAddedMsg:
		m.modalManager.GetTextAreaModal(modal.ModalAddComment).Reset()
		m.currentIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		m.submitValidation()
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case msgs.SelectIssueMsg:
		m.todoList.SelectIssueID(msg.IssueID)
		m.inProgList.SelectIssueID(msg.IssueID)
		m.blockedList.SelectIssueID(msg.IssueID)
		m.doneList.SelectIssueID(msg.IssueID)
		return m, nil

	case msgs.CreatedMsg:
		m.modalManager.GetTextInputModal(modal.ModalCreateIssue).Reset()
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

		selectedIssue := msg.Issue
		if selectedIssue.ID == "" {
			for _, issue := range allIssues {
				if issue.Title == msg.Issue.Title {
					selectedIssue = issue
					break
				}
			}
		}

		m.setDetailIssueWithComments(*selectedIssue)
		m.submitValidation()
		return m, tea.Sequence(todoCmd, inProgCmd, blockedCmd, doneCmd, func() tea.Msg {
			return msgs.SelectIssueMsg{IssueID: selectedIssue.ID}
		})

	case msgs.DeletedMsg:
		m.deleteIndex = -1
		m.currentIssueID = ""
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

		if len(todoIssues) == 0 && len(inProgIssues) == 0 && len(blockedIssues) == 0 && len(doneIssues) == 0 {
			m.setDetailIssueWithComments(models.Issue{})
			m.submitValidation()
			return m, tea.Sequence(todoCmd, inProgCmd, blockedCmd, doneCmd)
		}

		var targetIssues []*models.Issue
		switch m.focusManager.Current() {
		case modal.FocusColumn1:
			targetIssues = todoIssues
			if len(targetIssues) == 0 {
				if len(inProgIssues) > 0 {
					targetIssues = inProgIssues
					m.focusManager.SetCurrent(modal.FocusColumn2)
				} else if len(blockedIssues) > 0 {
					targetIssues = blockedIssues
					m.focusManager.SetCurrent(modal.FocusColumn3)
				} else if len(doneIssues) > 0 {
					targetIssues = doneIssues
					m.focusManager.SetCurrent(modal.FocusColumn4)
				}
			}
		case modal.FocusColumn2:
			targetIssues = inProgIssues
			if len(targetIssues) == 0 {
				if len(todoIssues) > 0 {
					targetIssues = todoIssues
					m.focusManager.SetCurrent(modal.FocusColumn1)
				} else if len(blockedIssues) > 0 {
					targetIssues = blockedIssues
					m.focusManager.SetCurrent(modal.FocusColumn3)
				} else if len(doneIssues) > 0 {
					targetIssues = doneIssues
					m.focusManager.SetCurrent(modal.FocusColumn4)
				}
			}
		case modal.FocusColumn3:
			targetIssues = blockedIssues
			if len(targetIssues) == 0 {
				if len(inProgIssues) > 0 {
					targetIssues = inProgIssues
					m.focusManager.SetCurrent(modal.FocusColumn2)
				} else if len(todoIssues) > 0 {
					targetIssues = todoIssues
					m.focusManager.SetCurrent(modal.FocusColumn1)
				} else if len(doneIssues) > 0 {
					targetIssues = doneIssues
					m.focusManager.SetCurrent(modal.FocusColumn4)
				}
			}
		case modal.FocusColumn4:
			targetIssues = doneIssues
			if len(targetIssues) == 0 {
				if len(blockedIssues) > 0 {
					targetIssues = blockedIssues
					m.focusManager.SetCurrent(modal.FocusColumn3)
				} else if len(inProgIssues) > 0 {
					targetIssues = inProgIssues
					m.focusManager.SetCurrent(modal.FocusColumn2)
				} else if len(todoIssues) > 0 {
					targetIssues = todoIssues
					m.focusManager.SetCurrent(modal.FocusColumn1)
				}
			}
		}

		if len(targetIssues) == 0 {
			m.setDetailIssueWithComments(models.Issue{})
			m.submitValidation()
			return m, tea.Sequence(todoCmd, inProgCmd, blockedCmd, doneCmd)
		}

		newIndex := msg.PreviousIndex
		if newIndex >= len(targetIssues) {
			newIndex = len(targetIssues) - 1
		}
		selectedIssue := targetIssues[newIndex]
		m.setDetailIssueWithComments(*selectedIssue)
		m.submitValidation()
		return m, tea.Sequence(todoCmd, inProgCmd, blockedCmd, doneCmd, func() tea.Msg {
			return msgs.SelectIssueMsg{IssueID: selectedIssue.ID}
		})

	case tea.KeyPressMsg:
		fl := m.FocusedIssueList()
		if fl.FilterState() == list.Filtering {
			cmd, _ := fl.Update(msg)
			return m, cmd
		}

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
		m.modalManager.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	fl := m.FocusedIssueList()
	cmd, changed := fl.Update(msg)
	if changed {
		if selected := fl.SelectedItem(); selected.ID != "" {
			m.setDetailIssueWithComments(selected.Issue)
		}
	}

	// Only propagate non-key messages to other lists (SetItems, etc.)
	// Key messages should only affect the focused list
	if _, isKeyMsg := msg.(tea.KeyPressMsg); !isKeyMsg {
		// Update all lists to ensure they receive commands like SetItems
		todoCmd, _ := m.todoList.Update(msg)
		inProgCmd, _ := m.inProgList.Update(msg)
		blockedCmd, _ := m.blockedList.Update(msg)
		doneCmd, _ := m.doneList.Update(msg)
		return m, tea.Sequence(cmd, todoCmd, inProgCmd, blockedCmd, doneCmd)
	}

	return m, cmd
}
