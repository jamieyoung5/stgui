package stgui

import (
	"strings"
	"time"

	"github.com/jamieyoung5/stgui/term"
)

func render(canvas *InteractiveCanvas, quit chan bool) {
	screen := canvas.Screens.Peek().Serialize()
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Millisecond * 500)
			newScreen := canvas.Screens.Peek().Serialize()
			if newScreen != screen {
				term.Clear()
				canvas.Draw()
			}
		}
	}
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
