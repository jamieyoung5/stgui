package stgui_test

import (
	"testing"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/widgets"
)

// sudokuStyle is the examples/sudoku look: box dividers, centred, fixed width.
func sudokuStyle() *stgui.GridStyle {
	return &stgui.GridStyle{
		VerticalDivider:   "│",
		HorizontalDivider: "─",
		Intersection:      "┼",
		NoValue:           " ",
		MinCellWidth:      3,
		Align:             stgui.AlignCenter,
		DrawColDivider:    func(col int) bool { return col%3 == 2 },
		DrawRowDivider:    func(row int) bool { return row%3 == 2 },
	}
}

func sudokuGrid() *stgui.Grid {
	g := stgui.NewEmptyGrid(9, 9, sudokuStyle())
	g.Each(func(row, col int, cell *stgui.Cell) {
		if (row+col)%3 != 0 {
			cell.Value = (row*3+col)%9 + 1
		}
	})

	return g
}

// boardGrid is a chessboard: every cell styled, no dividers, wide cells.
func boardGrid() *stgui.Grid {
	style := &stgui.GridStyle{
		NoValue:       " ",
		MinCellWidth:  3,
		Align:         stgui.AlignCenter,
		SelectedStyle: stgui.StyleBold + stgui.StyleBGGreen,
	}

	g := stgui.NewEmptyGrid(8, 8, style)
	pieces := []string{"♜", "♞", "♝", "♛", "♚", "♝", "♞", "♜"}
	g.Each(func(row, col int, cell *stgui.Cell) {
		if (row+col)%2 == 0 {
			cell.Style = stgui.Bg256(180) + stgui.Fg256(16)
		} else {
			cell.Style = stgui.Bg256(94) + stgui.Fg256(231)
		}
		switch row {
		case 0, 7:
			cell.Value = pieces[col]
		case 1, 6:
			cell.Value = "♟"
		}
	})

	return g
}

func BenchmarkGridRenderSudoku(b *testing.B) {
	g := sudokuGrid()
	b.ReportAllocs()
	for b.Loop() {
		_ = g.Render()
	}
}

func BenchmarkGridRenderBoard(b *testing.B) {
	g := boardGrid()
	b.ReportAllocs()
	for b.Loop() {
		_ = g.Render()
	}
}

func BenchmarkGridRenderLarge(b *testing.B) {
	g := stgui.NewEmptyGrid(40, 40, sudokuStyle())
	g.Each(func(row, col int, cell *stgui.Cell) { cell.Value = (row * col) % 10 })
	b.ReportAllocs()
	for b.Loop() {
		_ = g.RenderLines()
	}
}

// screen puts a board next to a sudoku with labels above and below, the shape a
// real app ends up with.
func benchScreen() *stgui.Screen {
	board := widgets.NewContainer(boardGrid())
	side := widgets.NewContainer(sudokuGrid())

	return stgui.NewScreen(
		board.Focus(0, 0),
		stgui.Row(widgets.NewLabel("stgui benchmark")),
		stgui.Row(board, side),
		stgui.Row(widgets.NewLabel("arrows to move, enter to select")),
	)
}

func BenchmarkScreenRender(b *testing.B) {
	s := benchScreen()
	b.ReportAllocs()
	for b.Loop() {
		_ = s.Render()
	}
}

func BenchmarkCellRenderLines(b *testing.B) {
	c := &stgui.Cell{Value: 42, Style: stgui.Bg256(180)}
	b.ReportAllocs()
	for b.Loop() {
		_ = c.RenderLines()
	}
}

func BenchmarkVisibleWidth(b *testing.B) {
	line := stgui.Bg256(180) + stgui.Fg256(16) + " ♞ " + "\033[0m" +
		stgui.Bg256(94) + " ♟ " + "\033[0m"
	b.ReportAllocs()
	for b.Loop() {
		_ = stgui.VisibleWidth(line)
	}
}
