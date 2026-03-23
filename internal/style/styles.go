package style

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

var (
	Primary = compat.AdaptiveColor{
		Light: lipgloss.Color("6"),
		Dark:  lipgloss.Color("6"),
	}
	Secondary = compat.AdaptiveColor{
		Light: lipgloss.Color("2"),
		Dark:  lipgloss.Color("2"),
	}

	Success = compat.AdaptiveColor{
		Light: lipgloss.Color("2"),
		Dark:  lipgloss.Color("2"),
	}
	Warning = compat.AdaptiveColor{
		Light: lipgloss.Color("3"),
		Dark:  lipgloss.Color("3"),
	}
	Error = compat.AdaptiveColor{
		Light: lipgloss.Color("9"),
		Dark:  lipgloss.Color("9"),
	}

	PrimaryText = compat.AdaptiveColor{
		Light: lipgloss.Color("7"),
		Dark:  lipgloss.Color("7"),
	}
	SecondaryText = compat.AdaptiveColor{
		Light: lipgloss.Color("8"),
		Dark:  lipgloss.Color("8"),
	}
	FaintText = compat.AdaptiveColor{
		Light: lipgloss.Color("8"),
		Dark:  lipgloss.Color("8"),
	}

	PrimaryBorder = compat.AdaptiveColor{
		Light: lipgloss.Color("8"),
		Dark:  lipgloss.Color("8"),
	}
	SecondaryBorder = compat.AdaptiveColor{
		Light: lipgloss.Color("8"),
		Dark:  lipgloss.Color("8"),
	}

	SelectedBackground = compat.AdaptiveColor{
		Light: lipgloss.Color("5"),
		Dark:  lipgloss.Color("5"),
	}
)

const (
	ListViewRatio     = 52 // Percentage of total width allocated to the list view
	LabelWidth        = 14
	MarginBottomSmall = 1
)

var DefaultBorder = lipgloss.ThickBorder()

var (
	HeaderStyle      = lipgloss.NewStyle().Foreground(Primary).Padding(0, 1).Bold(true)
	HeaderTitleStyle = lipgloss.NewStyle().Foreground(Primary).Bold(true).PaddingRight(1)
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
