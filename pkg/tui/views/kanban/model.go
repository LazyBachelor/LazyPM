package kanban
import (
	"context"

	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/components"
	"github.com/LazyBachelor/LazyPM/pkg/tui/issues"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	keyMap      KanbanKeyMap
	app         *app.App
	width       int
	height      int

	focusedColumn int  // 0 = To Do, 1 = In Progress, 2 = Blocked, 3 = Done
	focusOnDetail bool // true when detail pane is focused

	editingTitle   bool // true while we are editing a title
	titleInput     textinput.Model
	editingIssueID string

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
	editingAssignee  bool // true while editing assignee
	assigneeInput    textinput.Model
	assigneeIssueID  string
	choosingCloseReason bool // true while choosing a close reason
	closeReasonIssueID  string
	closingOtherReason  bool // true while entering a custom close reason
	closeReasonInput    textarea.Model
	feedbackChan        chan models.ValidationFeedback
	quitChan            chan bool
	submitChan          chan<- struct{}
	currentFeedback     models.ValidationFeedback
	showComplete        bool
}

func NewDashboard(app *app.App, feedbackChan chan models.ValidationFeedback, quitChan chan bool, submitChan chan<- struct{}) *Model {
	m := &Model{
		header:       components.NewHeader("Kanban Board"),
		keyMap:       defaultKanbanKeyMap,
		app:          app,
		width:        80,
		height:       24,
		focusedColumn: 0,
		focusOnDetail: false,
		feedbackChan:  feedbackChan,
		quitChan:      quitChan,
		submitChan:    submitChan,
	}

	allIssues, _ := app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
	todoIssues := components.StatusOnly(allIssues, models.StatusOpen)
	inProgIssues := components.StatusOnly(allIssues, models.StatusInProgress)
	blockedIssues := components.StatusOnly(allIssues, models.StatusBlocked)
	doneIssues := components.StatusOnly(allIssues, models.StatusClosed)

	m.todoList = components.NewIssueListFromIssues(app, todoIssues, 0, 0)
	m.inProgList = components.NewIssueListFromIssues(app, inProgIssues, 0, 0)
	m.blockedList = components.NewIssueListFromIssues(app, blockedIssues, 0, 0)
	m.doneList = components.NewIssueListFromIssues(app, doneIssues, 0, 0)
	m.issueDetail = components.NewIssueDetail()
	m.helpBar = components.NewHelpBar(components.ViewKanban)

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

	if selected := m.todoList.SelectedItem(); selected.ID != "" {
		m.issueDetail.SetIssue(selected.Issue)
	} else if selected := m.inProgList.SelectedItem(); selected.ID != "" {
		m.issueDetail.SetIssue(selected.Issue)
	} else if selected := m.blockedList.SelectedItem(); selected.ID != "" {
		m.issueDetail.SetIssue(selected.Issue)
	} else if selected := m.doneList.SelectedItem(); selected.ID != "" {
		m.issueDetail.SetIssue(selected.Issue)
	}

	return m
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
	return !m.focusOnDetail
}

func (m *Model) IsFocusedOnDetail() bool {
	return m.focusOnDetail
}

func (m *Model) FocusList() {
	m.focusOnDetail = false
	m.issueDetail.SetFocused(false)
}

func (m *Model) FocusDetail() {
	m.focusOnDetail = true
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
	switch m.focusedColumn {
	case 0:
		return &m.todoList
	case 1:
		return &m.inProgList
	case 2:
		return &m.blockedList
	case 3:
		return &m.doneList
	default:
		return &m.todoList
	}
}

// updateDetailFromSelection updates the detail pane based on the currently
// focused column's selected issue.
func (m *Model) updateDetailFromSelection() {
	selected := m.FocusedIssueList().SelectedItem()
	if selected.ID != "" {
		m.issueDetail.SetIssue(selected.Issue)
	}
}

// statusForColumn maps a board column index to a Status.
func statusForColumn(col int) models.Status {
	switch col {
	case 0:
		return models.StatusOpen
	case 1:
		return models.StatusInProgress
	case 2:
		return models.StatusBlocked
	case 3:
		return models.StatusClosed
	default:
		return models.StatusOpen
	}
}

// moveIssue moves the currently selected issue in the focused column horizontally
// to an adjacent column by updating its status.
func (m *Model) moveIssue(delta int) tea.Cmd {
	fl := m.FocusedIssueList()
	selected := fl.SelectedItem()
	if selected.ID == "" {
		return nil
	}

	newCol := m.focusedColumn + delta
	if newCol < 0 || newCol > 3 {
		return nil
	}

	newStatus := statusForColumn(newCol)
		return issues.UpdateIssueStatusCmd(m.app, selected.ID, string(newStatus))
}
