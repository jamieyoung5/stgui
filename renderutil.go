package stgui

import (
	"fmt"
	"strings"
)

func Render(l Layout) string {
	if l == nil {
		return ""
	}
	return strings.Join(l.RenderLines(), "\n")
}

// MaxLength finds the length of the longest string in a slice.
func MaxLength(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	return maxWidth
}

func GenerateVerticalDivider(
	size int,
	symbols *Symbols) string {

	squareDivider := CreateLine(size*3, symbols.VerticalDivider)
	divider := squareDivider

	for i := 0; i < size-2; i++ {
		divider += symbols.CrossSectionDivider + squareDivider + symbols.VerticalDivider
	}

	return divider + symbols.CrossSectionDivider + squareDivider + "\n"
}

func CreateLine(length int, symbol string) (line string) {
	for range length {
		line += symbol
	}

	return line
}

// PadBlock pads a text block (lines) to a target height and width.
// It right-pads with spaces and bottom-pads with empty lines.
func PadBlock(lines []string, height, width int, noValue string) []string {
	if len(lines) == 0 {
		if height == 0 && width == 0 {
			return nil
		}
		if height == 0 {
			height = 1
		}
		if width == 0 {
			width = len(noValue)
		}
		lines = []string{noValue}
	}

	padded := make([]string, height)
	emptyLine := strings.Repeat(" ", width)

	for i := 0; i < height; i++ {
		if i < len(lines) {
			padded[i] = fmt.Sprintf("%-*s", width, lines[i])
		} else {
			padded[i] = emptyLine
		}
	}
	return padded
}

func sideBySide(blocks []string, spacing int) string {
	// Split each block into lines
	linesPerBlock := make([][]string, len(blocks))
	maxLines := 0
	for i, block := range blocks {
		lines := strings.Split(block, "\n")
		linesPerBlock[i] = lines
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}

	// Build the result with each line containing parts from corresponding lines of each block
	var result strings.Builder
	for i := 0; i < maxLines; i++ {
		for _, lines := range linesPerBlock {
			if i < len(lines) {
				result.WriteString(lines[i])
			}
			// Add extra spacing between blocks
			result.WriteString(strings.Repeat(" ", spacing))
		}
		result.WriteString("\n")
	}

	return result.String()
}
