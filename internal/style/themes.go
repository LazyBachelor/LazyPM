package style

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func HuhCenterTheme() *huh.Styles {
	theme := huh.ThemeBase16(true)

	theme.Focused.Base = lipgloss.NewStyle().Align(lipgloss.Center)

	return theme
}
