package modal

import (
	"context"
	"io"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/app"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/style"
	"github.com/LazyBachelor/LazyPM/pkg/tui/msgs"
)

// struct for displaying an issue's dependencies and allows adding/removing new ones using keys 'a' and 'x'
type DependenciesModal struct {
	BaseModal
	app     *app.App
	issueID string
	deps    []*models.Issue
	width   int
	height  int
}

func NewDependenciesModal(a *app.App, issueID string) *DependenciesModal {
	return &DependenciesModal{
		BaseModal: NewBaseModal(ModalManageDependencies, TypeCustom),
		app:       a,
		issueID:   issueID,
		deps:      nil,
		width:     60,
		height:    20,
	}
}

func (d *DependenciesModal) Activate() tea.Cmd {
	d.BaseModal.activate()
	d.refreshDeps()
	return nil
}

func (d *DependenciesModal) Deactivate() {
	d.BaseModal.deactivate()
}

func (d *DependenciesModal) SetIssueID(issueID string) {
	d.issueID = issueID
	d.refreshDeps()
}

func (d *DependenciesModal) RefreshDeps() {
	d.refreshDeps()
}

func (d *DependenciesModal) refreshDeps() {
	if d.app == nil || d.issueID == "" {
		d.deps = nil
		return
	}
	deps, _ := d.app.Issues.GetDependencies(context.Background(), d.issueID)
	d.deps = deps
}

func (d *DependenciesModal) Update(msg tea.Msg) (tea.Cmd, bool) {
	if !d.IsActive() {
		return nil, false
	}

	// let msg pass through to the view
	if _, ok := msg.(msgs.DependencyAddRequestedMsg); ok {
		return nil, false
	}
	if _, ok := msg.(msgs.DependencyRemoveRequestedMsg); ok {
		return nil, false
	}
	if completed, ok := msg.(ModalCompletedMsg); ok && (completed.ModalID == ModalSelectDependency || completed.ModalID == ModalSelectRemoveDependency) {
		return nil, false
	}

	// refresh list when adding/removing dependencies
	if _, ok := msg.(msgs.DependencyAddedMsg); ok {
		d.refreshDeps()
		return nil, false
	}
	if _, ok := msg.(msgs.DependencyRemovedMsg); ok {
		d.refreshDeps()
		return nil, false
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			d.Deactivate()
			return func() tea.Msg { return ModalCancelledMsg{ModalID: d.ID()} }, true
		case "a":
			return func() tea.Msg { return msgs.DependencyAddRequestedMsg{IssueID: d.issueID} }, true
		case "x":
			return func() tea.Msg { return msgs.DependencyRemoveRequestedMsg{IssueID: d.issueID} }, true
		}
	}

	return nil, true // consume other keys when this modal is active (e.g. help key)
}

func (d *DependenciesModal) View() string {
	if d.width < 5 {
		return ""
	}

	boxWidth := max(min(60, d.width-4), 1)

	title := style.LabelStyle.Render("Dependencies for issue " + d.issueID)
	helpRow := lipgloss.NewStyle().Foreground(style.FaintText).Render("a = add, x = remove, esc = close")

	var depsContent string
	if len(d.deps) == 0 {
		depsContent = style.ValueStyle.Render("No dependencies.")
	} else {
		var lines []string
		for _, dep := range d.deps {
			if dep == nil {
				continue
			}
			line := "  " + dep.ID
			if dep.Title != "" {
				line += " — " + dep.Title
			}
			lines = append(lines, style.ValueStyle.Render(line))
		}
		depsContent = lipgloss.JoinVertical(lipgloss.Left, lines...)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		depsContent,
		"",
		helpRow,
	)

	return style.ModalContainerStyle.
		Width(boxWidth).
		Render(content)
}

func (d *DependenciesModal) SetSize(width, height int) {
	d.BaseModal.SetSize(width, height)
	d.width = width
	d.height = height
}

// dependencyListItem implements list.Item for the add-dependency list.
type dependencyListItem struct {
	Issue *models.Issue
}

func (i dependencyListItem) Title() string {
	if i.Issue == nil {
		return ""
	}
	return i.Issue.ID
}

func (i dependencyListItem) Description() string {
	if i.Issue == nil {
		return ""
	}
	return i.Issue.Title
}

func (i dependencyListItem) FilterValue() string {
	if i.Issue == nil {
		return ""
	}
	return i.Issue.ID + " " + i.Issue.Title
}

// DependencyListModal is a list-based modal for selecting an issue to add as a dependency.
// Use j/k and up/down to navigate, Enter to select.
type DependencyListModal struct {
	BaseModal
	list   list.Model
	width  int
	height int
}

type dependencyListDelegate struct {
	width int
}

func (d dependencyListDelegate) Height() int                               { return 1 }
func (d dependencyListDelegate) Spacing() int                              { return 0 }
func (d dependencyListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d dependencyListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	li, ok := item.(dependencyListItem)
	if !ok || li.Issue == nil {
		return
	}
	text := li.Issue.ID
	if li.Issue.Title != "" {
		text += " — " + li.Issue.Title
	}
	if index == m.Index() {
		text = style.HighlightKey("> ") + lipgloss.NewStyle().Bold(true).Render(text)
	} else {
		text = "  " + text
	}
	io.WriteString(w, lipgloss.NewStyle().Width(d.width-2).Render(text))
}

// NewDependencyListModal creates a modal with a list of issues to add as dependency.
func NewDependencyListModal(issues []*models.Issue, width, height int) *DependencyListModal {
	items := make([]list.Item, 0, len(issues))
	for _, iss := range issues {
		if iss != nil {
			items = append(items, dependencyListItem{Issue: iss})
		}
	}

	if width < 40 {
		width = 50
	}
	if height < 10 {
		height = 15
	}

	l := list.New(items, dependencyListDelegate{width: width}, width-4, height-6)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return &DependencyListModal{
		BaseModal: NewBaseModal(ModalSelectDependency, TypeCustom),
		list:      l,
		width:     width,
		height:    height,
	}
}

func (m *DependencyListModal) Activate() tea.Cmd {
	m.BaseModal.activate()
	return nil
}

func (m *DependencyListModal) Deactivate() {
	m.BaseModal.deactivate()
}

func (m *DependencyListModal) Update(msg tea.Msg) (tea.Cmd, bool) {
	if !m.IsActive() {
		return nil, false
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			m.Deactivate()
			return func() tea.Msg { return ModalCancelledMsg{ModalID: m.ID()} }, true
		case "enter":
			item := m.list.SelectedItem()
			if li, ok := item.(dependencyListItem); ok && li.Issue != nil {
				m.Deactivate()
				return func() tea.Msg {
					return ModalCompletedMsg{
						ModalID: m.ID(),
						Value:   SelectResult{SelectedValue: li.Issue.ID},
					}
				}, true
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return cmd, true
}

func (m *DependencyListModal) View() string {
	boxWidth := max(min(60, m.width-4), 1)
	helpRow := lipgloss.NewStyle().Foreground(style.FaintText).Render("↑/k up, ↓/j down, enter = add, esc = cancel")

	content := lipgloss.JoinVertical(lipgloss.Left,
		style.LabelStyle.Render("Add dependency (select issue):"),
		"",
		m.list.View(),
		"",
		helpRow,
	)

	return style.ModalContainerStyle.
		Width(boxWidth).
		Render(content)
}

func (m *DependencyListModal) SetSize(width, height int) {
	m.BaseModal.SetSize(width, height)
	m.width = width
	m.height = height
	m.list.SetSize(width-4, height-6)
}

// DependencyRemoveListModal is a list-based modal for selecting a dependency to remove.
// Use j/k and up/down to navigate, Enter to remove.
type DependencyRemoveListModal struct {
	BaseModal
	list   list.Model
	width  int
	height int
}

// NewDependencyRemoveListModal creates a modal with a list of dependencies to remove.
func NewDependencyRemoveListModal(deps []*models.Issue, width, height int) *DependencyRemoveListModal {
	items := make([]list.Item, 0, len(deps))
	for _, dep := range deps {
		if dep != nil {
			items = append(items, dependencyListItem{Issue: dep})
		}
	}

	if width < 40 {
		width = 50
	}
	if height < 10 {
		height = 15
	}

	l := list.New(items, dependencyListDelegate{width: width}, width-4, height-6)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return &DependencyRemoveListModal{
		BaseModal: NewBaseModal(ModalSelectRemoveDependency, TypeCustom),
		list:      l,
		width:     width,
		height:    height,
	}
}

func (m *DependencyRemoveListModal) Activate() tea.Cmd {
	m.BaseModal.activate()
	return nil
}

func (m *DependencyRemoveListModal) Deactivate() {
	m.BaseModal.deactivate()
}

func (m *DependencyRemoveListModal) Update(msg tea.Msg) (tea.Cmd, bool) {
	if !m.IsActive() {
		return nil, false
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			m.Deactivate()
			return func() tea.Msg { return ModalCancelledMsg{ModalID: m.ID()} }, true
		case "enter":
			item := m.list.SelectedItem()
			if li, ok := item.(dependencyListItem); ok && li.Issue != nil {
				m.Deactivate()
				return func() tea.Msg {
					return ModalCompletedMsg{
						ModalID: m.ID(),
						Value:   SelectResult{SelectedValue: li.Issue.ID},
					}
				}, true
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return cmd, true
}

func (m *DependencyRemoveListModal) View() string {
	boxWidth := max(min(60, m.width-4), 1)
	helpRow := lipgloss.NewStyle().Foreground(style.FaintText).Render("↑/k up, ↓/j down, enter = remove, esc = cancel")

	content := lipgloss.JoinVertical(lipgloss.Left,
		style.LabelStyle.Render("Remove dependency (select issue):"),
		"",
		m.list.View(),
		"",
		helpRow,
	)

	return style.ModalContainerStyle.
		Width(boxWidth).
		Render(content)
}

func (m *DependencyRemoveListModal) SetSize(width, height int) {
	m.BaseModal.SetSize(width, height)
	m.width = width
	m.height = height
	m.list.SetSize(width-4, height-6)
}

// EligibleDependencyIssues returns issues that can be added as dependencies:
// excludes current issue, existing dependencies, and dependents.
func EligibleDependencyIssues(a *app.App, issueID string) []*models.Issue {
	if a == nil || issueID == "" {
		return nil
	}
	ctx := context.Background()
	all, err := a.Issues.SearchIssues(ctx, "", models.IssueFilter{})
	if err != nil {
		return nil
	}
	deps, _ := a.Issues.GetDependencies(ctx, issueID)
	dependents, _ := a.Issues.GetDependents(ctx, issueID)

	exclude := map[string]struct{}{issueID: {}}
	for _, d := range deps {
		if d != nil {
			exclude[d.ID] = struct{}{}
		}
	}
	for _, d := range dependents {
		if d != nil {
			exclude[d.ID] = struct{}{}
		}
	}

	var out []*models.Issue
	for _, iss := range all {
		if iss != nil {
			if _, skip := exclude[iss.ID]; !skip {
				out = append(out, iss)
			}
		}
	}
	return out
}
