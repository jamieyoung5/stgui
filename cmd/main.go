package main

import (
	"log"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/keyboard"
)

// GridWidget wraps stgui.Grid to satisfy the stgui.Widget interface.
// We embed *stgui.Grid so that methods like Render() and Size() are automatically promoted
// to satisfy the interface, leaving us to implement only Select().
type GridWidget struct {
	*stgui.Grid
}

// Select handles input. Returns exit=true to close the screen.
func (gw *GridWidget) Select(c *stgui.Cursor, input string) (*stgui.Screen, bool) {
	// Simple navigation or logic mapping
	switch input {
	case "UP":
		c.Up()
	case "DOWN":
		c.Down()
	case "LEFT":
		c.Left()
	case "RIGHT":
		c.Right()
	case "QUIT":
		return nil, true // Exit the screen
	}
	return nil, false // Stay on this screen
}

func main() {
	// 1. Create Data
	data := [][]any{
		{"ID", "Name", "Role", "Status"},
		{"01", "Alice", "Dev", "Active"},
		{"02", "Bob", "Design", "Offline"},
		{"03", "Charlie", "Product", "Active"},
	}

	// 2. Create Grid
	grid, err := stgui.NewGrid(data, stgui.WithGridSymbols())
	if err != nil {
		log.Fatalf("Failed to create grid: %v", err)
	}

	// 3. Create Controls Mapping
	controls := map[string]string{
		keyboard.UpArrowKey:    "UP",
		keyboard.DownArrowKey:  "DOWN",
		keyboard.LeftArrowKey:  "LEFT",
		keyboard.RightArrowKey: "RIGHT",
		"q":                    "QUIT",
		//keyboard.EscapeKey:     "QUIT",
	}

	// 5. Wrap Grid in Widget
	// The GridWidget struct defined above allows the Grid to satisfy the Widget interface.
	gridWidget := &GridWidget{Grid: grid}

	// 4. Create Cursor
	// We attach the cursor to the root cell of the grid so it starts at top-left.
	grid.Root.Selected = true
	cursor := stgui.NewCursor(gridWidget, grid.Root, controls)

	// 6. Create Screen
	// We pass a slice of cursors (just one here) and a 2D layout of widgets.
	screen := stgui.NewScreen(
		[]*stgui.Cursor{cursor},
		[][]stgui.Widget{{gridWidget}},
	)

	// 7. Create and Run App
	app := stgui.NewApp(screen)
	app.Run()
}
