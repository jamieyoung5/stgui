package stgui

import "fmt"

type Cursor struct {
	widget      Widget
	currentCell *Cell
	parentCell  *Cell

	controls       map[string]string
	selectedColour string
}

func NewCursor(widget Widget, startingCell *Cell, controls map[string]string) *Cursor {
	return &Cursor{
		widget:      widget,
		currentCell: startingCell,
		controls:    controls,
	}
}

func (c *Cursor) Select(input string) (screen *Screen, exit bool) {
	return c.widget.Select(c, input)
}

func (c *Cursor) Up() {
	if c.currentCell.Up != nil {
		c.currentCell.Selected = false
		c.currentCell = c.currentCell.Up
		c.currentCell.Selected = true
	}
}

func (c *Cursor) Down() {
	if c.currentCell.Down != nil {
		c.currentCell.Selected = false
		c.currentCell = c.currentCell.Down
		c.currentCell.Selected = true
	}
}

func (c *Cursor) Right() {
	fmt.Printf("RIGHT CALLED")
	if c.currentCell.Right != nil {
		c.currentCell.Selected = false
		c.currentCell = c.currentCell.Right
		c.currentCell.Selected = true
	}
}

func (c *Cursor) Left() {
	if c.currentCell.Left != nil {
		c.currentCell.Selected = false
		c.currentCell = c.currentCell.Left
		c.currentCell.Selected = true
	}
}
