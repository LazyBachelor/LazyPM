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
	backlogList IssueList
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

	currentSprintNum int
	backlogNum       int

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

	m.issueDetail = components.NewIssueDetail()
	m.helpBar = components.NewHelpBar(components.ViewKanban)

	ctx := context.Background()
	backlogNum, _ := app.Issues.GetBacklogSprint(ctx)
	m.backlogNum = backlogNum

	backlogIssues, _ := app.Issues.GetIssuesNotInAnySprint(ctx)
	m.backlogList = components.NewIssueListFromIssues(app, backlogIssues, 20, 10)

	sprints, _ := app.Issues.GetSprints(ctx)
	m.currentSprintNum = 0
	for _, sprintNum := range sprints {
		if sprintNum != backlogNum {
			m.currentSprintNum = sprintNum
			break
		}
	}

	var sprintIssues []*models.Issue
	if m.currentSprintNum > 0 {
		sprintIssues, _ = app.Issues.GetIssuesBySprint(ctx, m.currentSprintNum)
	}
	todoIssues := components.StatusOnly(sprintIssues, models.StatusOpen)
	inProgIssues := components.StatusOnly(sprintIssues, models.StatusInProgress)
	blockedIssues := components.StatusOnly(sprintIssues, models.StatusBlocked)
	doneIssues := components.StatusOnly(sprintIssues, models.StatusClosed)

	m.todoList = components.NewIssueListFromIssues(app, todoIssues, 20, 10)
	m.inProgList = components.NewIssueListFromIssues(app, inProgIssues, 20, 10)
	m.blockedList = components.NewIssueListFromIssues(app, blockedIssues, 20, 10)
	m.doneList = components.NewIssueListFromIssues(app, doneIssues, 20, 10)

	// Setup focus areas for kanban columns
	m.focusManager.EnableArea(modal.FocusColumn0)
	m.focusManager.EnableArea(modal.FocusColumn1)
	m.focusManager.EnableArea(modal.FocusColumn2)
	m.focusManager.EnableArea(modal.FocusColumn3)
	m.focusManager.EnableArea(modal.FocusColumn4)
	m.focusManager.SetCurrent(modal.FocusColumn0)

	// Register modals
	m.registerModals()

	if selected := m.backlogList.SelectedItem(); selected.ID != "" {
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
		m.focusManager.SetCurrent(modal.FocusColumn0)
		m.issueDetail.SetFocused(false)
	} else {
		m.focusManager.SetCurrent(modal.FocusDetail)
		m.issueDetail.SetFocused(true)
	}
}

func (m *Model) FocusedIssueList() *IssueList {
	switch m.focusManager.Current() {
	case modal.FocusColumn0:
		return &m.backlogList
	case modal.FocusColumn1:
		return &m.todoList
	case modal.FocusColumn2:
		return &m.inProgList
	case modal.FocusColumn3:
		return &m.blockedList
	case modal.FocusColumn4:
		return &m.doneList
	default:
		return &m.backlogList
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
	case modal.FocusColumn0:
		return models.StatusOpen
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
	case modal.FocusColumn0: // Backlog
		if delta > 0 {
			return m.moveIssueToSprint(selected.ID, m.currentSprintNum)
		}
	case modal.FocusColumn1: // To Do
		if delta > 0 {
			newCol = modal.FocusColumn2 // To In Progress
		} else if delta < 0 {
			return m.moveIssueToBacklog(selected.ID)
		}
	case modal.FocusColumn2: // In Progress
		if delta > 0 {
			newCol = modal.FocusColumn3 // To Blocked
		} else {
			newCol = modal.FocusColumn1 // To To Do
		}
	case modal.FocusColumn3: // Blocked
		if delta > 0 {
			newCol = modal.FocusColumn4 // To Done
		} else {
			newCol = modal.FocusColumn2 // To In Progress
		}
	case modal.FocusColumn4: // Done
		if delta < 0 {
			newCol = modal.FocusColumn3 // To Blocked
		}
	}

	if newCol == modal.FocusNone {
		return nil
	}

	newStatus := statusForColumn(newCol)
	return msgs.UpdateIssueStatusCmd(m.app, selected.ID, string(newStatus))
}

// moveIssueToSprint moves an issue to a sprint
func (m *Model) moveIssueToSprint(issueID string, sprintNum int) tea.Cmd {
	return func() tea.Msg {
		err := m.app.Issues.AddIssueToSprint(context.Background(), issueID, sprintNum)
		m.submitValidation()
		if err != nil {
			return m.refreshIssueListsAndSelectIssue(issueID)()
		}
		return m.refreshIssueListsAndSelectIssue(issueID)()
	}
}

// moveIssueToBacklog removes an issue from the current sprint (moves it to backlog)
func (m *Model) moveIssueToBacklog(issueID string) tea.Cmd {
	return func() tea.Msg {
		err := m.app.Issues.RemoveIssueFromSprint(context.Background(), issueID, m.currentSprintNum)
		m.submitValidation()
		if err != nil {
			return m.refreshIssueListsAndSelectIssue(issueID)()
		}
		return m.refreshIssueListsAndSelectIssue(issueID)()
	}
}
