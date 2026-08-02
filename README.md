# stgui

A small, lightweight grid-based terminal GUI library for Go.

## Installation

```bash
go get github.com/jamieyoung5/stgui
```

## Quick start

```go
package main

import (
	"fmt"
	"os"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/widgets"
)

func main() {
	lbl := widgets.NewLabel("Hello World")
	btn := widgets.NewQuitButton("Quit")

	grid, err := stgui.NewGrid([][]any{
		{lbl},
		{btn},
	}, stgui.WithGridSymbols())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// One grid, cursor starting on the button.
	screen := widgets.NewScreen(grid, 1, 0)

	if err := stgui.NewApp(screen).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```

`make run EXAMPLE=<name>` runs one of `counter`, `login`, `sudoku`, `chess`.

## Known limits

- Every character counts as one column, so double-width CJK and combining
  characters won't line up. Box-drawing characters and chess pieces are fine.
- Widgets side by side in the same screen row are aligned by byte length. Keep
  styled or non-ASCII content to one widget per row.
- One cursor per screen: `ActiveCursor` takes the first non-nil one.
- Frames are clipped to the terminal, so there is no wrapping, no scrolling.
