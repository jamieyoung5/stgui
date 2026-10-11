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
	grid := stgui.NewGrid([][]any{
		{widgets.NewLabel("Hello World")},
		{widgets.NewQuitButton("Quit")},
	}, stgui.WithGridSymbols())

	// One grid, cursor starting on the button.
	screen := widgets.NewGridScreen(grid, 1, 0)

	if err := stgui.NewApp(screen).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```

Arrow keys move the cursor, Enter activates whatever it's on, Ctrl+C quits.

## Beyond one grid

A grid on its own is a screen. For more, wrap it in a `Container`, take a cursor
off it, and stack the rows:

```go
board := widgets.NewContainer(grid).OnKey(onKey)

screen := stgui.NewScreen(board.Focus(0, 0),
	stgui.Row(widgets.NewLabel("Sudoku - 1-9 to fill, q to quit")),
	stgui.Row(board),
	stgui.Row(status),
)

stgui.NewApp(screen).QuitOn("q").Run()
```

`OnKey` gets the input and the cursor, so one handler covers every cell.

## Keys and controls

A key is what was pressed, a control is what it means. Widgets act on controls,
so bindings are swappable:

```go
app.Bind("w", stgui.ControlUp)
```

An unbound key still reaches the widget, with no control attached - that's how
typing works. `Input.Rune` reports false for bound keys, so a key doing a job in
your interface is never typed into a text field as well.

## Actions

Every handler returns an `Action`. Screens are a stack:

```go
stgui.Stay()          // redraw, nothing else
stgui.Push(screen)    // open on top; Close comes back here
stgui.Replace(screen) // swap this screen out
stgui.Close()         // back to the screen underneath, quits if last
stgui.Quit()          // exit now
```

## Events

The loop runs on events, and you can add your own:

```go
type tickEvent struct {
	events.Msg // embed this and it's an event
	at time.Time
}

app.OnEvent(func(event events.Event) stgui.Action {
	if tick, ok := event.(tickEvent); ok {
		clock.Text = tick.at.Format(time.TimeOnly)
	}
	return stgui.Stay()
})
```

`app.Send(tickEvent{at: t})` is safe from any goroutine. The screen is redrawn
after `OnEvent`, so writing to a label animates it.

Widgets are `Label`, `Button`, `Input` and `Container`.

`make run EXAMPLE=<name>` runs one of `counter`, `login`, `sudoku`, `clock`.

## Known limits

- Every character counts as one column, so double-width CJK and combining
  characters won't line up. Box-drawing characters and chess pieces are fine.
- Widgets side by side in the same screen row are aligned by byte length. Keep
  styled or non-ASCII content to one widget per row.
- Frames are clipped to the terminal, so there is no wrapping, no scrolling.
- One cursor per screen, so there's no tabbing between two separate grids.
