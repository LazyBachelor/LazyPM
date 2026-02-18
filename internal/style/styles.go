package style

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	PrimaryColor   = lipgloss.Color("6")
	SecondaryColor = lipgloss.Color("2")
	AccentColor    = lipgloss.Color("7")
	TextColor      = lipgloss.Color("15")

	BorderColor = lipgloss.Color("8")
)

var (
	AppStyle = lipgloss.NewStyle().Padding(1, 2).Foreground(TextColor)
)

var (
	DefaultBorder = lipgloss.NormalBorder()
	BorderStyle   = lipgloss.NewStyle().Border(DefaultBorder).BorderForeground(BorderColor)
)

var (
	TitleStyle = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)
	TextStyle  = lipgloss.NewStyle().Foreground(TextColor)
	HelpStyle  = lipgloss.NewStyle().Align(lipgloss.Center).Foreground(AccentColor)
)
