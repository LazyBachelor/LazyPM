package style

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	PrimaryColor   = lipgloss.AdaptiveColor{Light: "#007acc", Dark: "#1e90ff"}
	SecondaryColor = lipgloss.AdaptiveColor{Light: "#ff6f61", Dark: "#ff6347"}
	AccentColor    = lipgloss.AdaptiveColor{Light: "#6a5acd", Dark: "#9370db"}
	Background     = lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#1e1e1e"}
	TextColor      = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}
)

var (
	AppStyle = lipgloss.NewStyle().Padding(1, 2).Background(Background).Foreground(TextColor)
)

var (
	DefaultBorder = lipgloss.NormalBorder()
	BorderStyle   = lipgloss.NewStyle().Border(DefaultBorder).BorderForeground(PrimaryColor)
)

var (
	TitleStyle       = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)
	DescriptionStyle = lipgloss.NewStyle().Foreground(TextColor).Italic(true)
	DetailStyle      = lipgloss.NewStyle().Foreground(SecondaryColor)
	HelpStyle        = lipgloss.NewStyle().Align(lipgloss.Center).Foreground(AccentColor)
)
