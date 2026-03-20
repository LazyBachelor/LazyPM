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
		runes := []rune(text)
		if len(runes) > 0 {
			firstChar := string(runes[0])
			if lipgloss.Width(firstChar) <= maxWidth {
				return firstChar
			}
		}
		return ""
	}

	runes := []rune(text)
	lastSafe := 0
	for i := range runes {
		candidate := string(runes[:i+1])
		if lipgloss.Width(candidate)+ellipsisWidth > maxWidth {
			break
		}
		lastSafe = i + 1
	}

	current := string(runes[:lastSafe])
	return current + ellipsis
}
