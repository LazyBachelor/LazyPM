package style

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func HuhCenterTheme() *huh.Theme {
	theme := huh.ThemeBase16()

	theme.Focused.Base = lipgloss.NewStyle().Align(lipgloss.Center)

	return theme
}
