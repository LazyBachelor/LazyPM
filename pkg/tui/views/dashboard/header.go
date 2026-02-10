package dashboard

import (
	"github.com/LazyBachelor/LazyPM/pkg/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

type Header struct {
	title string
}

func NewHeader(title string) Header {
	return Header{
		title: title,
	}
}

func (h Header) View(width int) string {
	title := styles.HeaderTitleStyle.Render(h.title)

	return lipgloss.PlaceHorizontal(
		width,
		lipgloss.Left,
		title,
		lipgloss.WithWhitespaceChars("─"),
		lipgloss.WithWhitespaceForeground(styles.Primary),
	)
}

func (h Header) Height() int {
	return 2
}
