package stgui

import (
	"strings"

	"github.com/jamieyoung5/gostrc/strutil"
)

// Widget is a top-level piece of a screen. It draws itself and decides what each
// input means. widgets.Container is the usual one, backed by a grid.
type Widget interface {
	// Size in (columns, rows).
	Size() (width, height int)
	// Render draws the widget, lines separated by "\n".
	Render() string
	// Select handles one input. Return a screen to navigate to it; return exit
	// to close this screen, which quits if it's the last one open.
	Select(cursor *Cursor, input string) (screen *Screen, exit bool)
}

// Screen is one terminal-full: widgets laid out in rows, plus the cursors moving
// through them.
type Screen struct {
	// Persist keeps this screen on the stack when you navigate away, so the new
	// one goes on top and closing it lands back here.
	Persist bool
	Cursors []*Cursor

	widgets [][]Widget
}

func NewScreen(cursors []*Cursor, widgets [][]Widget) *Screen {
	return &Screen{
		Cursors: cursors,
		widgets: widgets,
	}
}

// ActiveCursor is the first cursor that isn't nil.
func (s *Screen) ActiveCursor() *Cursor {
	for _, c := range s.Cursors {
		if c != nil {
			return c
		}
	}
	return nil
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
