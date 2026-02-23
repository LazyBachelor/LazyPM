package dashboard

import (
	"context"
	"io"

	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type IssueList struct {
	list   list.Model
	app    *service.App
	width  int
	height int
}

type ListIssue struct {
	models.Issue
}

func (l ListIssue) Title() string       { return l.Issue.Title }
func (l ListIssue) Description() string { return l.Issue.Description }
func (l ListIssue) FilterValue() string { return l.Issue.ID + " " + l.Issue.Title }

type TableColumn struct {
	width uint
	label string
	key   string
}

func getTableColumns(width int) []TableColumn {
	switch {
	case width < 45:
		return []TableColumn{
			{width: 10, label: "ID", key: "id"},
			{width: uint(width - 10), label: "TITLE", key: "title"},
		}
	case width < 60:
		return []TableColumn{
			{width: 10, label: "ID", key: "id"},
			{width: 20, label: "TITLE", key: "title"},
			{width: 15, label: "STATUS", key: "status"},
		}
	default:
		return []TableColumn{
			{width: 12, label: "ID", key: "id"},
			{width: 20, label: "TITLE", key: "title"},
			{width: 15, label: "STATUS", key: "status"},
			{width: 10, label: "TYPE", key: "type"},
		}
	}
}

func renderHeaders(cols []TableColumn) string {
	var parts []string
	headerStyle := lipgloss.NewStyle().Foreground(styles.FaintText).Bold(true)

	for _, col := range cols {
		colWidth := col.width
		if colWidth == 0 {
			colWidth = 1
		}
		style := lipgloss.NewStyle().Width(int(colWidth))
		headerText := headerStyle.Render(truncate.StringWithTail(col.label, colWidth, "..."))
		parts = append(parts, style.Render(headerText))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

func NewIssueList(app *service.App, width, height int) IssueList {
	issues, err := app.Issues.AllIssues(context.Background())
	if err != nil {
		return IssueList{}
	}

	listIssues := []ListIssue{}
	for _, issue := range issues {
		listIssues = append(listIssues, ListIssue{Issue: issue})
	}

	items := make([]list.Item, len(listIssues))
	for i, issue := range listIssues {
		items[i] = issue
	}

	l := list.New(items, NewIssueListDelegate(width), width, height)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.FilterInput.PromptStyle = styles.FilterPromptStyle
	l.FilterInput.Cursor.Style = styles.FilterStyle
	l.FilterInput.TextStyle = styles.FilterInputStyle
	l.FilterInput.Prompt = "🔍 "

	return IssueList{
		list:   l,
		app:    app,
		width:  width,
		height: height,
	}
}

func NewIssueListFromIssues(app *service.App, issues []models.Issue, width, height int) IssueList {
	// for making an IssueList from a pre-existing list of issues.
	listIssues := make([]list.Item, len(issues))
	for i, issue := range issues {
		listIssues[i] = ListIssue{Issue: issue}
	}
	l := list.New(listIssues, NewIssueListDelegate(width), width, height)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.FilterInput.PromptStyle = styles.FilterPromptStyle
	l.FilterInput.Cursor.Style = styles.FilterStyle
	l.FilterInput.TextStyle = styles.FilterInputStyle
	l.FilterInput.Prompt = "🔍 "
	return IssueList{
		list:   l,
		app:    app,
		width:  width,
		height: height,
	}
}

func OpenAndInProgressOnly(issues []models.Issue) []models.Issue {
	// used to display open & in-progress issues in the first window in the dashboard
	out := make([]models.Issue, 0, len(issues))
	for _, issue := range issues {
		if issue.Status == models.StatusOpen || issue.Status == models.StatusInProgress {
			out = append(out, issue)
		}
	}
	return out
}

func ClosedOnly(issues []models.Issue) []models.Issue {
	// used to display issues in the second window in the dashboard
	out := make([]models.Issue, 0, len(issues))
	for _, issue := range issues {
		if issue.Status == models.StatusClosed {
			out = append(out, issue)
		}
	}
	return out
}

func (l *IssueList) Update(msg tea.Msg) (tea.Cmd, bool) {
	var cmd tea.Cmd
	oldIndex := l.list.Index()
	l.list, cmd = l.list.Update(msg)
	changed := l.list.Index() != oldIndex
	return cmd, changed
}

func (l *IssueList) SetSize(width, height int) {
	l.width = width
	l.height = height

	l.list.SetSize(width, height)
	l.list.SetDelegate(NewIssueListDelegate(width))
}

func (l IssueList) View() string {
	return l.renderResponsive()
}

func (l IssueList) renderResponsive() string {
	cols := getTableColumns(l.width)
	header := renderHeaders(cols)

	var content []string

	if l.list.FilterState() == list.Filtering {
		filterText := l.list.FilterInput.Value()
		filterView := styles.FilterStyle.Render("🔍 " + filterText)
		content = append(content, filterView)
	}

	itemsView := l.renderFilteredItems()
	content = append(content, header, itemsView)

	return styles.ContainerStyle.
		Width(l.width).
		MaxWidth(l.width).
		MaxHeight(l.height).
		Render(lipgloss.JoinVertical(lipgloss.Left, content...))
}

func (l IssueList) renderFilteredItems() string {
	var items []string

	var visibleItems []list.Item
	if l.list.FilterState() == list.Filtering || l.list.FilterState() == list.FilterApplied {
		visibleItems = l.list.VisibleItems()
	} else {
		visibleItems = l.list.Items()
	}

	start, end := l.list.Paginator.GetSliceBounds(len(visibleItems))

	cursor := l.list.Index()

	for i := start; i < end && i < len(visibleItems); i++ {
		isSelected := i == cursor
		if issue, ok := visibleItems[i].(ListIssue); ok {
			cols := getTableColumns(l.width)
			row := renderRow(issue, isSelected, cols)
			items = append(items, row)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

func (l IssueList) SelectedItem() ListIssue {
	if item, ok := l.list.SelectedItem().(ListIssue); ok {
		return item
	}
	return ListIssue{}
}

func (l IssueList) Index() int {
	return l.list.Index()
}

func (l IssueList) FilterState() list.FilterState {
	return l.list.FilterState()
}

func (l *IssueList) SetIssues(issues []models.Issue) tea.Cmd {
	listIssues := make([]list.Item, len(issues))
	for i, issue := range issues {
		listIssues[i] = ListIssue{Issue: issue}
	}
	return l.list.SetItems(listIssues)
}

func (l *IssueList) SelectIssueID(issueID string) {
	items := l.list.Items()
	for i := 0; i < len(items); i++ {
		if item, ok := items[i].(ListIssue); ok && item.ID == issueID {
			l.list.Select(i)
			return
		}
	}
}

type IssueListDelegate struct {
	width int
}

func NewIssueListDelegate(width int) IssueListDelegate {
	return IssueListDelegate{width: width}
}

func (d IssueListDelegate) Height() int                               { return 1 }
func (d IssueListDelegate) Spacing() int                              { return 0 }
func (d IssueListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d IssueListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	issue, ok := listItem.(ListIssue)
	if !ok {
		return
	}

	isSelected := index == m.Index()
	cols := getTableColumns(d.width)

	row := renderRow(issue, isSelected, cols)
	io.WriteString(w, row)
}

func renderRow(issue ListIssue, isSelected bool, cols []TableColumn) string {
	var parts []string

	for _, col := range cols {
		value := getColumnValue(col, issue)
		colWidth := col.width
		if colWidth == 0 {
			colWidth = 1
		}

		style := lipgloss.NewStyle().Width(int(colWidth))
		if isSelected {
			style = style.Background(styles.SelectedBackground).Bold(true)
		}

		truncated := truncate.StringWithTail(value, colWidth, "...")
		parts = append(parts, style.Render(truncated))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

func getColumnValue(col TableColumn, issue ListIssue) string {
	switch col.key {
	case "id":
		return issue.ID
	case "title":
		return issue.Title()
	case "status":
		return string(issue.Issue.Status)
	case "type":
		return string(issue.Issue.IssueType)
	default:
		return ""
	}
}
