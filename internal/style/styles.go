package style

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

var (
	Primary   = compat.AdaptiveColor{Light: lipgloss.Color("#5A56E0"), Dark: lipgloss.Color("#7571F9")}
	Secondary = compat.AdaptiveColor{Light: lipgloss.Color("#02BA84"), Dark: lipgloss.Color("#02BF87")}

	Success = compat.AdaptiveColor{Light: lipgloss.Color("#02BA84"), Dark: lipgloss.Color("#02BF87")}
	Warning = compat.AdaptiveColor{Light: lipgloss.Color("#F59E0B"), Dark: lipgloss.Color("#F59E0B")}
	Error   = compat.AdaptiveColor{Light: lipgloss.Color("#FE5F86"), Dark: lipgloss.Color("#FE5F86")}

	PrimaryText   = compat.AdaptiveColor{Light: lipgloss.Color("#1A1A1A"), Dark: lipgloss.Color("#E0E0E0")}
	SecondaryText = compat.AdaptiveColor{Light: lipgloss.Color("#666666"), Dark: lipgloss.Color("#999999")}
	FaintText     = compat.AdaptiveColor{Light: lipgloss.Color("#999999"), Dark: lipgloss.Color("#666666")}

	PrimaryBorder   = compat.AdaptiveColor{Light: lipgloss.Color("#5A56E0"), Dark: lipgloss.Color("#7571F9")}
	SecondaryBorder = compat.AdaptiveColor{Light: lipgloss.Color("#CCCCCC"), Dark: lipgloss.Color("#444444")}

	SelectedBackground = compat.AdaptiveColor{Light: lipgloss.Color("#E8E8E8"), Dark: lipgloss.Color("#333333")}
)

const (
	ListViewRatio     = 52 // Percentage of total width allocated to the list view
	LabelWidth        = 14
	MarginBottomSmall = 1
)

var DefaultBorder = lipgloss.ThickBorder()

var (
	HeaderStyle      = lipgloss.NewStyle().Foreground(Primary).Padding(0, 1).Bold(true)
	HeaderTitleStyle = lipgloss.NewStyle().Foreground(Primary).Bold(true).Padding(0)
)

var ContainerStyle = lipgloss.NewStyle().
	Border(DefaultBorder, true, false, false, false).
	BorderForeground(SecondaryBorder).
	Padding(1)

var ModalContainerStyle = lipgloss.NewStyle().
	Border(DefaultBorder).
	BorderForeground(PrimaryBorder).
	Padding(2, 3)

var DetailsContainerStyle = lipgloss.NewStyle().
	Border(DefaultBorder, true, false, false, true).
	BorderForeground(SecondaryBorder).
	Padding(1)

var (
	RowStyle       = lipgloss.NewStyle().MarginBottom(MarginBottomSmall)
	TitleStyle     = lipgloss.NewStyle().Foreground(Primary).Bold(true)
	LabelStyle     = lipgloss.NewStyle().Foreground(SecondaryText)
	ValueStyle     = lipgloss.NewStyle().Foreground(PrimaryText)
	IssueTypeStyle = lipgloss.NewStyle().Foreground(Primary).Bold(true)
)

var (
	FilterStyle       = lipgloss.NewStyle().Foreground(Primary).Bold(true).Padding(0, 1)
	FilterInputStyle  = lipgloss.NewStyle().Foreground(PrimaryText).Padding(0, 1)
	FilterPromptStyle = lipgloss.NewStyle().Foreground(Primary).Bold(true)
)

func StatusStyle(status string) lipgloss.Style {
	style := lipgloss.NewStyle().Bold(true)
	switch status {
	case "open":
		return style.Foreground(Secondary)
	case "closed":
		return style.Foreground(FaintText)
	case "in_progress":
		return style.Foreground(Warning)
	case "blocked":
		return style.Foreground(Error)
	default:
		return style.Foreground(SecondaryText)
	}
}

func HighlightKey(key string) string {
	return lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true).
		Padding(0, 1).
		Render(key)
}

var (
	TextStyle   = lipgloss.NewStyle().Foreground(PrimaryText)
	BorderStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(SecondaryBorder)
	ErrorStyle  = lipgloss.NewStyle().Foreground(Error).Bold(true)
	HelpStyle   = lipgloss.NewStyle().Align(lipgloss.Center).Foreground(Secondary)
)

var SecondaryColor = lipgloss.Color("#02BA84")
