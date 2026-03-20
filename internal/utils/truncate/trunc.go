package truncate

import "charm.land/lipgloss/v2"

// TruncateToWidth trims the given text so that its rendered width does not
// exceed maxWidth. If truncation occurs and there is room, an ellipsis is
// appended to indicate that the text was shortened.
func TruncateToWidth(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	if lipgloss.Width(text) <= maxWidth {
		return text
	}

	const ellipsis = "…"
	ellipsisWidth := lipgloss.Width(ellipsis)
	if ellipsisWidth > maxWidth {
		// Not enough space even for an ellipsis; return empty.
		return ""
	}

	runes := []rune(text)
	current := ""
	for _, r := range runes {
		next := current + string(r)
		if lipgloss.Width(next)+ellipsisWidth > maxWidth {
			break
		}
		current = next
	}

	return current + ellipsis
}
