package stgui

import "strings"

type Widget interface {
	GetDimensions() (height int, width int)
	Print(cursor *Cursor)
	Render(cursor *Cursor) string
	Select(cursor *Cursor, macro string) (screen *Screen, exit bool)
}

type Element struct {
	cursor *Cursor
	widget Widget
}

func (vc *Element) Render() string {
	return vc.widget.Render(vc.cursor)
}

type Screen struct {
	Elements [][]*Element
	Cursors  []*Cursor
	Persist  bool
}

func (s *Screen) Render() string {
	var builder strings.Builder

	for _, row := range s.Elements {
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
