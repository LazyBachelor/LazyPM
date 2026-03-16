package styles

import "github.com/charmbracelet/lipgloss"

var (
	Primary   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	Secondary = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}

	Success = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	Warning = lipgloss.AdaptiveColor{Light: "#F59E0B", Dark: "#F59E0B"}
	Error   = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}

	PrimaryText   = lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#E0E0E0"}
	SecondaryText = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"}
	FaintText     = lipgloss.AdaptiveColor{Light: "#999999", Dark: "#666666"}

	PrimaryBorder   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	SecondaryBorder = lipgloss.AdaptiveColor{Light: "#CCCCCC", Dark: "#444444"}

	SelectedBackground = lipgloss.AdaptiveColor{Light: "#E8E8E8", Dark: "#333333"}
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
