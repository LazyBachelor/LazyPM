package dashboard

import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ValidationFeedbackMsg struct {
	Feedback models.ValidationFeedback
}

type Model struct {
	header            Header
	issueList         IssueList
	issueDetail       IssueDetail
	closedIssueList   IssueList
	helpBar           HelpBar
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

	choosingStatus   bool // true while choosing a status
	statusIssueID    string
	choosingPriority bool // true while choosing a priority
	priorityIssueID  string
	choosingType     bool // true while choosing a type
	typeIssueID      string
	addingComment    bool // true while adding a comment
	commentInput     textarea.Model
	commentIssueID   string
	feedbackChan     chan models.ValidationFeedback
	quitChan         chan bool
	currentFeedback  models.ValidationFeedback
	showComplete     bool
}

func NewDashboard(app *app.App, feedbackChan chan models.ValidationFeedback, quitChan chan bool) *Model {
	m := &Model{
		header:            NewHeader("Project Manager Dashboard"),
		keyMap:            defaultDashboardKeyMap,
		app:               app,
		width:             80,
		height:            24,
		focusedWindow:     0,
		focusedPaneMain:   0,
		focusedPaneClosed: 0,
		feedbackChan:      feedbackChan,
		quitChan:          quitChan,
	}

	allIssues, _ := app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
	m.issueList = NewIssueListFromIssues(app, OpenAndInProgressOnly(allIssues), 0, 0)
	m.issueDetail = NewIssueDetail()
	m.closedIssueList = NewIssueListFromIssues(app, ClosedOnly(allIssues), 0, 0)
	m.helpBar = NewHelpBar(m.keyMap)

	ti := textinput.New()
	ti.Placeholder = "Issue title ..."
	ti.CharLimit = 256
	m.titleInput = ti

	createTi := textinput.New()
	createTi.Placeholder = "New issue title ..."
	createTi.CharLimit = 256
	m.createTitleInput = createTi

	descTa := textarea.New()
	descTa.Placeholder = "Issue description..."
	descTa.SetWidth(56)
	descTa.SetHeight(8)
	m.descriptionInput = descTa

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

func (m *Model) Init() tea.Cmd {
	return m.listenForValidation()
}

func (m *Model) listenForValidation() tea.Cmd {
	return func() tea.Msg {
		feedback := <-m.feedbackChan
		return ValidationFeedbackMsg{Feedback: feedback}
	}
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
