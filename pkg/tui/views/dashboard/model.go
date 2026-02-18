package dashboard

import (
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/task"
	tea "github.com/charmbracelet/bubbletea"
)

type ValidationFeedbackMsg struct {
	Feedback task.ValidationFeedback
}

type Model struct {
	header          Header
	issueList       IssueList
	issueDetail     IssueDetail
	helpBar         HelpBar
	keyMap          DashboardKeyMap
	svc             *service.Services
	width           int
	height          int
	focusedPane     int // 0 = list, 1 = detail
	feedbackChan    chan task.ValidationFeedback
	quitChan        chan bool
	currentFeedback task.ValidationFeedback
	showComplete    bool
}

func NewDashboard(svc *service.Services, feedbackChan chan task.ValidationFeedback, quitChan chan bool) *Model {
	m := &Model{
		header:       NewHeader("Project Manager Dashboard"),
		keyMap:       defaultDashboardKeyMap,
		svc:          svc,
		width:        80,
		height:       24,
		focusedPane:  0,
		feedbackChan: feedbackChan,
		quitChan:     quitChan,
	}

	m.issueList = NewIssueList(svc, 0, 0)
	m.issueDetail = NewIssueDetail()
	m.helpBar = NewHelpBar(m.keyMap)

	if selected := m.issueList.SelectedItem(); selected.ID != "" {
		m.issueDetail.SetIssue(selected.Issue)
	}

	return m
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
	return m.focusedPane == 0
}

func (m *Model) IsFocusedOnDetail() bool {
	return m.focusedPane == 1
}

func (m *Model) FocusList() {
	m.focusedPane = 0
	m.issueDetail.SetFocused(false)
}

func (m *Model) FocusDetail() {
	m.focusedPane = 1
	m.issueDetail.SetFocused(true)
}

func (m *Model) ToggleFocus() {
	if m.focusedPane == 0 {
		m.FocusDetail()
	} else {
		m.FocusList()
	}
}
