// Sudoku: one key handler for the whole board. The grid is the display, a plain
// array is the game, and Container.OnKey joins the two.
package main

import (
	"fmt"
	"os"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/keyboard"
	"github.com/jamieyoung5/stgui/widgets"
)

const size = 9

// Left to right, top to bottom. 0 is an empty square.
const puzzle = "" +
	"530070000" +
	"600195000" +
	"098000060" +
	"800060003" +
	"400803001" +
	"700020006" +
	"060000280" +
	"000419005" +
	"000080079"

type board struct {
	grid   *stgui.Grid
	values [size][size]int
	given  [size][size]bool
	status *widgets.Label
}

func newBoard(puzzle string) *board {
	// Dividers on the box boundaries only. Without this it reads as a
	// spreadsheet.
	style := &stgui.GridStyle{
		VerticalDivider:   "│",
		HorizontalDivider: "─",
		Intersection:      "┼",
		NoValue:           " ",
		MinCellWidth:      3,
		Align:             stgui.AlignCenter,
		DrawColDivider:    func(col int) bool { return col%3 == 2 },
		DrawRowDivider:    func(row int) bool { return row%3 == 2 },
	}

	b := &board{
		grid:   stgui.NewEmptyGrid(size, size, style),
		status: widgets.NewLabel(""),
	}

	for i, digit := range puzzle {
		row, col := i/size, i%size
		if digit == '0' {
			continue
		}

		b.values[row][col] = int(digit - '0')
		b.given[row][col] = true
	}

	b.draw()
	return b
}

// onKey gets every key the container didn't use to move, plus the cursor, so it
// always knows which square we're on.
func (b *board) onKey(cursor *stgui.Cursor, input string) (*stgui.Screen, bool) {
	row, col := cursor.Row, cursor.Col
	if b.given[row][col] {
		return nil, false
	}

	switch {
	case len(input) == 1 && input[0] >= '1' && input[0] <= '9':
		b.values[row][col] = int(input[0] - '0')
	case input == "0" || input == " " || input == keyboard.BackspaceKey:
		b.values[row][col] = 0
	default:
		return nil, false
	}

	b.draw()
	return nil, false
}

// draw pushes the game into the grid: what's in each square, and a colour for
// whether it was given, typed in, or clashing.
func (b *board) draw() {
	clashes := b.clashes()
	filled := 0

	b.grid.Each(func(row, col int, cell *stgui.Cell) {
		value := b.values[row][col]
		if value == 0 {
			b.grid.Set(row, col, nil)
			cell.Style = ""
			return
		}

		filled++
		b.grid.Set(row, col, value)

		switch {
		case clashes[row][col]:
			cell.Style = stgui.StyleRed + stgui.StyleBold
		case b.given[row][col]:
			cell.Style = stgui.StyleBold
		default:
			cell.Style = stgui.StyleCyan
		}
	})

	b.status.Text = b.summary(filled, clashes)
}

func (b *board) summary(filled int, clashes [size][size]bool) string {
	clashing := 0
	for row := range clashes {
		for _, clash := range clashes[row] {
			if clash {
				clashing++
			}
		}
	}

	switch {
	case clashing > 0:
		return fmt.Sprintf("%d squares clash", clashing)
	case filled == size*size:
		return "Solved."
	default:
		return fmt.Sprintf("%d of %d filled", filled, size*size)
	}
}

// clashes finds digits repeated in a row, a column or a box.
func (b *board) clashes() [size][size]bool {
	var clashing [size][size]bool

	for row := range size {
		for col := range size {
			value := b.values[row][col]
			if value == 0 {
				continue
			}

			for other := range size {
				if other != col && b.values[row][other] == value {
					clashing[row][col] = true
				}
				if other != row && b.values[other][col] == value {
					clashing[row][col] = true
				}
			}

			boxRow, boxCol := (row/3)*3, (col/3)*3
			for r := boxRow; r < boxRow+3; r++ {
				for c := boxCol; c < boxCol+3; c++ {
					if (r != row || c != col) && b.values[r][c] == value {
						clashing[row][col] = true
					}
				}
			}
		}
	}

	return clashing
}

func main() {
	b := newBoard(puzzle)

	container := widgets.NewContainer(b.grid)
	container.OnKey = b.onKey
	cursor := container.Focus(0, 0)

	screen := stgui.NewScreen(
		[]*stgui.Cursor{cursor},
		[][]stgui.Widget{
			{widgets.NewText(widgets.NewLabel("Sudoku - 1-9 to fill, 0 to clear, q to quit"))},
			{container},
			{widgets.NewText(b.status)},
		},
	)

	app := stgui.NewApp(screen)
	app.QuitKeys = append(app.QuitKeys, "q")

	if err := app.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
