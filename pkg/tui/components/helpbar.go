package components

import (
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
)

type ViewKind int

const (
	ViewIssues ViewKind = iota
	ViewKanban
)

type ShortItem struct {
	Key  string
	Desc string
}

type FullRow struct {
	LeftKey   string
	LeftDesc  string
	RightKey  string
	RightDesc string
}

type HelpBarConfig struct {
	ShortItems []ShortItem
	FullRows   []FullRow
}

type HelpBar struct {
	view    ViewKind
	config  HelpBarConfig
	showAll bool
	width   int
}

func NewHelpBar(view ViewKind) HelpBar {
	return HelpBar{view: view, config: helpBarConfig(view)}
}

func (h *HelpBar) SetWidth(width int) {
	h.width = width
}

func (h HelpBar) View() string {
	if h.width == 0 {
		return ""
	}
	if h.showAll {
		return h.fullHelp()
	}
	return h.shortHelp()
}

func (h HelpBar) shortHelp() string {
	keys := make([]string, 0, len(h.config.ShortItems))
	for _, item := range h.config.ShortItems {
		keys = append(keys, styles.HighlightKey(item.Key)+item.Desc+"  ")
	}
	content := lipgloss.JoinHorizontal(lipgloss.Left, keys...)
	return lipgloss.NewStyle().
		Border(lipgloss.Border{Top: "─"}, true, false, false, false).
		BorderForeground(styles.SecondaryBorder).
		Padding(0, 1).
		Width(h.width).
		Render(content)
}

func (h HelpBar) fullHelp() string {
	keyStyle := lipgloss.NewStyle().Width(8).Align(lipgloss.Right)
	descStyle := lipgloss.NewStyle().Width(12)

	renderHelpItem := func(key, desc string) string {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			keyStyle.Render(styles.HighlightKey(key)),
			" ",
			descStyle.Render(desc),
		)
	}

	renderRow := func(leftKey, leftDesc, rightKey, rightDesc string) string {
		leftItem := renderHelpItem(leftKey, leftDesc)
		rightItem := renderHelpItem(rightKey, rightDesc)
		return lipgloss.JoinHorizontal(lipgloss.Left, leftItem, "  ", rightItem)
	}

	rows := make([]string, 0, len(h.config.FullRows))
	for _, r := range h.config.FullRows {
		rows = append(rows, renderRow(r.LeftKey, r.LeftDesc, r.RightKey, r.RightDesc))
	}
	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return lipgloss.NewStyle().
		Border(lipgloss.Border{Top: "─"}, true, false, false, false).
		BorderForeground(styles.SecondaryBorder).
		Padding(0, 1).
		Width(h.width).
		Render(content)
}

func (h HelpBar) Height() int {
	if h.width == 0 {
		return 0
	}
	return lipgloss.Height(h.View())
}

func (h HelpBar) IsExpanded() bool {
	return h.showAll
}

func (h *HelpBar) ToggleHelp() {
	h.showAll = !h.showAll
}

func helpBarConfig(view ViewKind) HelpBarConfig {
	switch view {
	case ViewIssues:
		return HelpBarConfig{
			ShortItems: []ShortItem{
				{Key: "tab", Desc: "switch"},
				{Key: "v", Desc: "kanban"},
				{Key: "↑/k", Desc: "up"},
				{Key: "↓/j", Desc: "down"},
				{Key: "pgup/pgdn", Desc: "page"},
				{Key: "a", Desc: "add"},
				{Key: "c", Desc: "comment"},
				{Key: "e/d/s/p/t/A", Desc: "edit"},
				{Key: "x", Desc: "delete"},
				{Key: "S", Desc: "submit"},
				{Key: "q", Desc: "quit"},
				{Key: "?", Desc: "help"},
			},
			FullRows: []FullRow{
				{LeftKey: "tab", LeftDesc: "switch window", RightKey: "↑/k", RightDesc: "up"},
				{LeftKey: "enter", LeftDesc: "view issue", RightKey: "↓/j", RightDesc: "down"},
				{LeftKey: "pgup", LeftDesc: "page up", RightKey: "pgdn", RightDesc: "page down"},
				{LeftKey: "b", LeftDesc: "back to list", RightKey: "a", RightDesc: "add issue"},
				{LeftKey: "c", LeftDesc: "add comment", RightKey: "", RightDesc: ""},
				{LeftKey: "e", LeftDesc: "edit title", RightKey: "d", RightDesc: "edit description"},
				{LeftKey: "s", LeftDesc: "change status", RightKey: "p", RightDesc: "change priority"},
				{LeftKey: "t", LeftDesc: "change type", RightKey: "A", RightDesc: "change assignee"},
				{LeftKey: "x", LeftDesc: "delete issue", RightKey: "v", RightDesc: "kanban"},
				{LeftKey: "q", LeftDesc: "quit", RightKey: "S", RightDesc: "submit"},
				{LeftKey: "?", LeftDesc: "help", RightKey: "", RightDesc: ""},
			},
		}
	case ViewKanban:
		return HelpBarConfig{
			ShortItems: []ShortItem{
				{Key: "v", Desc: "list view"},
				{Key: "↑/k", Desc: "up"},
				{Key: "↓/j", Desc: "down"},
				{Key: "pgup/pgdn", Desc: "page"},
				{Key: "h/l", Desc: "column"},
				{Key: "←/→", Desc: "move"},
				{Key: "a", Desc: "add"},
				{Key: "e/d/s/p/t/A", Desc: "edit"},
				{Key: "x", Desc: "delete"},
				{Key: "q", Desc: "quit"},
				{Key: "S", Desc: "submit"},
				{Key: "?", Desc: "help"},
			},
			FullRows: []FullRow{
				{LeftKey: "v", LeftDesc: "list view", RightKey: "↑/k", RightDesc: "up"},
				{LeftKey: "enter", LeftDesc: "view issue", RightKey: "↓/j", RightDesc: "down"},
				{LeftKey: "pgup", LeftDesc: "page up", RightKey: "pgdn", RightDesc: "page down"},
				{LeftKey: "h/l", LeftDesc: "switch column", RightKey: "←/→", RightDesc: "move issue"},
				{LeftKey: "b", LeftDesc: "back to list", RightKey: "a", RightDesc: "add issue"},
				{LeftKey: "e", LeftDesc: "edit title", RightKey: "d", RightDesc: "edit description"},
				{LeftKey: "s", LeftDesc: "change status", RightKey: "p", RightDesc: "change priority"},
				{LeftKey: "t", LeftDesc: "change type", RightKey: "A", RightDesc: "change assignee"},
				{LeftKey: "x", LeftDesc: "delete issue", RightKey: "q", RightDesc: "quit"},
				{LeftKey: "S", LeftDesc: "submit", RightKey: "?", RightDesc: "help"},
			},
		}
	default:
		return HelpBarConfig{}
	}
}
