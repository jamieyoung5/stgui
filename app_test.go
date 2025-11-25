package stgui

import (
	"fmt"
	"testing"
)

// GridWidget is a wrapper to make a Grid satisfy the Widget interface
type GridWidget struct {
	*Grid
}

// Render adapts the Grid.Render to Widget.Render (which accepts a cursor)
func (gw *GridWidget) Render(cursor *Cursor) string {
	return gw.Grid.Render()
}

// Select implements the Widget interface
func (gw *GridWidget) Select(cursor *Cursor, input string) (*Screen, bool) {
	// Simple navigation logic could go here
	// For this playground, we just return nil (no new screen) and false (don't exit)

	fmt.Println(input)
	return nil, false
}

func TestApp_Playground(t *testing.T) {
	// 1. Setup some sample data for the Grid
	data := [][]any{
		{"ID", "Name", "Role"},
		{"1", "Alice", "Engineer"},
		{"2", "Bob", "Designer"},
		{"3", "Charlie", "Manager"},
	}

	// 2. Create the Grid
	grid, err := NewGrid(data, WithGridSymbols())
	if err != nil {
		t.Fatalf("failed to create grid: %v", err)
	}

	// 3. Create a Cursor
	// Map raw keys to abstract signals if needed
	controls := map[string]string{
		"w": "UP",
		"s": "DOWN",
		"a": "LEFT",
		"d": "RIGHT",
	}
	// Cursor at 0,0 inside the widget, and 0,0 on the screen layout
	cursor := NewCursor(0, 0, 0, 0, controls)

	// 4. Wrap Grid in our adapter to make it a Widget
	widget := &GridWidget{Grid: grid}

	// 5. Build the Screen
	// Note: We are accessing the unexported 'elements' field here,
	// which is why this test file must be in package 'stgui'
	screen := &Screen{
		Cursors: []*Cursor{cursor},
		elements: [][]*Element{
			{
				{
					cursor: cursor,
					widget: widget,
				},
			},
		},
	}

	// 6. Initialize the App
	app := NewApp(screen)

	// 7. Display the app
	// We call Display() directly to see the output without blocking on Run()
	t.Log("Rendering App Screen:")
	app.Run()
}
