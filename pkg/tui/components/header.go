package components

import (
	"charm.land/lipgloss/v2"
	"github.com/LazyBachelor/LazyPM/internal/style"
)

type Header struct {
	Title string
}

func NewHeader(title string) Header {
	return Header{Title: title}
}

func (h Header) View(width int) string {
	title := style.HeaderTitleStyle.Render(h.Title)

	return lipgloss.PlaceHorizontal(
		width,
		lipgloss.Left,
		title,
		lipgloss.WithWhitespaceChars("─"),
	)
}

func (h Header) Height() int {
	return 1
}
