package dashboard

import (
	"fmt"

	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (d Model) Update(m tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := m.(type) {
	case tea.KeyMsg:
		if d.issueList.FilterState() == list.Filtering {
			break
		}
		cmd := d.handleKeyMsg(msg)
		if cmd != nil {
			return d, cmd
		}
	case tea.WindowSizeMsg:
		d.updateSizes(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	oldIndex := d.issueList.Index()
	d.issueList, cmd = d.issueList.Update(m)

	if d.issueList.Index() != oldIndex {
		d.updateIssueView()
	}

	return d, cmd
}

func (d *Model) updateIssueView() {
	if item, ok := d.issueList.SelectedItem().(ListIssue); ok {
		d.issueView.ID = item.ID
		d.issueView.Title = item.Issue.Title
		d.issueView.Description = item.Issue.Description
		d.issueView.Status = string(item.Issue.Status)
		d.issueView.IssueType = string(item.Issue.IssueType)

		content := fmt.Sprintf("ID: %s\nTitle: %s\nDescription: %s\nStatus: %s\nType: %s",
			d.issueView.ID, d.issueView.Title, d.issueView.Description, d.issueView.Status, d.issueView.IssueType)

		d.issueView.Viewport.SetContent(content)
	}
}

func (d *Model) updateSizes(width, height int) {
	listWidth := width / 2
	issueWidth := width - listWidth

	w, h := styles.AppStyle.GetFrameSize()
	d.issueList.SetSize(listWidth-w, height-h)
	d.issueView.SetSize(issueWidth, height-h)
	d.issueView.Viewport.Width = issueWidth
	d.issueView.Viewport.Height = height - h
}
