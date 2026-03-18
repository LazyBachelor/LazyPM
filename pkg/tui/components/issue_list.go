package components

import (
	"context"
	"fmt"
	"io"
	"sort"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/muesli/reflow/truncate"
)

type IssueList struct {
	list              list.Model
	app               *app.App
	width             int
	height            int
	highlightSelected bool
}

type ListIssue struct {
	models.Issue
}

func (l ListIssue) Title() string       { return l.Issue.Title }
func (l ListIssue) Description() string { return l.Issue.Description }
func (l ListIssue) FilterValue() string { return l.Issue.ID + " " + l.Issue.Title }

type tableColumn struct {
	width uint
	label string
	key   string
}

func getTableColumns(width int) []tableColumn {
	switch {
	case width < 45:
		return []tableColumn{
			{width: 10, label: "ID", key: "id"},
			{width: uint(width - 10), label: "TITLE", key: "title"},
		}
	case width < 60:
		return []tableColumn{
			{width: 10, label: "ID", key: "id"},
			{width: 20, label: "TITLE", key: "title"},
			{width: 15, label: "STATUS", key: "status"},
		}
	case width < 75:
		return []tableColumn{
			{width: 12, label: "ID", key: "id"},
			{width: 20, label: "TITLE", key: "title"},
			{width: 15, label: "STATUS", key: "status"},
			{width: 10, label: "TYPE", key: "type"},
			{width: 15, label: "PRIORITY", key: "priority"},
		}
	default:
		return []tableColumn{
			{width: 12, label: "ID", key: "id"},
			{width: 20, label: "TITLE", key: "title"},
			{width: 15, label: "STATUS", key: "status"},
			{width: 11, label: "TYPE", key: "type"},
			{width: 14, label: "PRIORITY", key: "priority"},
			{width: 12, label: "ASSIGNEE", key: "assignee"},
		}
	}
}

// PriorityCodeName returns the human-readable name for a priority code.
// Exported for use by IssueDetail.
func PriorityCodeName(priority int) string {
	if name, ok := priorityCodeNames[priority]; ok {
		return name
	}
	return fmt.Sprintf("%d", priority)
}

var priorityCodeNames = map[int]string{
	0: "irrelevant",
	1: "low",
	2: "normal",
	3: "high",
	4: "critical",
}

func renderHeaders(cols []tableColumn) string {
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

// IssueInputs bundles the common text inputs used for issue title, creation, description, and assignee.
type IssueInputs struct {
	Title       textinput.Model
	CreateTitle textinput.Model
	Description textarea.Model
	Assignee    textinput.Model
}

// NewIssueInputs creates initialized inputs for issue title, new issue title, and description.
func NewIssueInputs() IssueInputs {
	ti := textinput.New()
	ti.Placeholder = "Issue title ..."
	ti.CharLimit = 256

	createTi := textinput.New()
	createTi.Placeholder = "New issue title ..."
	createTi.CharLimit = 256

	descTa := textarea.New()
	descTa.Placeholder = "Issue description..."
	descTa.SetWidth(56)
	descTa.SetHeight(8)

	assigneeTi := textinput.New()
	assigneeTi.Placeholder = "Assignee name..."
	assigneeTi.CharLimit = 64

	return IssueInputs{
		Title:       ti,
		CreateTitle: createTi,
		Description: descTa,
		Assignee:    assigneeTi,
	}
}

// ValidationFeedbackMsg is a shared message carrying validation results from the app.
type ValidationFeedbackMsg struct {
	Feedback models.ValidationFeedback
}

// ListenForValidation returns a command that waits for a validation feedback message
// on the given channel and wraps it in a ValidationFeedbackMsg.
func ListenForValidation(ch chan models.ValidationFeedback) tea.Cmd {
	return func() tea.Msg {
		feedback := <-ch
		return ValidationFeedbackMsg{Feedback: feedback}
	}
}

func NewIssueList(app *app.App, width, height int) IssueList {
	// NewIssueList creates an IssueList populated from the app.
	issues, err := app.Issues.SearchIssues(context.Background(), "", models.IssueFilter{})
	if err != nil {
		return IssueList{}
	}

	listIssues := []ListIssue{}
	for _, issue := range issues {
		listIssues = append(listIssues, ListIssue{Issue: *issue})
	}

	items := make([]list.Item, len(listIssues))
	for i, issue := range listIssues {
		items[i] = issue
	}

	l := list.New(items, newIssueListDelegate(width), width, height)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return IssueList{
		list:              l,
		app:               app,
		width:             width,
		height:            height,
		highlightSelected: true,
	}
}

// SetHighlightSelected controls whether the selected row is visually highlighted.
// Use false for lists that are not currently focused (e.g. non-focused Kanban columns).
func (l *IssueList) SetHighlightSelected(show bool) {
	l.highlightSelected = show
}

func NewIssueListFromIssues(app *app.App, issues []*models.Issue, width, height int) IssueList {
	// NewIssueListFromIssues creates an IssueList from a slice of issues.
	listIssues := make([]list.Item, len(issues))
	for i, issue := range issues {
		listIssues[i] = ListIssue{Issue: *issue}
	}
	l := list.New(listIssues, newIssueListDelegate(width), width, height)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	return IssueList{
		list:              l,
		app:               app,
		width:             width,
		height:            height,
		highlightSelected: true,
	}
}

func OpenAndInProgressOnly(issues []*models.Issue) []*models.Issue {
	// used to display open & in-progress issues in the first window in the dashboard
	out := make([]*models.Issue, 0, len(issues))
	for _, issue := range issues {
		if issue.Status == models.StatusOpen ||
			issue.Status == models.StatusInProgress ||
			issue.Status == models.StatusBlocked ||
			issue.Status == models.StatusReadyToSprint {
			out = append(out, issue)
		}
	}
	sortByPriorityDesc(out)
	return out
}

func ClosedOnly(issues []*models.Issue) []*models.Issue {
	// used to display issues in the second window in the dashboard
	out := make([]*models.Issue, 0, len(issues))
	for _, issue := range issues {
		if issue.Status == models.StatusClosed {
			out = append(out, issue)
		}
	}
	sortByPriorityDesc(out)
	return out
}

// StatusOnly returns issues that exactly match the given status, sorted by priority.
func StatusOnly(issues []*models.Issue, status models.Status) []*models.Issue {
	out := make([]*models.Issue, 0, len(issues))
	for _, issue := range issues {
		if issue.Status == status ||
			(status == models.StatusOpen && issue.Status == models.StatusReadyToSprint) {
			out = append(out, issue)
		}
	}
	sortByPriorityDesc(out)
	return out
}

func sortByPriorityDesc(issues []*models.Issue) {
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Priority > issues[j].Priority
	})
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
	l.list.SetDelegate(newIssueListDelegate(width))
}

func (l IssueList) View() string {
	// renders the list
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
		isSelected := l.highlightSelected && i == cursor
		if issue, ok := visibleItems[i].(ListIssue); ok {
			cols := getTableColumns(l.width)
			row := renderRow(issue, isSelected, cols)
			items = append(items, row)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// selectedItem returns the currently selected issue.
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

func (l *IssueList) SetIssues(issues []*models.Issue) tea.Cmd {
	listIssues := make([]list.Item, len(issues))
	for i, issue := range issues {
		listIssues[i] = ListIssue{Issue: *issue}
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

type issueListDelegate struct {
	width int
}

func newIssueListDelegate(width int) issueListDelegate {
	return issueListDelegate{width: width}
}

func (d issueListDelegate) Height() int                               { return 1 }
func (d issueListDelegate) Spacing() int                              { return 0 }
func (d issueListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d issueListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	issue, ok := listItem.(ListIssue)
	if !ok {
		return
	}

	isSelected := index == m.Index()
	cols := getTableColumns(d.width)
	row := renderRow(issue, isSelected, cols)
	io.WriteString(w, row)
}

func renderRow(issue ListIssue, isSelected bool, cols []tableColumn) string {
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

func getColumnValue(col tableColumn, issue ListIssue) string {
	switch col.key {
	case "id":
		return issue.ID
	case "title":
		return issue.Title()
	case "status":
		return string(issue.Issue.Status)
	case "type":
		return string(issue.Issue.IssueType)
	case "priority":
		return PriorityCodeName(issue.Issue.Priority)
	case "assignee":
		return issue.Assignee
	default:
		return ""
	}
}
