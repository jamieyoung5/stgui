package widgets

import (
	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/keyboard"
)

// Container is the standard widget: a Grid you can move around. It takes the
// navigation keys for itself and passes everything else to whatever is in the
// focused cell.
type Container struct {
	*stgui.Grid

	// OnKey picks up what's left: not a navigation command, and not wanted by
	// the focused widget. It gets the cursor too, so a single handler can drive
	// a whole board instead of every cell needing its own widget.
	//
	// Return values mean what they do for Select.
	OnKey func(cursor *stgui.Cursor, input string) (next *stgui.Screen, exit bool)
}

func NewContainer(grid *stgui.Grid) *Container {
	return &Container{Grid: grid}
}

// NewScreen wraps a grid in a container and a cursor at (row, col), for when the
// screen is nothing but the one grid.
//
// Build it yourself for anything more. A board with a title and a status line:
//
//	board := widgets.NewContainer(grid)
//	cursor := board.Focus(0, 0)
//	screen := stgui.NewScreen([]*stgui.Cursor{cursor}, [][]stgui.Widget{
//		{widgets.NewText(title)},
//		{board},
//		{widgets.NewText(status)},
//	})
func NewScreen(grid *stgui.Grid, row, col int) *stgui.Screen {
	container := NewContainer(grid)
	cursor := container.Focus(row, col)

	return stgui.NewScreen(
		[]*stgui.Cursor{cursor},
		[][]stgui.Widget{{container}},
	)
}

// Focus gives you a cursor on this container, starting at (row, col), arrow keys
// wired up.
func (c *Container) Focus(row, col int) *stgui.Cursor {
	return stgui.NewCursor(c, c.Grid, row, col, stgui.DefaultDirectionalControls)
}

// Size is the grid's size, in cells.
func (c *Container) Size() (width, height int) {
	return c.Grid.Size()
}

// Render draws the grid.
func (c *Container) Render() string {
	return c.Grid.Render()
}

// Select moves the cursor on a navigation command. Otherwise it looks at the
// focused cell: Enter activates a Navigator or Clickable child, anything else
// goes to an InputHandler child. Whatever is left over falls through to OnKey.
func (c *Container) Select(cursor *stgui.Cursor, input string) (*stgui.Screen, bool) {
	switch input {
	case stgui.ControlUp:
		cursor.Up()
		return nil, false
	case stgui.ControlDown:
		cursor.Down()
		return nil, false
	case stgui.ControlLeft:
		cursor.Left()
		return nil, false
	case stgui.ControlRight:
		cursor.Right()
		return nil, false
	}

	if cell := cursor.Cell(); cell != nil && cell.Child != nil {
		if input == keyboard.EnterKey {
			switch child := cell.Child.(type) {
			case Navigator:
				return child.OnActivate()
			case Clickable:
				child.OnClick()
				return nil, false
			}
		}

		if handler, ok := cell.Child.(InputHandler); ok {
			handler.HandleInput(input)
			return nil, false
		}
	}

	if c.OnKey != nil {
		return c.OnKey(cursor, input)
	}

	return nil, false
}
