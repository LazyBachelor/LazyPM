package dashboard

import (
	"context"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/modal"
)

// Use shared types from components for consistency.
type (
	Header      = components.Header
	IssueList   = components.IssueList
	IssueDetail = components.IssueDetail
	ListIssue   = components.ListIssue
)

type Model struct {
	header          Header
	issueList       IssueList
	issueDetail     IssueDetail
	closedIssueList IssueList
	helpBar         components.HelpBar
	keyMap          KeyMap
	app             *app.App
	width           int
	height          int

	// Modal and Focus management
	modalManager *modal.Manager
	focusManager *modal.FocusManager

	// Inputs (still needed for modals that require them)
	titleInput       textinput.Model
	descriptionInput textarea.Model
	createTitleInput textinput.Model
	assigneeInput    textinput.Model
	commentInput     textarea.Model
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
		header:       components.NewHeader("Project Manager Dashboard"),
		keyMap:       defaultDashboardKeyMap,
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

	commentTa := textarea.New()
	commentTa.Placeholder = "Write your comment..."
	commentTa.SetWidth(56)
	commentTa.SetHeight(6)
	m.commentInput = commentTa

	// Setup lists
	allIssues, _ := app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
	m.issueList = components.NewIssueListFromIssues(app, components.SortedIssues(allIssues), 0, 0)
	m.issueDetail = components.NewIssueDetail()
	m.helpBar = components.NewHelpBar(components.ViewIssues)

	// Setup focus
	m.focusManager.EnableArea(modal.FocusList)
	m.focusManager.EnableArea(modal.FocusDetail)
	m.focusManager.SetCurrent(modal.FocusList)

	// Register modals
	m.registerModals()

	if selected := m.issueList.SelectedItem(); selected.ID != "" {
		m.setDetailIssueWithComments(selected.Issue)
	}

	return m
}

// registerModals sets up all modals (handlers are in operations.go via ModalCompletedMsg)
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

	// Add Comment Modal
	m.modalManager.RegisterModal(modal.NewTextAreaModal(modal.TextAreaConfig{
		ID:          modal.ModalAddComment,
		Label:       "Add comment (Ctrl+S or Enter to save, Esc to cancel):",
		Placeholder: "Write your comment...",
		SaveKeys:    []string{"ctrl+s", "enter"},
		InputHeight: 8,
	}))

	// Close Reason TextArea Modal
	m.modalManager.RegisterModal(modal.NewTextAreaModal(modal.TextAreaConfig{
		ID:          modal.ModalCloseReason,
		Label:       "Enter closing reason (Enter or Ctrl+S to save, Esc to cancel):",
		Placeholder: "Enter closing reason...",
		SaveKeys:    []string{"enter", "ctrl+s"},
		InputHeight: 4,
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
}

func (m *Model) Init() tea.Cmd {
	if m.submitChan != nil {
		m.submitChan <- struct{}{}
		m.logAction("tui submitted validation")
	}
	return components.ListenForValidation(m.feedbackChan)
}

// setDetailIssueWithComments sets the issue in the detail pane and loads its comments.
func (m *Model) setDetailIssueWithComments(issue models.Issue) {
	m.issueDetail.SetIssue(issue)
	if issue.ID == "" {
		m.issueDetail.SetComments(nil)
		return
	}
	comments, _ := m.app.Issues.GetIssueComments(context.Background(), issue.ID)
	m.issueDetail.SetComments(comments)
}

func (m *Model) logAction(action string) {
	if m.app != nil {
		m.app.LogAction(models.EncodeActionEvent(models.ActionEvent{
			Source: "tui",
			Action: action,
		}))
	}
}

// submitValidation sends a validation request to the submit channel.
func (m *Model) submitValidation() {
	if m.submitChan != nil {
		select {
		case m.submitChan <- struct{}{}:
			m.logAction("tui submitted validation")
		default:
		}
	}
}

// IsInModal returns true when a modal is active (kept for backwards compatibility).
func (m *Model) IsInModal() bool {
	return m.modalManager.IsModalActive()
}

func (m *Model) IsFocusedOnList() bool {
	return m.focusManager.IsListFocused()
}

func (m *Model) IsFocusedOnDetail() bool {
	return m.focusManager.IsDetailFocused()
}
