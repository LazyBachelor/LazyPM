package styles

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")).
			Bold(true).Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")).Padding(1)

	CommandStyle = lipgloss.NewStyle().Padding(1)
)
