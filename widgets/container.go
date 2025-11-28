package widgets

import (
	"github.com/jamieyoung5/stgui"
)

// Container is a standard widget implementation that wraps a Grid.
// It handles basic navigation (UP, DOWN, LEFT, RIGHT) and delegates
// specific interactions (Enter, text input) to the active child widget.
type Container struct {
	*stgui.Grid
}

func NewContainer(grid *stgui.Grid) *Container {
	return &Container{Grid: grid}
}

// Size returns the dimensions of the underlying grid.
func (c *Container) Size() (width, height int) {
	return c.Grid.Size()
}

// Render returns the string representation of the grid.
func (c *Container) Render() string {
	return c.Grid.Render()
}

// Select handles input events.
// It processes navigation commands to move the cursor.
// For other inputs, it checks the active cell:
// - "ENTER" triggers OnClick if the child is Clickable.
// - Any other input is passed to HandleInput if the child is an InputHandler.
func (c *Container) Select(cursor *stgui.Cursor, input string) (*stgui.Screen, bool) {
	switch input {
	case "UP":
		cursor.Up()
		return nil, false
	case "DOWN":
		cursor.Down()
		return nil, false
	case "LEFT":
		cursor.Left()
		return nil, false
	case "RIGHT":
		cursor.Right()
		return nil, false
	}

	if cursor.Row >= 0 && cursor.Row < len(c.Grid.Cells) {
		row := c.Grid.Cells[cursor.Row]
		if cursor.Col >= 0 && cursor.Col < len(row) {
			cell := row[cursor.Col]

			if input == "ENTER" {
				if clickable, ok := cell.Child.(Clickable); ok {
					clickable.OnClick()
					return nil, false
				}
			}

			if handler, ok := cell.Child.(InputHandler); ok {
				handler.HandleInput(input)
				return nil, false
			}
		}
	}

	return nil, false
}
