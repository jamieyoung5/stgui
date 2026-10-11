package widgets

import (
	"github.com/jamieyoung5/stgui"
)

// KeyHandler takes an input the container had no use for. It gets the cursor
// too, so one handler can drive a whole board.
type KeyHandler func(cursor *stgui.Cursor, input stgui.Input) stgui.Action

// Container is the standard widget: a Grid you can move around. It takes the
// navigation controls for itself and passes everything else to whatever is in
// the focused cell.
type Container struct {
	*stgui.Grid

	onKey KeyHandler
}

func NewContainer(grid *stgui.Grid) *Container {
	return &Container{Grid: grid}
}

// OnKey sets the handler for what's left over: not a control the cursor used,
// and not wanted by the focused widget.
func (c *Container) OnKey(handler KeyHandler) *Container {
	c.onKey = handler
	return c
}

// NewGridScreen wraps a grid in a container and a cursor at (row, col), for when
// the screen is nothing but the one grid.
//
// Build it yourself for anything more. A board with a title and a status line:
//
//	board := widgets.NewContainer(grid)
//	cursor := board.Focus(0, 0)
//	screen := stgui.NewScreen(cursor,
//		stgui.Row(widgets.NewLabel(title)),
//		stgui.Row(board),
//		stgui.Row(status),
//	)
func NewGridScreen(grid *stgui.Grid, row, col int) *stgui.Screen {
	container := NewContainer(grid)

	return stgui.NewScreen(
		container.Focus(row, col),
		stgui.Row(container),
	)
}

// Focus gives you a cursor on this container, starting at (row, col).
func (c *Container) Focus(row, col int) *stgui.Cursor {
	return stgui.NewCursor(c, c.Grid, row, col)
}

// Render draws the grid.
func (c *Container) Render() string {
	return c.Grid.Render()
}

// Select moves the cursor on a directional control. Otherwise it looks at the
// focused cell: ControlSelect activates a Navigator or Clickable child, anything
// else goes to an InputHandler child. Whatever is left falls through to the
// OnKey handler.
func (c *Container) Select(cursor *stgui.Cursor, input stgui.Input) stgui.Action {
	switch input.Control {
	case stgui.ControlUp:
		cursor.Up()
		return stgui.Stay()
	case stgui.ControlDown:
		cursor.Down()
		return stgui.Stay()
	case stgui.ControlLeft:
		cursor.Left()
		return stgui.Stay()
	case stgui.ControlRight:
		cursor.Right()
		return stgui.Stay()
	}

	if cell := cursor.Cell(); cell != nil && cell.Child != nil {
		if input.Is(stgui.ControlSelect) {
			switch child := cell.Child.(type) {
			case Navigator:
				return child.OnActivate()
			case Clickable:
				child.OnClick()
				return stgui.Stay()
			}
		}

		if handler, ok := cell.Child.(InputHandler); ok {
			handler.HandleInput(input)
			return stgui.Stay()
		}
	}

	if c.onKey != nil {
		return c.onKey(cursor, input)
	}

	return stgui.Stay()
}
