package dashboard

import (
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
		err := m.svc.DeleteIssue(selected.ID)
		return deleteResultMsg{err: err}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		// 🔴 Hvis vi er i delete-confirmation mode
		if m.showDeleteConfirm {
			switch msg.String() {
			case "y":
				return m, m.deleteIssueCmd()
			case "n", "esc":
				m.showDeleteConfirm = false
				return m, nil
			}
		}

		// Hvis listen er i filter mode
		if m.issueList.FilterState() == list.Filtering {
			cmd, _ := m.issueList.Update(msg)
			return m, cmd
		}

		// 🔴 Trigger delete
		if msg.String() == "d" {
			selected := m.issueList.SelectedItem()
			if selected.ID != "" {
				m.showDeleteConfirm = true
			}
			return m, nil
		}

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

		// Fjern item fra bubbles list
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
