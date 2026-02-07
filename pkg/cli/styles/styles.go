// Package styles defines the styling for the CLI output using the lipgloss library.
package styles

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().Bold(true).Padding(1)

	CommandStyle = lipgloss.NewStyle().Padding(1)
)
