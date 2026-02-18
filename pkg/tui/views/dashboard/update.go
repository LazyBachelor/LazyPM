package dashboard

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)

type deleteResultMsg struct {
	err error
}

func (m *Model) deleteIssueCmd() tea.Cmd {
	selected := m.issueList.SelectedItem()
	if selected.ID == "" {
		return nil
	}

	return func() tea.Msg {
		err := m.svc.Beads.DeleteIssue(context.Background(), selected.ID)
		return deleteResultMsg{err: err}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.showDeleteConfirm {
			switch msg.String() {
			case "y":
				return m, m.deleteIssueCmd()
			case "n", "esc":
				m.showDeleteConfirm = false
				return m, nil
			default:
				return m, nil
			}
		}

		if m.issueList.FilterState() == list.Filtering {
			cmd, _ := m.issueList.Update(msg)
			return m, cmd
		}

		if key.Matches(msg, m.keyMap.Delete) {
			m.deleteErr = nil
			selected := m.issueList.SelectedItem()
			if selected.ID != "" {
				m.showDeleteConfirm = true
			}
			return m, nil
		}

		m.deleteErr = nil
		cmd := m.handleKeyMsg(msg)
		if cmd != nil {
			return m, cmd
		}

	case deleteResultMsg:
		m.showDeleteConfirm = false

		if msg.err != nil {
			m.deleteErr = msg.err
			return m, nil
		}

		m.deleteErr = nil
		index := m.issueList.Index()
		m.issueList = m.issueList.Remove(index)

		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	cmd, changed := m.issueList.Update(msg)
	if changed {
		if selected := m.issueList.SelectedItem(); selected.ID != "" {
			m.issueDetail.SetIssue(selected.Issue)
		}
	}

	return m, cmd
}
