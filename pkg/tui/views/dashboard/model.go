package dashboard

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Use shared types from components for consistency.
type (
	Header      = components.Header
	IssueList   = components.IssueList
	IssueDetail = components.IssueDetail
	ListIssue   = components.ListIssue
)

type Model struct {
	header            Header
	issueList         IssueList
	issueDetail       IssueDetail
	closedIssueList   IssueList
	helpBar           components.HelpBar
	keyMap            DashboardKeyMap
	app               *app.App
	width             int
	height            int
	focusedWindow     int // 0 = main (display issues), 1 = closed issues
	focusedPaneMain   int // 0 = list, 1 = detail
	focusedPaneClosed int
	editingTitle      bool // true while we are editing a title
	titleInput        textinput.Model
	editingIssueID    string

	editingDescription bool // true while editing a description
	descriptionInput   textarea.Model
	editingDescIssueID string
	creatingIssue      bool // true while creating a new issue
	createTitleInput   textinput.Model

	confirmingDelete   bool // true while confirming a delete
	deleteConfirmID    string
	deleteConfirmIndex int

	choosingStatus      bool 
	statusIssueID       string
	choosingPriority    bool 
	priorityIssueID     string
	choosingType        bool 
	typeIssueID         string
	editingAssignee     bool 
	assigneeInput       textinput.Model
	assigneeIssueID     string
	addingComment       bool 
	commentInput        textarea.Model
	commentIssueID      string
	choosingCloseReason bool 
	closeReasonIssueID  string
	closingOtherReason  bool 
	closeReasonInput    textarea.Model
	feedbackChan        chan models.ValidationFeedback
	quitChan            chan bool
	currentFeedback     models.ValidationFeedback
	showComplete        bool
	submitChan          chan<- struct{}
}

func NewDashboard(app *app.App, feedbackChan chan models.ValidationFeedback, quitChan chan bool, submitChan chan<- struct{}) *Model {
	m := &Model{
		header:            components.NewHeader("Project Manager Dashboard"),
		keyMap:            defaultDashboardKeyMap,
		app:               app,
		width:             80,
		height:            24,
		focusedWindow:     0,
		focusedPaneMain:   0,
		focusedPaneClosed: 0,
		feedbackChan:      feedbackChan,
		quitChan:          quitChan,
		submitChan:        submitChan,
	}

	allIssues, _ := app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
	m.issueList = components.NewIssueListFromIssues(app, components.OpenAndInProgressOnly(allIssues), 0, 0)
	m.issueDetail = components.NewIssueDetail()
	m.closedIssueList = components.NewIssueListFromIssues(app, components.ClosedOnly(allIssues), 0, 0)
	m.helpBar = components.NewHelpBar(components.ViewIssues)

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

	if selected := m.issueList.SelectedItem(); selected.ID != "" {
		m.setDetailIssueWithComments(selected.Issue)
	} else if selected := m.closedIssueList.SelectedItem(); selected.ID != "" {
		m.setDetailIssueWithComments(selected.Issue)
	}

	return m
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

func (m *Model) startAddComment(selected ListIssue) {
	m.addingComment = true
	m.commentIssueID = selected.ID
	m.commentInput.SetValue("")
	m.commentInput.Reset()
}

func (m *Model) logAction(action string) {
	if m.app != nil {
		m.app.LogAction(models.EncodeActionEvent(models.ActionEvent{
			Source: "tui",
			Action: action,
		}))
	}
}

func (m *Model) startEditTitle(selected ListIssue) {
	m.editingTitle = true
	m.editingIssueID = selected.ID
	m.titleInput.SetValue(selected.Issue.Title)
	m.titleInput.CursorEnd()
}

func (m *Model) startEditDescription(selected ListIssue) {
	m.editingDescription = true
	m.editingDescIssueID = selected.ID
	m.descriptionInput.SetValue(selected.Issue.Description)
	m.descriptionInput.CursorEnd()
}

func (m *Model) startCreateIssue() {
	m.creatingIssue = true
	m.createTitleInput.SetValue("")
	m.createTitleInput.Reset()
}

func (m *Model) startConfirmDelete(issueID string, index int) {
	m.confirmingDelete = true
	m.deleteConfirmID = issueID
	m.deleteConfirmIndex = index
}

func (m *Model) startChooseStatus(selected ListIssue) {
	m.choosingStatus = true
	m.statusIssueID = selected.ID
}

func (m *Model) startChoosePriority(selected ListIssue) {
	m.choosingPriority = true
	m.priorityIssueID = selected.ID
}

func (m *Model) startChooseType(selected ListIssue) {
	m.choosingType = true
	m.typeIssueID = selected.ID
}

func (m *Model) startEditAssignee(selected ListIssue) {
	m.editingAssignee = true
	m.assigneeIssueID = selected.ID
	m.assigneeInput.SetValue(selected.Assignee)
	m.assigneeInput.CursorEnd()
}

func (m *Model) Init() tea.Cmd {
	return components.ListenForValidation(m.feedbackChan)
}

// IsInModal returns true when a modal (edit, create, delete confirm, choose status/priority/type) is active.
func (m *Model) IsInModal() bool {
	return m.editingTitle || m.creatingIssue || m.editingDescription ||
		m.choosingStatus || m.choosingPriority || m.confirmingDelete ||
		m.choosingType || m.editingAssignee ||
		m.choosingCloseReason || m.closingOtherReason
}

func (m *Model) IsFocusedOnList() bool {
	if m.focusedWindow == 0 {
		return m.focusedPaneMain == 0
	}
	return m.focusedPaneClosed == 0
}

func (m *Model) IsFocusedOnDetail() bool {
	if m.focusedWindow == 0 {
		return m.focusedPaneMain == 1
	}
	return m.focusedPaneClosed == 1
}

func (m *Model) FocusList() {
	if m.focusedWindow == 0 {
		m.focusedPaneMain = 0
	} else {
		m.focusedPaneClosed = 0
	}
	m.issueDetail.SetFocused(false)
}

func (m *Model) FocusDetail() {
	if m.focusedWindow == 0 {
		m.focusedPaneMain = 1
	} else {
		m.focusedPaneClosed = 1
	}
	m.issueDetail.SetFocused(true)
}

func (m *Model) ToggleFocus() {
	if m.IsFocusedOnList() {
		m.FocusDetail()
	} else {
		m.FocusList()
	}
}

func (m *Model) FocusedIssueList() *IssueList {
	// return the issue list of the currently focused window so we can use two tui windows for issues
	if m.focusedWindow == 0 {
		return &m.issueList
	}
	return &m.closedIssueList
}

func (m *Model) ToggleFocusedWindow() {
	// switch focus between open/in-progress and closed issues window
	m.focusedWindow = 1 - m.focusedWindow
	if m.focusedWindow == 0 {
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			m.setDetailIssueWithComments(selected.Issue)
		}
	} else {
		if selected := m.closedIssueList.SelectedItem(); selected.ID != "" {
			m.setDetailIssueWithComments(selected.Issue)
		}
	}
	m.issueDetail.SetFocused(m.IsFocusedOnDetail())
}
