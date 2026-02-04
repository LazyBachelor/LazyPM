package styles

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().Padding(3, 3)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Padding(0, 1).Bold(true).Border(lipgloss.NormalBorder())

	SelectedItemStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				Foreground(lipgloss.Color("2")).
				Padding(0, 0, 0, 1)

	ItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Padding(0, 0, 0, 2)

	IssueStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("8")).
			Padding(1)

	FocusedIssueStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("2")).
				Padding(1)
)
