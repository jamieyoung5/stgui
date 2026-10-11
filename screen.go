package stgui

import (
	"strings"

	"github.com/jamieyoung5/gostrc/strutil"
)

// Widget is a top-level piece of a screen. It draws itself and decides what each
// input means. widgets.Container is the usual one, backed by a grid.
type Widget interface {
	// Render draws the widget, lines separated by "\n".
	Render() string
	// Select handles one input and returns what the app should do next.
	// Stay() carries on with the current screen.
	Select(cursor *Cursor, input Input) Action
}

// Screen is one terminal-full: widgets laid out in rows, plus the cursor moving
// through them.
type Screen struct {
	// Cursor is where input goes. A screen without one draws but ignores
	// every key.
	Cursor *Cursor

	widgets [][]Widget
}

// NewScreen lays rows of widgets out top to bottom, with cursor taking the
// input.
func NewScreen(cursor *Cursor, rows ...[]Widget) *Screen {
	return &Screen{
		Cursor:  cursor,
		widgets: rows,
	}
}

// Row is one row of a screen: widgets sitting side by side.
func Row(widgets ...Widget) []Widget {
	return widgets
}

// Render draws the lot: widgets sharing a row go side by side, rows get a blank
// line between them.
func (s *Screen) Render() string {
	rows := make([]string, 0, len(s.widgets))

	for _, widgetRow := range s.widgets {
		if len(widgetRow) == 0 {
			continue
		}

		// SideBySide terminates every line including the last, so drop that
		// one - otherwise joining rows gives two blank lines, not one.
		block := strutil.SideBySide(4, renderWidgets(widgetRow)...)
		rows = append(rows, strings.TrimSuffix(block, "\n"))
	}

	return strings.Join(rows, "\n\n")
}

func renderWidgets(widgets []Widget) []string {
	renderedElements := make([]string, 0, len(widgets))
	for _, elem := range widgets {
		if elem == nil {
			continue
		}
		renderedElements = append(renderedElements, elem.Render())
	}

	return renderedElements
}
