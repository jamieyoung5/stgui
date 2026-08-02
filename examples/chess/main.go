// Chess: an interaction spread over two cells. Enter picks a piece up, Enter
// again puts it down, and what the key does depends on whether anything is held.
//
// A board, not a game - nothing stops you playing Nh1xa8.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/keyboard"
	"github.com/jamieyoung5/stgui/widgets"
)

const size = 8

// Colours. The squares have no divider between them and are three columns wide,
// so the backgrounds touch and read as a checkerboard.
var (
	lightSquare  = stgui.Bg256(180)
	darkSquare   = stgui.Bg256(94)
	heldSquare   = stgui.Bg256(178)
	cursorSquare = stgui.Bg256(107) + stgui.Fg256(232)

	whitePiece = stgui.Fg256(255) + stgui.StyleBold
	blackPiece = stgui.Fg256(232) + stgui.StyleBold
)

// Starting position, black's back rank first.
var startingPosition = [size]string{
	"♜♞♝♛♚♝♞♜",
	"♟♟♟♟♟♟♟♟",
	"        ",
	"        ",
	"        ",
	"        ",
	"♙♙♙♙♙♙♙♙",
	"♖♘♗♕♔♗♘♖",
}

type game struct {
	grid   *stgui.Grid
	pieces [size][size]string
	status *widgets.Label

	holding          string
	fromRow, fromCol int
}

func newGame() *game {
	style := &stgui.GridStyle{
		NoValue:       " ",
		MinCellWidth:  3,
		Align:         stgui.AlignCenter,
		SelectedStyle: cursorSquare,
	}

	g := &game{
		grid:   stgui.NewEmptyGrid(size, size, style),
		status: widgets.NewLabel(""),
	}

	for row, rank := range startingPosition {
		for col, piece := range []rune(rank) {
			if piece != ' ' {
				g.pieces[row][col] = string(piece)
			}
		}
	}

	g.status.Text = "Enter to pick up a piece"
	g.draw()

	return g
}

// onKey runs for anything the container didn't use to move the cursor. Only Enter
// means anything, and it means two different things.
func (g *game) onKey(cursor *stgui.Cursor, input string) (*stgui.Screen, bool) {
	if input != keyboard.EnterKey {
		return nil, false
	}

	row, col := cursor.Row, cursor.Col

	switch {
	case g.holding == "":
		if g.pieces[row][col] == "" {
			g.status.Text = fmt.Sprintf("%s is empty", square(row, col))
			break
		}

		g.holding = g.pieces[row][col]
		g.fromRow, g.fromCol = row, col
		g.status.Text = fmt.Sprintf("Holding %s from %s", g.holding, square(row, col))

	case row == g.fromRow && col == g.fromCol:
		g.status.Text = fmt.Sprintf("Put %s back on %s", g.holding, square(row, col))
		g.holding = ""

	default:
		g.pieces[g.fromRow][g.fromCol] = ""
		g.pieces[row][col] = g.holding
		g.status.Text = fmt.Sprintf("%s to %s", g.holding, square(row, col))
		g.holding = ""
	}

	g.draw()
	return nil, false
}

// draw lays the pieces out and colours the squares, marking the one a held piece
// came from. A cell's style covers the padding too, so each square comes out a
// solid block.
func (g *game) draw() {
	g.grid.Each(func(row, col int, cell *stgui.Cell) {
		piece := g.pieces[row][col]
		g.grid.Set(row, col, piece)

		background := lightSquare
		if (row+col)%2 == 1 {
			background = darkSquare
		}
		if g.holding != "" && row == g.fromRow && col == g.fromCol {
			background = heldSquare
		}

		cell.Style = background + pieceColour(piece)
	})
}

// pieceColour: outlined glyphs are white's, solid ones are black's.
func pieceColour(piece string) string {
	if strings.ContainsAny(piece, "♔♕♖♗♘♙") {
		return whitePiece
	}
	return blackPiece
}

// square as a chess player would write it, "e4".
func square(row, col int) string {
	return fmt.Sprintf("%c%d", 'a'+col, size-row)
}

func main() {
	g := newGame()

	board := widgets.NewContainer(g.grid)
	board.OnKey = g.onKey
	cursor := board.Focus(size-2, 4) // e2

	screen := stgui.NewScreen(
		[]*stgui.Cursor{cursor},
		[][]stgui.Widget{
			{widgets.NewText(widgets.NewLabel("Chess - Enter to pick up and put down, q to quit"))},
			{board},
			{widgets.NewText(g.status)},
		},
	)

	app := stgui.NewApp(screen)
	app.QuitKeys = append(app.QuitKeys, "q")

	if err := app.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
