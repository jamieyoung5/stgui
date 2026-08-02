package widgets

import (
	"strings"

	"github.com/jamieyoung5/stgui"
)

// Text is a widget that draws a Renderable and ignores input - titles and status
// lines sitting above or below the interactive part of a screen.
type Text struct {
	Content stgui.Renderable
}

// NewText puts a Renderable, a Label say, straight onto a screen.
func NewText(content stgui.Renderable) *Text {
	return &Text{Content: content}
}

func (t *Text) lines() []string {
	if t.Content == nil {
		return nil
	}
	return t.Content.RenderLines()
}

func (t *Text) Size() (width, height int) {
	lines := t.lines()
	return maxWidth(lines), len(lines)
}

func (t *Text) Render() string {
	return strings.Join(t.lines(), "\n")
}

// Select does nothing. Text isn't interactive.
func (t *Text) Select(cursor *stgui.Cursor, input string) (*stgui.Screen, bool) {
	return nil, false
}

// maxWidth of the longest line.
func maxWidth(lines []string) int {
	widest := 0
	for _, line := range lines {
		if w := stgui.VisibleWidth(line); w > widest {
			widest = w
		}
	}

	return widest
}
