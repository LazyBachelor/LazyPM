package dashboard

import (
	"context"

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

	// Current issue being operated on
	currentIssueID string
	deleteIndex    int

	dependenciesModal *modal.DependenciesModal

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

func (m *Model) Init() tea.Cmd {
	if m.submitChan != nil {
		m.submitChan <- struct{}{}
		m.logAction("tui submitted validation")
	}
	return components.ListenForValidation(m.feedbackChan)
}

// registerModals sets up all modals using the common registration helper
func (m *Model) registerModals() {
	modal.RegisterCommonModals(m.modalManager)
}

// setDetailIssueWithComments sets the issue in the detail pane and loads its comments and dependencies.
func (m *Model) setDetailIssueWithComments(issue models.Issue) {
	m.issueDetail.SetIssue(issue)
	if issue.ID == "" {
		m.issueDetail.SetComments(nil)
		m.issueDetail.SetDependencies(nil)
		return
	}
	ctx := context.Background()
	comments, _ := m.app.Issues.GetIssueComments(ctx, issue.ID)
	m.issueDetail.SetComments(comments)
	deps, _ := m.app.Issues.GetDependencies(ctx, issue.ID)
	m.issueDetail.SetDependencies(deps)
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

// IsInModal returns true when a modal is active
func (m *Model) IsInModal() bool {
	return m.modalManager.IsModalActive()
}

func (m *Model) IsFocusedOnList() bool {
	return m.focusManager.IsListFocused()
}

func (m *Model) IsFocusedOnDetail() bool {
	return m.focusManager.IsDetailFocused()
}
