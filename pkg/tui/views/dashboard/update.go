package dashboard

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.issueList.FilterState() == list.Filtering {
			cmd, _ := m.issueList.Update(msg)
			return m, cmd
		}

		cmd := m.handleKeyMsg(msg)
		if cmd != nil {
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case ValidationFeedbackMsg:
		m.currentFeedback = msg.Feedback
		if msg.Feedback.Success {
			m.showComplete = true
			return m, tea.Quit
		}
		return m, m.listenForValidation()
	}

	cmd, changed := m.issueList.Update(msg)
	if changed {
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			m.issueDetail.SetIssue(selected.Issue)
		}
	}

	return m, cmd
}
