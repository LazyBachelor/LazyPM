package components

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/style"
)

type ViewKind int

const (
	ViewIssues ViewKind = iota
	ViewKanban
)

const (
	fullHelpKeyWidth  = 10
	fullHelpDescWidth = 18
	fullHelpColGap    = 4
	maxFullHelpCols   = 4
)

type HelpItem struct {
	Key  string
	Desc string
}

type HelpBarConfig struct {
	ShortItems []HelpItem
	FullItems  []HelpItem
}

type HelpBar struct {
	view    ViewKind
	config  HelpBarConfig
	showAll bool
	width   int
}

func NewHelpBar(view ViewKind) HelpBar {
	return HelpBar{
		view:   view,
		config: helpBarConfig(view),
	}
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
	items := make([]string, 0, len(h.config.ShortItems))
	for _, item := range h.config.ShortItems {
		items = append(
			items,
			style.HighlightKey(item.Key)+item.Desc+"  ",
		)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Left, items...)

	return lipgloss.NewStyle().
		Border(lipgloss.Border{Top: "─"}, true, false, false, false).
		BorderForeground(style.SecondaryBorder).
		Padding(0, 1).
		Width(h.width).
		Render(content)
}

func (h HelpBar) fullHelp() string {
	if len(h.config.FullItems) == 0 {
		return ""
	}

	keyStyle := lipgloss.NewStyle().Width(fullHelpKeyWidth).Align(lipgloss.Right)
	descStyle := lipgloss.NewStyle().Width(fullHelpDescWidth)
	cellStyle := lipgloss.NewStyle().Width(fullHelpKeyWidth + 1 + fullHelpDescWidth)

	renderItem := func(item HelpItem) string {
		return cellStyle.Render(
			keyStyle.Render(style.HighlightKey(item.Key)) + " " + descStyle.Render(item.Desc),
		)
	}

	innerWidth := max(h.width-2, 1)
	cellWidth := fullHelpKeyWidth + 1 + fullHelpDescWidth
	cols := fitHelpColumns(innerWidth, cellWidth, fullHelpColGap, maxFullHelpCols)
	rows := (len(h.config.FullItems) + cols - 1) / cols
	gap := strings.Repeat(" ", fullHelpColGap)

	var result []string
	for r := range rows {
		var cells []string
		for c := range cols {
			if idx := r + c*rows; idx < len(h.config.FullItems) {
				cells = append(cells, renderItem(h.config.FullItems[idx]))
			}
		}
		result = append(result, joinHorizontalWithGap(cells, gap))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, result...)
	return lipgloss.NewStyle().
		Border(lipgloss.Border{Top: "─"}, true, false, false, false).
		BorderForeground(style.SecondaryBorder).
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

func fitHelpColumns(
	availableWidth int,
	cellWidth int,
	gapWidth int,
	maxCols int,
) int {
	for cols := maxCols; cols >= 1; cols-- {
		neededWidth := cols*cellWidth + (cols-1)*gapWidth
		if neededWidth <= availableWidth {
			return cols
		}
	}

	return 1
}

func joinHorizontalWithGap(cells []string, gap string) string {
	if len(cells) == 0 {
		return ""
	}

	parts := make([]string, 0, len(cells)*2-1)
	for i, cell := range cells {
		if i > 0 {
			parts = append(parts, gap)
		}
		parts = append(parts, cell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

func helpBarConfig(view ViewKind) HelpBarConfig {
	switch view {
	case ViewIssues:
		return HelpBarConfig{
			ShortItems: []HelpItem{
				{Key: "v", Desc: "kanban"},
				{Key: "↑/k", Desc: "up"},
				{Key: "↓/j", Desc: "down"},
				{Key: "/", Desc: "search"},
				{Key: "a", Desc: "add"},
				{Key: "c", Desc: "comment"},
				{Key: "e/d/s/p/t/A/D", Desc: "edit"},
				{Key: "x", Desc: "delete"},
				{Key: "q", Desc: "quit"},
				{Key: "?", Desc: "help"},
			},
			FullItems: []HelpItem{
				{Key: "v", Desc: "kanban"},
				{Key: "↑/k", Desc: "up"},
				{Key: "↓/j", Desc: "down"},
				{Key: "/", Desc: "search"},
				{Key: "a", Desc: "add issue"},
				{Key: "c", Desc: "add comment"},
				{Key: "x", Desc: "delete issue"},
				{Key: "e", Desc: "edit title"},
				{Key: "d", Desc: "edit description"},
				{Key: "s", Desc: "change status"},
				{Key: "p", Desc: "change priority"},
				{Key: "t", Desc: "change type"},
				{Key: "A", Desc: "change assignee"},
				{Key: "D", Desc: "edit dependencies"},
				{Key: "q", Desc: "quit"},
				{Key: "?", Desc: "help"},
			},
		}

	case ViewKanban:
		return HelpBarConfig{
			ShortItems: []HelpItem{
				{Key: "v", Desc: "list view"},
				{Key: "↑/k", Desc: "up"},
				{Key: "↓/j", Desc: "down"},
				{Key: "/", Desc: "search"},
				{Key: "n", Desc: "new sprint"},
				{Key: "S", Desc: "select sprint"},
				{Key: "pgup/pgdn", Desc: "page"},
				{Key: "h/l", Desc: "column"},
				{Key: "←/→", Desc: "move"},
				{Key: "a", Desc: "add"},
				{Key: "e/d/s/p/t/A/D", Desc: "edit"},
				{Key: "x", Desc: "delete"},
				{Key: "q", Desc: "quit"},
				{Key: "?", Desc: "help"},
			},
			FullItems: []HelpItem{
				{Key: "v", Desc: "list view"},
				{Key: "h/l", Desc: "switch column"},
				{Key: "←/→", Desc: "move issue"},
				{Key: "↑/k", Desc: "up"},
				{Key: "↓/j", Desc: "down"},
				{Key: "/", Desc: "search"},
				{Key: "pgup", Desc: "page up"},
				{Key: "pgdn", Desc: "page down"},
				{Key: "a", Desc: "add issue"},
				{Key: "x", Desc: "delete issue"},
				{Key: "e", Desc: "edit title"},
				{Key: "d", Desc: "edit description"},
				{Key: "s", Desc: "change status"},
				{Key: "p", Desc: "change priority"},
				{Key: "t", Desc: "change type"},
				{Key: "A", Desc: "change assignee"},
				{Key: "D", Desc: "manage dependencies"},
				{Key: "S", Desc: "select sprint"},
				{Key: "n", Desc: "new sprint"},
				{Key: "q", Desc: "quit"},
				{Key: "?", Desc: "help"},
			},
		}

	default:
		return HelpBarConfig{}
	}
}
