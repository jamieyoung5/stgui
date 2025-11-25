package stgui

import (
	"strings"

	"github.com/jamieyoung5/gostrc/strutil"
)

type Widget interface {
	Size() (width, height int)
	Render(cursor *Cursor) string
	Select(cursor *Cursor, input string) (screen *Screen, exit bool)
}

type Element struct {
	cursor *Cursor
	widget Widget
}

func NewElement(cursor *Cursor, widget Widget) *Element {
	return &Element{
		cursor: cursor,
		widget: widget,
	}
}

func (vc *Element) Render() string {
	return vc.widget.Render(vc.cursor)
}

type Screen struct {
	Persist bool
	Cursors []*Cursor

	elements [][]*Element
}

func NewScreen(cursors []*Cursor, elements [][]*Element) *Screen {
	return &Screen{
		Cursors:  cursors,
		elements: elements,
	}
}

func (s *Screen) SelectElement(cursor *Cursor, input string) (screen *Screen, exit bool) {
	return s.elements[cursor.gridX][cursor.gridY].widget.Select(cursor, input)
}

func (s *Screen) Render() string {
	var builder strings.Builder

	for _, elemSet := range s.elements {
		elems := renderElements(elemSet)

		// Draw each row of Components side by side, with a specified number of spaces in between
		builder.WriteString(strutil.SideBySide(4, elems...))
		builder.WriteString("\n\n") // Add spacing between rows
	}

	return builder.String()
}

func renderElements(elements []*Element) []string {
	var renderedElements []string
	for _, elem := range elements {
		renderedElements = append(renderedElements, elem.Render())
	}

	return renderedElements
}
