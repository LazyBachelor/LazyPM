package dashboard

import (
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/tui/views/issue"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	issueView      issue.Model
	issueList      list.Model
	keyMap         DashboardKeyMap
	svc            *service.Services
	focusedOnIssue bool
}

func (d Model) Init() tea.Cmd {
	return nil
}

type ListIssue struct {
	models.Issue
}

func (i ListIssue) Title() string       { return i.Issue.Title }
func (i ListIssue) Description() string { return i.Issue.Description }
func (i ListIssue) FilterValue() string { return i.Issue.ID + " " + i.Issue.Title }
