package style

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

// BaseTheme wraps huh.ThemeBase for use as a huh.Theme
type BaseTheme struct{}

func (t BaseTheme) Theme(isDark bool) *huh.Styles {
	return huh.ThemeBase(isDark)
}

// Base16Theme wraps huh.ThemeBase16 for use as a huh.Theme
type Base16Theme struct{}

func (t Base16Theme) Theme(isDark bool) *huh.Styles {
	return huh.ThemeBase16(isDark)
}

// CenterTheme is a theme that centers form content while preserving ThemeBase16 styling
type CenterTheme struct{}

func (t CenterTheme) Theme(isDark bool) *huh.Styles {
	theme := huh.ThemeBase16(isDark)
	// Center the field content
	theme.Focused.Base = theme.Focused.Base.Align(lipgloss.Center)
	theme.Focused.Card = lipgloss.NewStyle() // Remove card border
	theme.Group.Base = lipgloss.NewStyle()   // Remove group padding

	return theme
}

// HuhCenterTheme returns a center-aligned theme for questionnaires
func HuhCenterTheme() huh.Theme {
	return CenterTheme{}
}

// HuhBaseTheme returns a base theme
func HuhBaseTheme() huh.Theme {
	return BaseTheme{}
}

// HuhBase16Theme returns a base16 theme
func HuhBase16Theme() huh.Theme {
	return Base16Theme{}
}
