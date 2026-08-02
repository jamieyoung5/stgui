package stgui

import "github.com/jamieyoung5/stgui/keyboard"

// Navigation commands. Controls maps keys onto these and widgets act on the
// commands, not the keys, which is what keeps bindings swappable.
const (
	ControlUp    = "UP"
	ControlDown  = "DOWN"
	ControlLeft  = "LEFT"
	ControlRight = "RIGHT"
)

// DefaultDirectionalControls is just the arrow keys. Copy and extend it to add
// your own, e.g. "w": ControlUp.
var DefaultDirectionalControls = map[string]string{
	keyboard.UpArrowKey:    ControlUp,
	keyboard.DownArrowKey:  ControlDown,
	keyboard.LeftArrowKey:  ControlLeft,
	keyboard.RightArrowKey: ControlRight,
}

// Cursor is the focused cell of a grid. A screen's input goes through the
// cursor's widget, which decides what to make of it.
type Cursor struct {
	widget Widget
	grid   *Grid
	Row    int
	Col    int

	Controls map[string]string
}

// NewCursor focuses (row, col) and marks that cell selected. The position is
// clamped to the grid; nil controls means the arrow keys.
func NewCursor(widget Widget, grid *Grid, row, col int, controls map[string]string) *Cursor {
	if controls == nil {
		controls = DefaultDirectionalControls
	}

	c := &Cursor{
		widget:   widget,
		grid:     grid,
		Row:      row,
		Col:      col,
		Controls: controls,
	}

	c.clamp()
	c.setSelected(true)

	return c
}

// Select hands input to the widget: where to go next, and whether this screen is
// done with.
func (c *Cursor) Select(input string) (screen *Screen, exit bool) {
	if c.widget == nil {
		return nil, false
	}
	return c.widget.Select(c, input)
}

func (c *Cursor) Up()    { c.MoveTo(c.Row-1, c.Col) }
func (c *Cursor) Down()  { c.MoveTo(c.Row+1, c.Col) }
func (c *Cursor) Left()  { c.MoveTo(c.Row, c.Col-1) }
func (c *Cursor) Right() { c.MoveTo(c.Row, c.Col+1) }

// MoveTo focuses (row, col). Off-grid positions are ignored, so the cursor stops
// at the edges instead of wrapping.
func (c *Cursor) MoveTo(row, col int) {
	if c.grid == nil || row < 0 || row >= c.grid.height || col < 0 || col >= c.grid.width {
		return
	}

	c.setSelected(false)
	c.Row, c.Col = row, col
	c.setSelected(true)
}

// Cell is the focused cell, or nil if there's no grid or the cursor is off it.
func (c *Cursor) Cell() *Cell {
	if c.grid == nil || c.Row < 0 || c.Row >= len(c.grid.Cells) {
		return nil
	}

	row := c.grid.Cells[c.Row]
	if c.Col < 0 || c.Col >= len(row) {
		return nil
	}

	return row[c.Col]
}

// clamp drags an out-of-range start back onto the grid.
func (c *Cursor) clamp() {
	if c.grid == nil {
		return
	}

	c.Row = min(max(c.Row, 0), max(c.grid.height-1, 0))
	c.Col = min(max(c.Col, 0), max(c.grid.width-1, 0))
}

func (c *Cursor) setSelected(selected bool) {
	if cell := c.Cell(); cell != nil {
		cell.Selected = selected
	}
}
