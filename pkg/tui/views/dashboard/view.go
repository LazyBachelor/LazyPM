package dashboard

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views/issue"
)

func NewDashboard(svc *service.Services) Model {

	issues, err := svc.Beads.AllIssues(context.Background())

	listIssues := []ListIssue{}
	for _, issue := range issues {
		listIssues = append(listIssues, ListIssue{Issue: issue})
	}

	if err != nil {
		listIssues = []ListIssue{}
	}

	items := make([]list.Item, len(listIssues))
	for i, issue := range listIssues {
		items[i] = issue
	}

	issueList := list.New(items, IssueListDelegate{}, 0, 0)

	issueList.Title = "Issues"
	issueList.Styles.Title = styles.TitleStyle

	issueView := issue.NewIssueView(issue.Model{})

	m := Model{
		issueView: issueView,
		issueList: issueList,
		keyMap:    defaultDashboardKeyMap,
		svc:       svc,
	}

	m.updateIssueView()

	return m
}

func (d Model) View() string {
	issueView := d.issueView.View()

	if d.focusedOnIssue {
		issueView = styles.FocusedIssueStyle.Render(issueView)
	} else {
		issueView = styles.IssueStyle.Render(issueView)
	}

	str := lipgloss.JoinHorizontal(lipgloss.Left, styles.AppStyle.Render(d.issueList.View()), issueView)

	return str
}
