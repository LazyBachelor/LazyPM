package dashboard

import (
	"context"
	"strconv"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/utils/user"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/modal"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
)

func (m *Model) refreshIssueListsAndSelectIssue(issueID string) tea.Cmd {
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

func (m *Model) refreshAndSubmit(issueID string) tea.Cmd {
	refreshCmd := m.refreshIssueListsAndSelectIssue(issueID)
	m.submitValidation()
	return refreshCmd
}

func (m *Model) startConfirmExit() tea.Cmd {
	return m.modalManager.ShowModal(modal.ModalConfirmExit)
}

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

func (m *Model) startManageDependencies(selected ListIssue) tea.Cmd {
	m.currentIssueID = selected.ID
	depModal := modal.NewDependenciesModal(m.app, selected.ID)
	m.dependenciesModal = depModal
	return m.modalManager.PushModal(depModal)
}

// handleModalCompleted handles all modal completion messages
func (m *Model) handleModalCompleted(msg modal.ModalCompletedMsg) tea.Cmd {
	switch msg.ModalID {
	case modal.ModalConfirmExit:
		if r, ok := msg.Value.(modal.ConfirmResult); ok && r.Confirmed {
			m.logAction("tui confirmed exit")
			return tea.Quit
		}
	case modal.ModalEditTitle:
		if r, ok := msg.Value.(modal.TextInputResult); ok {
			m.logAction("tui submitted issue title edit")
			cmd := msgs.UpdateIssueTitleCmd(m.app, m.currentIssueID, r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalCreateIssue:
		if r, ok := msg.Value.(modal.TextInputResult); ok && r.Value != "" {
			m.logAction("tui submitted new issue")
			cmd := msgs.CreateIssueCmd(m.app, r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalEditAssignee:
		if r, ok := msg.Value.(modal.TextInputResult); ok {
			m.logAction("tui submitted assignee edit")
			cmd := msgs.UpdateIssueAssigneeCmd(m.app, m.currentIssueID, r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalEditDescription:
		if r, ok := msg.Value.(modal.TextAreaResult); ok {
			m.logAction("tui submitted issue description edit")
			cmd := msgs.UpdateIssueDescriptionCmd(m.app, m.currentIssueID, r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalAddComment:
		if r, ok := msg.Value.(modal.TextAreaResult); ok && r.Value != "" {
			cmd := msgs.AddIssueCommentCmd(m.app, m.currentIssueID, user.GetOsUsername(), r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalCloseReason:
		if r, ok := msg.Value.(modal.TextAreaResult); ok && r.Value != "" {
			m.logAction("tui submitted custom close reason")
			cmd := msgs.CloseIssueCmd(m.app, m.currentIssueID, r.Value)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalConfirmDelete:
		if r, ok := msg.Value.(modal.ConfirmResult); ok && r.Confirmed {
			m.logAction("tui confirmed issue deletion")
			idx := m.deleteIndex
			issueID := m.currentIssueID
			m.deleteIndex = -1
			m.currentIssueID = ""
			cmd := msgs.DeleteIssueCmd(m.app, issueID, idx)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalSelectStatus:
		if r, ok := msg.Value.(modal.SelectResult); ok {
			m.logAction("tui selected issue status")
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
			m.logAction("tui selected issue priority")
			priority, _ := strconv.Atoi(r.SelectedValue)
			cmd := msgs.UpdateIssuePriorityCmd(m.app, m.currentIssueID, priority)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalSelectType:
		if r, ok := msg.Value.(modal.SelectResult); ok {
			m.logAction("tui selected issue type")
			issueType := models.IssueType(r.SelectedValue)
			cmd := msgs.UpdateIssueTypeCmd(m.app, m.currentIssueID, issueType)
			return func() tea.Msg { return cmd() }
		}
	case modal.ModalSelectDependency:
		if r, ok := msg.Value.(modal.SelectResult); ok && r.SelectedValue != "" {
			m.logAction("tui added dependency")
			m.modalManager.PopModal()
			return msgs.AddDependencyCmd(m.app, m.currentIssueID, r.SelectedValue)
		}
	case modal.ModalSelectRemoveDependency:
		if r, ok := msg.Value.(modal.SelectResult); ok && r.SelectedValue != "" {
			m.logAction("tui removed dependency")
			m.modalManager.PopModal()
			return msgs.RemoveDependencyCmd(m.app, m.currentIssueID, r.SelectedValue)
		}
	}
	return nil
}

// handleModalCancelled handles all modal cancellation messages
func (m *Model) handleModalCancelled(msg modal.ModalCancelledMsg) {
	switch msg.ModalID {
	case modal.ModalConfirmExit:
		m.logAction("tui canceled exit")
	case modal.ModalEditTitle:
		m.currentIssueID = ""
		m.logAction("tui canceled issue title edit")
	case modal.ModalCreateIssue:
		m.logAction("tui canceled issue creation")
	case modal.ModalEditAssignee:
		m.currentIssueID = ""
		m.logAction("tui canceled assignee edit")
	case modal.ModalEditDescription:
		m.currentIssueID = ""
		m.logAction("tui canceled issue description edit")
	case modal.ModalAddComment:
		m.currentIssueID = ""
	case modal.ModalCloseReason:
		m.currentIssueID = ""
		m.logAction("tui canceled close reason")
	case modal.ModalConfirmDelete:
		m.deleteIndex = -1
		m.currentIssueID = ""
		m.logAction("tui canceled issue deletion")
	case modal.ModalSelectStatus:
		m.currentIssueID = ""
		m.logAction("tui canceled status picker")
	case modal.ModalSelectCloseReason:
		m.currentIssueID = ""
		m.logAction("tui canceled close reason")
	case modal.ModalSelectPriority:
		m.currentIssueID = ""
		m.logAction("tui canceled priority picker")
	case modal.ModalSelectType:
		m.currentIssueID = ""
		m.logAction("tui canceled type picker")
	case modal.ModalManageDependencies:
		m.dependenciesModal = nil
		m.modalManager.PopModal()
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			m.setDetailIssueWithComments(selected.Issue)
		}
		m.logAction("tui closed dependencies management")
	case modal.ModalSelectDependency:
		m.modalManager.PopModal()
		m.logAction("tui canceled add dependency")
	case modal.ModalSelectRemoveDependency:
		m.modalManager.PopModal()
		m.logAction("tui canceled remove dependency")
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

	case msgs.DependencyAddRequestedMsg:
		issues := modal.EligibleDependencyIssues(m.app, msg.IssueID)
		if len(issues) == 0 {
			m.logAction("tui no issues available to add as dependency")
			return m, nil
		}
		listModal := modal.NewDependencyListModal(issues, m.width, m.height)
		return m, m.modalManager.PushModal(listModal)

	case msgs.DependencyRemoveRequestedMsg:
		deps, _ := m.app.Issues.GetDependencies(context.Background(), msg.IssueID)
		if len(deps) == 0 {
			m.logAction("tui no dependencies to remove")
			return m, nil
		}
		listModal := modal.NewDependencyRemoveListModal(deps, m.width, m.height)
		return m, m.modalManager.PushModal(listModal)

	case msgs.DependencyAddedMsg:
		if m.dependenciesModal != nil {
			m.dependenciesModal.RefreshDeps()
		}
		if msg.Err != nil {
			m.logAction("tui failed to add dependency")
			return m, nil
		}
		m.logAction("tui added dependency")
		m.submitValidation()
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case msgs.DependencyRemovedMsg:
		if m.dependenciesModal != nil {
			m.dependenciesModal.RefreshDeps()
		}
		if msg.Err != nil {
			m.logAction("tui failed to remove dependency")
			return m, nil
		}
		m.logAction("tui removed dependency")
		m.submitValidation()
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case msgs.TitleUpdatedMsg:
		m.modalManager.GetTextInputModal(modal.ModalEditTitle).Reset()
		m.currentIssueID = ""
		if msg.Err != nil {
			m.logAction("tui failed to update issue title")
			return m, nil
		}
		m.logAction("tui updated issue title")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.DescriptionUpdatedMsg:
		m.modalManager.GetTextAreaModal(modal.ModalEditDescription).Reset()
		m.currentIssueID = ""
		if msg.Err != nil {
			m.logAction("tui failed to update issue description")
			return m, nil
		}
		m.logAction("tui updated issue description")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.StatusUpdatedMsg:
		m.currentIssueID = ""
		if msg.Err != nil {
			m.logAction("tui failed to update issue status")
			return m, nil
		}
		m.logAction("tui updated issue status")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.PriorityUpdatedMsg:
		m.currentIssueID = ""
		if msg.Err != nil {
			m.logAction("tui failed to update issue priority")
			return m, nil
		}
		m.logAction("tui updated issue priority")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.TypeUpdatedMsg:
		m.currentIssueID = ""
		if msg.Err != nil {
			m.logAction("tui failed to update issue type")
			return m, nil
		}
		m.logAction("tui updated issue type")
		return m, m.refreshAndSubmit(msg.IssueID)

	case msgs.AssigneeUpdatedMsg:
		m.modalManager.GetTextInputModal(modal.ModalEditAssignee).Reset()
		m.currentIssueID = ""
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
		m.modalManager.GetTextInputModal(modal.ModalCreateIssue).Reset()
		if msg.Err != nil || msg.Issue == nil {
			m.logAction("tui failed to create issue")
			return m, nil
		}
		allIssues, err := m.app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
		if err != nil {
			return m, nil
		}
		setItemsCmd := m.issueList.SetIssues(components.SortedIssues(allIssues))

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
		m.logAction("tui created issue")
		m.submitValidation()
		return m, tea.Sequence(setItemsCmd, func() tea.Msg { return msgs.SelectIssueMsg{IssueID: selectedIssue.ID} })

	case msgs.IssueCommentAddedMsg:
		m.modalManager.GetTextAreaModal(modal.ModalAddComment).Reset()
		m.currentIssueID = ""
		if msg.Err != nil {
			return m, nil
		}
		m.submitValidation()
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case msgs.DeletedMsg:
		m.deleteIndex = -1
		m.currentIssueID = ""
		if msg.Err != nil {
			m.logAction("tui failed to delete issue")
			return m, nil
		}
		m.logAction("tui deleted issue")
		m.submitValidation()
		return m, m.refreshIssueListsAndSelectIssue(msg.IssueID)

	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, nil
		}

		if m.issueList.FilterState() == list.Filtering {
			cmd, _ := m.issueList.Update(msg)
			return m, cmd
		}

		if msg.String() == "esc" {
			return m, nil
		}

		cmd := m.handleKeyPressMsg(msg)
		if cmd != nil || m.IsInModal() {
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

	if m.focusManager.IsFocused(modal.FocusList) {
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
