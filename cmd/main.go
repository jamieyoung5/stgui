package main

import (
	"fmt"
	"log"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/keyboard"
)

// GridWidget wraps stgui.Grid to satisfy the stgui.Widget interface.
// Since Grid.Render() doesn't take a cursor and Grid lacks Select(),
// we adapt it here.
type GridWidget struct {
	*stgui.Grid
}

// Render adapts the interface.
func (gw *GridWidget) Render(c *stgui.Cursor) string {
	return gw.Grid.Render()
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
		keyboard.EscapeKey:     "QUIT",
	}

	// 4. Create Cursor
	cursor := stgui.NewCursor(0, 0, 0, 0, controls)

	// 5. Wrap Grid in Widget Adapter
	widget := &GridWidget{Grid: grid}

	// 6. Create Element
	element := stgui.NewElement(cursor, widget)

	// 7. Create Screen
	// Layout is a 1x1 grid of elements
	screenLayout := [][]*stgui.Element{
		{element},
	}
	screen := stgui.NewScreen([]*stgui.Cursor{cursor}, screenLayout)

	// 8. Run App
	fmt.Println("Starting stgui demo... Press 'q' or 'ESC' to quit.")
	app := stgui.NewApp(screen)
	app.Run()
}
