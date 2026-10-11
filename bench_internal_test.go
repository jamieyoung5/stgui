package stgui

import (
	"bufio"
	"io"
	"testing"
)

// gridWidget is the grid as a screen widget, standing in for widgets.Container
// (which an internal test can't import).
type gridWidget struct {
	*Grid
}

func (g *gridWidget) Select(cursor *Cursor, input Input) Action { return Stay() }

// boardScreen is a styled 8x8 board with a cursor on it, the shape of a real app.
func boardScreen() (*Screen, *Cursor) {
	style := &GridStyle{
		NoValue:       " ",
		MinCellWidth:  3,
		Align:         AlignCenter,
		SelectedStyle: StyleBold + StyleBGGreen,
	}

	grid := NewEmptyGrid(8, 8, style)
	grid.Each(func(row, col int, cell *Cell) {
		if (row+col)%2 == 0 {
			cell.Style = Bg256(180) + Fg256(16)
		} else {
			cell.Style = Bg256(94) + Fg256(231)
		}
		if row < 2 || row > 5 {
			cell.Value = "♟"
		}
	})

	widget := &gridWidget{Grid: grid}
	cursor := NewCursor(widget, grid, 0, 0)

	return NewScreen(cursor, Row(widget)), cursor
}

// benchApp sends frames to a discarding writer instead of the terminal.
func benchApp(screen *Screen) *App {
	app := NewApp(screen)
	app.out = bufio.NewWriter(io.Discard)
	app.width, app.height = 120, 40

	return app
}

// BenchmarkFrame is a full repaint: render the screen, clip it, diff it against
// the last one and write what changed.
func BenchmarkFrame(b *testing.B) {
	screen, _ := boardScreen()
	app := benchApp(screen)

	b.ReportAllocs()
	for b.Loop() {
		app.lastBuffer = nil
		app.Display()
	}
}

// BenchmarkFrameDiffed is the common case: one cell changed, so the diff writes
// two lines and skips the rest.
func BenchmarkFrameDiffed(b *testing.B) {
	screen, cursor := boardScreen()
	app := benchApp(screen)
	app.Display()

	b.ReportAllocs()
	for b.Loop() {
		cursor.Down()
		cursor.Up()
		app.Display()
	}
}
