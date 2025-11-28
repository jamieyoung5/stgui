package stgui

import "github.com/jamieyoung5/stgui/keyboard"

var DefaultDirectionalControls = map[string]string{
	keyboard.UpArrowKey:    "UP",
	keyboard.DownArrowKey:  "DOWN",
	keyboard.LeftArrowKey:  "LEFT",
	keyboard.RightArrowKey: "RIGHT",
}

type Cursor struct {
	widget Widget
	grid   *Grid
	Row    int
	Col    int

	Controls       map[string]string
	selectedColour string
}

func NewCursor(widget Widget, grid *Grid, row, col int, controls map[string]string) *Cursor {
	return &Cursor{
		widget:   widget,
		grid:     grid,
		Row:      row,
		Col:      col,
		Controls: controls,
	}
}

func (c *Cursor) Select(input string) (screen *Screen, exit bool) {
	return c.widget.Select(c, input)
}

func (c *Cursor) Up() {
	if c.Row > 0 {
		c.grid.Cells[c.Row][c.Col].Selected = false
		c.Row--
		c.grid.Cells[c.Row][c.Col].Selected = true
	}
}

func (c *Cursor) Down() {
	if c.Row < c.grid.height-1 {
		c.grid.Cells[c.Row][c.Col].Selected = false
		c.Row++
		c.grid.Cells[c.Row][c.Col].Selected = true
	}
}

func (c *Cursor) Right() {
	if c.Col < c.grid.width-1 {
		c.grid.Cells[c.Row][c.Col].Selected = false
		c.Col++
		c.grid.Cells[c.Row][c.Col].Selected = true
	}
}

func (c *Cursor) Left() {
	if c.Col > 0 {
		c.grid.Cells[c.Row][c.Col].Selected = false
		c.Col--
		c.grid.Cells[c.Row][c.Col].Selected = true
	}
}
