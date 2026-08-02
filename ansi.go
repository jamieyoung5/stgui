package stgui

import (
	"strings"
	"unicode/utf8"
)

// VisibleWidth is how many columns s takes up, skipping over escape sequences.
// Every rune counts as one column, so double-width CJK comes out short.
func VisibleWidth(s string) int {
	width := 0

	for i := 0; i < len(s); {
		if seq := ansiSequenceAt(s[i:]); seq != "" {
			i += len(seq)
			continue
		}

		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
		width++
	}

	return width
}

func maxVisibleWidth(lines []string) int {
	widest := 0
	for _, line := range lines {
		if w := VisibleWidth(line); w > widest {
			widest = w
		}
	}

	return widest
}

// truncateVisible cuts line down to width columns. Escapes are kept, since they
// cost no columns, and a reset goes on the end if we cut through a style.
func truncateVisible(line string, width int) string {
	if VisibleWidth(line) <= width {
		return line
	}

	var (
		sb      strings.Builder
		visible int
		styled  bool
	)

	for i := 0; i < len(line); {
		if seq := ansiSequenceAt(line[i:]); seq != "" {
			sb.WriteString(seq)
			styled = true
			i += len(seq)
			continue
		}

		if visible == width {
			break
		}

		r, size := utf8.DecodeRuneInString(line[i:])
		sb.WriteRune(r)
		visible++
		i += size
	}

	if styled {
		sb.WriteString(resetStyle)
	}

	return sb.String()
}

// ansiSequenceAt returns the escape sequence at the start of s, or "" if there
// isn't one there.
func ansiSequenceAt(s string) string {
	if len(s) < 2 || s[0] != '\x1b' || s[1] != '[' {
		return ""
	}

	for i := 2; i < len(s); i++ {
		if s[i] >= '@' && s[i] <= '~' {
			return s[:i+1]
		}
	}

	return s
}
