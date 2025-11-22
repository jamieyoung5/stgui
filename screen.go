package stgui

import "strings"

type Widget interface {
	Size() (width, height int)
	Print(cursor *Cursor)
	Render(cursor *Cursor) string
	Select(cursor *Cursor, input string) (screen *Screen, exit bool)
}

type Element struct {
	cursor *Cursor
	widget Widget
}

func (vc *Element) Render() string {
	return vc.widget.Render(vc.cursor)
}

type Screen struct {
	Persist bool
	Cursors []*Cursor

	elements [][]*Element
}

func (s *Screen) Render() string {
	var builder strings.Builder

	for _, row := range s.elements {
		var items []string
		for _, componentNode := range row {
			items = append(items, componentNode.Render())
		}

		// Draw each row of Components side by side, with a specified number of spaces in between
		builder.WriteString(sideBySide(items, 4))
		builder.WriteString("\n\n") // Add spacing between rows
	}

	return builder.String()
}

func (s *Screen) SelectElement(cursor *Cursor, input string) (screen *Screen, exit bool) {
	return s.elements[cursor.gridX][cursor.gridY].widget.Select(cursor, input)
}
