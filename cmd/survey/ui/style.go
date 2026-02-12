package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func theme() *huh.Theme {
	theme := huh.ThemeBase16()

	theme.Focused.Base = lipgloss.NewStyle().Align(lipgloss.Center)

	return theme
}
