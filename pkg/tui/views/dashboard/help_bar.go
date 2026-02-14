package dashboard

import (
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

type HelpBar struct {
	keyMap  DashboardKeyMap
	showAll bool
	width   int
}

func NewHelpBar(keyMap DashboardKeyMap) HelpBar {
	return HelpBar{keyMap: keyMap}
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
	keys := []string{
		styles.HighlightKey("↑/k") + " up",
		styles.HighlightKey("↓/j") + " down",
		styles.HighlightKey("a") + " add",
		styles.HighlightKey("e/d/s") + " edit",
		styles.HighlightKey("x") + " delete",
		styles.HighlightKey("q") + " quit",
		styles.HighlightKey("?") + " help",
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

	rows := []string{
		renderRow("↑/k", "up", "enter", "view issue"),
		renderRow("↓/j", "down", "b", "back to list"),
		renderRow("a", "add issue", "e", "edit title"),
		renderRow("d", "edit desc", "s", "change status"),
		renderRow("x", "delete issue", "?", "help"),
		renderRow("q", "quit", "", ""),
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
