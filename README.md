# stgui

A simple, grid-based terminal GUI library for Go.

## Installation

```bash
go get [github.com/jamieyoung5/stgui](https://github.com/jamieyoung5/stgui)
```

## Quick Start

```go
package main

import (
	"fmt"
	"os"

	"[github.com/jamieyoung5/stgui](https://github.com/jamieyoung5/stgui)"
	"[github.com/jamieyoung5/stgui/widgets](https://github.com/jamieyoung5/stgui/widgets)"
)

func main() {
	// Create Widgets
	lbl := widgets.NewLabel("Hello World")
	btn := widgets.NewButton("Quit", func() {
		os.Exit(0)
	})

	// Define Grid Layout
	gridData := [][]any{
		{lbl},
		{btn},
	}
	grid, _ := stgui.NewGrid(gridData, stgui.WithGridSymbols())

	// Setup Container and Cursor
	container := widgets.NewContainer(grid)
	cursor := stgui.NewCursor(container, grid, 1, 0, stgui.DefaultDirectionalControls)
	grid.Cells[1][0].Selected = true // Select the button initially

	// Initialize Screen and App
	screen := stgui.NewScreen(
		[]*stgui.Cursor{cursor},
		[][]stgui.Widget{{container}},
	)

	if err := stgui.NewApp(screen).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```
