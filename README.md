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

A grid on its own is a screen. For anything more, wrap the grid in a `Container`,
take a cursor off it, and stack the rows up:

```go
board := widgets.NewContainer(grid).OnKey(onKey)

screen := stgui.NewScreen(board.Focus(0, 0),
	stgui.Row(widgets.NewLabel("Sudoku - 1-9 to fill, q to quit")),
	stgui.Row(board),
	stgui.Row(status), // a *Label you write to as the app runs
)

stgui.NewApp(screen).QuitOn("q").Run()
```

`Container.OnKey` gets the input and the cursor, so one handler covers every
cell instead of each cell needing its own widget.

`OnKey`, `OnEvent`, `Bind` and `QuitOn` all chain off their constructor.

## Keys and controls

A **key** is what was pressed. A **control** is what it means. Widgets act on
controls, so bindings are swappable without any widget knowing:

```go
app.Bind("w", stgui.ControlUp)     // now w moves up, everywhere
app.Bind(" ", stgui.ControlSelect) // space activates, as well as Enter
```

Bindings live on the `App`, so one call covers every screen. Keys arrive at
widgets either way. An unbound key just has no control attached, which is how
typing works:

```go
func onKey(cursor *stgui.Cursor, in stgui.Input) stgui.Action {
	if in.Is(stgui.ControlSelect) { ... }        // a control
	if digit, typed := in.Rune(); typed { ... }  // a character
	if in.Key == stgui.KeyBackspace { ... }      // a named key
}
```

`Rune` reports false for named keys *and* for bound keys, so a key doing a job in
your interface never gets typed into a text field as well.

Controls are lowercase (`"up"`) and keys are uppercase (`"UP"`), so the two can't
be mistaken for each other.

## Actions

Every handler returns an `Action` saying what the app should do next. Screens are
a stack:

```go
stgui.Stay()          // handled it, redraw, nothing else
stgui.Push(screen)    // open on top; Close comes back here
stgui.Replace(screen) // swap this screen out, no going back
stgui.Close()         // back to the screen underneath, quits if last
stgui.Quit()          // tear it all down now
```

The zero `Action` is `Stay`, so a handler with nothing to say can return
`stgui.Action{}`.

## Events

The loop is driven by events, and you can add your own: a ticker, a network
reply, anything from another goroutine.

```go
type tickEvent struct {
	events.Msg // embed this and it's an event
	at time.Time
}

app := stgui.NewApp(screen).OnEvent(func(event events.Event) stgui.Action {
	if tick, ok := event.(tickEvent); ok {
		clock.Text = tick.at.Format(time.TimeOnly)
	}
	return stgui.Stay()
})

go func() {
	for t := range time.Tick(time.Second) {
		app.Send(tickEvent{at: t}) // safe from any goroutine
	}
}()
```

The screen is redrawn after `OnEvent`, so writing to a label is enough to animate
it. `Send` gives up once the app has quit, so a sending goroutine never leaks.

## The pieces

| | |
|---|---|
| `App` | Owns the terminal, runs the event loop over a stack of screens. |
| `Screen` | One terminal-full: rows of widgets, plus the cursor taking input. |
| `Grid` | Rows and columns of cells. `Set`/`Get`/`Each` to drive it. |
| `Cell` | Holds a value or a widget, plus a `Style`. |
| `Cursor` | The focused cell. Feeds input to its widget. |
| `Input` | One key press: which `Key`, and what `Control` it's bound to. |
| `Action` | What a handler wants the app to do next. |
| `Widget` | Anything with `Render() string` and `Select(...) Action`. |
| `Renderable` | Anything with `RenderLines() []string`. Goes in a cell. |

Built in: `Label`, `Button`, `Input`, `Container`. A `Label` is both a
`Renderable` and a `Widget`, so it works in a cell and as a screen row.

`NewButton` runs a callback, `NewNavButton` opens a screen, `NewQuitButton`
closes one, and `NewActionButton` decides which to do when pressed:

```go
widgets.NewActionButton("Login", func() stgui.Action {
	if !valid() {
		status.Text = "Denied"
		return stgui.Stay()
	}
	return stgui.Push(welcome)
})
```

`make run EXAMPLE=<name>` runs one of `counter`, `login`, `sudoku`, `chess`,
`clock`.

## Known limits

- Every character counts as one column, so double-width CJK and combining
  characters won't line up. Box-drawing characters and chess pieces are fine.
- Widgets side by side in the same screen row are aligned by byte length. Keep
  styled or non-ASCII content to one widget per row.
- Frames are clipped to the terminal, so there is no wrapping, no scrolling.
- One cursor per screen, so there's no tabbing between two separate grids.
