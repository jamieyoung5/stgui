package stgui

import (
	"strings"

	"github.com/jamieyoung5/gostrc/strutil"
)

type Widget interface {
	Size() (width, height int)
	Render() string
	Select(cursor *Cursor, input string) (screen *Screen, exit bool)
}

type Screen struct {
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

// ActiveCursor returns the first non-nil cursor on the screen.
func (s *Screen) ActiveCursor() *Cursor {
	for _, c := range s.Cursors {
		if c != nil {
			return c
		}
	}
	return nil
}

func (s *Screen) Render() string {
	var builder strings.Builder

	for _, elemSet := range s.widgets {
		elems := renderWidgets(elemSet)

		// draw each row of Components side by side, with a specified number of spaces in between
		builder.WriteString(strutil.SideBySide(4, elems...))
		builder.WriteString("\n\n") // add spacing between rows
	}

	return builder.String()
}

func renderWidgets(widgets []Widget) []string {
	var renderedElements []string
	for _, elem := range widgets {
		renderedElements = append(renderedElements, elem.Render())
	}

	return renderedElements
}
