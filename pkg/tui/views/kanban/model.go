package kanban

import (
	"context"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/modal"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
)

type (
	Header      = components.Header
	IssueList   = components.IssueList
	IssueDetail = components.IssueDetail
	ListIssue   = components.ListIssue
)

type Model struct {
	header      Header
	todoList    IssueList
	inProgList  IssueList
	blockedList IssueList
	doneList    IssueList
	issueDetail IssueDetail
	helpBar     components.HelpBar
	keyMap      KeyMap
	app         *app.App
	width       int
	height      int

	// Modal and Focus management
	modalManager *modal.Manager
	focusManager *modal.FocusManager

	// Inputs
	titleInput       textinput.Model
	descriptionInput textarea.Model
	createTitleInput textinput.Model
	assigneeInput    textinput.Model
	closeReasonInput textarea.Model

	// Current issue being operated on
	currentIssueID string
	deleteIndex    int

	feedbackChan    chan models.ValidationFeedback
	quitChan        chan bool
	currentFeedback models.ValidationFeedback
	showComplete    bool
	submitChan      chan<- struct{}
}

func NewDashboard(app *app.App, feedbackChan chan models.ValidationFeedback, quitChan chan bool, submitChan chan<- struct{}) *Model {
	m := &Model{
		header:       components.NewHeader("Kanban Board"),
		keyMap:       defaultKanbanKeyMap,
		app:          app,
		width:        80,
		height:       24,
		feedbackChan: feedbackChan,
		quitChan:     quitChan,
		submitChan:   submitChan,
		modalManager: modal.NewManager(),
		focusManager: modal.NewFocusManager(),
		deleteIndex:  -1,
	}

	// Setup inputs
	inputs := components.NewIssueInputs()
	m.titleInput = inputs.Title
	m.createTitleInput = inputs.CreateTitle
	m.descriptionInput = inputs.Description
	m.assigneeInput = inputs.Assignee

	closeReasonTa := textarea.New()
	closeReasonTa.Placeholder = "Enter closing reason..."
	closeReasonTa.SetWidth(56)
	closeReasonTa.SetHeight(4)
	m.closeReasonInput = closeReasonTa

	// Setup lists
	m.issueDetail = components.NewIssueDetail()
	m.helpBar = components.NewHelpBar(components.ViewKanban)

	allIssues, _ := app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
	todoIssues := components.StatusOnly(allIssues, models.StatusOpen)
	inProgIssues := components.StatusOnly(allIssues, models.StatusInProgress)
	blockedIssues := components.StatusOnly(allIssues, models.StatusBlocked)
	doneIssues := components.StatusOnly(allIssues, models.StatusClosed)

	m.todoList = components.NewIssueListFromIssues(app, todoIssues, 20, 10)
	m.inProgList = components.NewIssueListFromIssues(app, inProgIssues, 20, 10)
	m.blockedList = components.NewIssueListFromIssues(app, blockedIssues, 20, 10)
	m.doneList = components.NewIssueListFromIssues(app, doneIssues, 20, 10)

	// Setup focus areas for kanban columns
	m.focusManager.EnableArea(modal.FocusColumn1)
	m.focusManager.EnableArea(modal.FocusColumn2)
	m.focusManager.EnableArea(modal.FocusColumn3)
	m.focusManager.EnableArea(modal.FocusColumn4)
	m.focusManager.SetCurrent(modal.FocusColumn1)

	// Register modals
	m.registerModals()

	if selected := m.todoList.SelectedItem(); selected.ID != "" {
		m.issueDetail.SetIssue(selected.Issue)
	}

	return m
}

func (m *Model) registerModals() {
	// Edit Title Modal
	m.modalManager.RegisterModal(modal.NewTextInputModal(modal.TextInputConfig{
		ID:           modal.ModalEditTitle,
		Label:        "Edit title (Enter to save, Esc to cancel):",
		Placeholder:  "Issue title...",
		SaveKeys:     []string{"enter"},
		CharLimit:    256,
		InitialValue: "",
	}))

	// Create Issue Modal
	m.modalManager.RegisterModal(modal.NewTextInputModal(modal.TextInputConfig{
		ID:          modal.ModalCreateIssue,
		Label:       "New issue (Enter to create, Esc to cancel):",
		Placeholder: "New issue title...",
		SaveKeys:    []string{"enter"},
		CharLimit:   256,
	}))

	// Edit Assignee Modal
	m.modalManager.RegisterModal(modal.NewTextInputModal(modal.TextInputConfig{
		ID:          modal.ModalEditAssignee,
		Label:       "Edit assignee (Enter to save, Esc to cancel):",
		Placeholder: "Assignee name...",
		SaveKeys:    []string{"enter"},
		CharLimit:   64,
	}))

	// Edit Description Modal
	m.modalManager.RegisterModal(modal.NewTextAreaModal(modal.TextAreaConfig{
		ID:          modal.ModalEditDescription,
		Label:       "Edit description (Ctrl+S to save, Esc to cancel):",
		Placeholder: "Issue description...",
		SaveKeys:    []string{"ctrl+s"},
		InputHeight: 10,
	}))

	// Delete Confirm Modal
	m.modalManager.RegisterModal(modal.NewConfirmModal(modal.ConfirmConfig{
		ID:      modal.ModalConfirmDelete,
		Message: "Delete issue?",
		YesKeys: []string{"y", "Y"},
		NoKeys:  []string{"n", "N", "esc"},
	}))

	// Status Select Modal
	m.modalManager.RegisterModal(modal.NewSelectModal(modal.SelectConfig{
		ID:      modal.ModalSelectStatus,
		Label:   "Change status:",
		Options: modal.StatusOptions(),
	}))

	// Close Reason Select Modal
	m.modalManager.RegisterModal(modal.NewSelectModal(modal.SelectConfig{
		ID:      modal.ModalSelectCloseReason,
		Label:   "Choose closing reason:",
		Options: modal.CloseReasonOptions(),
	}))

	// Priority Select Modal
	m.modalManager.RegisterModal(modal.NewSelectModal(modal.SelectConfig{
		ID:      modal.ModalSelectPriority,
		Label:   "Change priority:",
		Options: modal.PriorityOptions(),
	}))

	// Type Select Modal
	m.modalManager.RegisterModal(modal.NewSelectModal(modal.SelectConfig{
		ID:      modal.ModalSelectType,
		Label:   "Change type:",
		Options: modal.TypeOptions(),
	}))

	// Close Reason TextArea Modal
	m.modalManager.RegisterModal(modal.NewTextAreaModal(modal.TextAreaConfig{
		ID:          modal.ModalCloseReason,
		Label:       "Enter closing reason (Enter or Ctrl+S to save, Esc to cancel):",
		Placeholder: "Enter closing reason...",
		SaveKeys:    []string{"enter", "ctrl+s"},
		InputHeight: 4,
	}))
}

func (m *Model) logAction(action string) {
	if m.app != nil {
		m.app.LogAction(models.EncodeActionEvent(models.ActionEvent{
			Source: "tui",
			Action: action,
		}))
	}
}

func (m *Model) submitValidation() {
	if m.submitChan != nil {
		select {
		case m.submitChan <- struct{}{}:
			m.logAction("tui submitted validation")
		default:
		}
	}
}

func (m *Model) Init() tea.Cmd {
	if m.submitChan != nil {
		m.submitChan <- struct{}{}
		m.logAction("tui submitted validation")
	}
	return components.ListenForValidation(m.feedbackChan)
}

func (m *Model) IsInModal() bool {
	return m.modalManager.IsModalActive()
}

func (m *Model) IsFocusedOnList() bool {
	return m.focusManager.IsListFocused()
}

func (m *Model) IsFocusedOnDetail() bool {
	return m.focusManager.IsDetailFocused()
}

func (m *Model) FocusList() {
	m.focusManager.SetCurrent(modal.FocusColumn1)
	m.issueDetail.SetFocused(false)
}

func (m *Model) FocusDetail() {
	m.focusManager.SetCurrent(modal.FocusDetail)
	m.issueDetail.SetFocused(true)
}

func (m *Model) ToggleFocus() {
	if m.focusManager.IsDetailFocused() {
		m.FocusList()
	} else {
		m.FocusDetail()
	}
}

func (m *Model) FocusedIssueList() *IssueList {
	switch m.focusManager.Current() {
	case modal.FocusColumn1:
		return &m.todoList
	case modal.FocusColumn2:
		return &m.inProgList
	case modal.FocusColumn3:
		return &m.blockedList
	case modal.FocusColumn4:
		return &m.doneList
	default:
		return &m.todoList
	}
}

func (m *Model) updateDetailFromSelection() {
	selected := m.FocusedIssueList().SelectedItem()
	if selected.ID != "" {
		m.issueDetail.SetIssue(selected.Issue)
	}
}

func statusForColumn(col modal.FocusArea) models.Status {
	switch col {
	case modal.FocusColumn1:
		return models.StatusOpen
	case modal.FocusColumn2:
		return models.StatusInProgress
	case modal.FocusColumn3:
		return models.StatusBlocked
	case modal.FocusColumn4:
		return models.StatusClosed
	default:
		return models.StatusOpen
	}
}

func (m *Model) moveIssue(delta int) tea.Cmd {
	fl := m.FocusedIssueList()
	selected := fl.SelectedItem()
	if selected.ID == "" {
		return nil
	}

	currentCol := m.focusManager.Current()
	var newCol modal.FocusArea
	switch currentCol {
	case modal.FocusColumn1:
		if delta > 0 {
			newCol = modal.FocusColumn2
		}
	case modal.FocusColumn2:
		if delta > 0 {
			newCol = modal.FocusColumn3
		} else {
			newCol = modal.FocusColumn1
		}
	case modal.FocusColumn3:
		if delta > 0 {
			newCol = modal.FocusColumn4
		} else {
			newCol = modal.FocusColumn2
		}
	case modal.FocusColumn4:
		if delta < 0 {
			newCol = modal.FocusColumn3
		}
	}

	if newCol == modal.FocusNone {
		return nil
	}

	newStatus := statusForColumn(newCol)
	return msgs.UpdateIssueStatusCmd(m.app, selected.ID, string(newStatus))
}
