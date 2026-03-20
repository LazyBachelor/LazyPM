package kanban

import (
	"context"

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
		m.setDetailIssueWithComments(selected.Issue)
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	if m.submitChan != nil {
		m.submitChan <- struct{}{}
		m.logAction("tui submitted validation")
	}
	return components.ListenForValidation(m.feedbackChan)
}

func (m *Model) registerModals() {
	modal.RegisterCommonModals(m.modalManager)
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

func (m *Model) IsInModal() bool {
	return m.modalManager.IsModalActive()
}

func (m *Model) IsFocusedOnList() bool {
	return m.focusManager.IsListFocused()
}

func (m *Model) IsFocusedOnDetail() bool {
	return m.focusManager.IsDetailFocused()
}

func (m *Model) ToggleFocus() {
	if m.focusManager.IsDetailFocused() {
		m.focusManager.SetCurrent(modal.FocusColumn1)
		m.issueDetail.SetFocused(false)
	} else {
		m.focusManager.SetCurrent(modal.FocusDetail)
		m.issueDetail.SetFocused(true)
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
		m.setDetailIssueWithComments(selected.Issue)
	}
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
