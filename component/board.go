package component

import (
	"fmt"

	"github.com/jamieyoung5/stgui"
)

type BoardControls struct {
	Up    []byte
	Down  []byte
	Left  []byte
	Right []byte
}

type Board struct {
	grid              *stgui.Grid
	styling           *stgui.Symbols
	selector          func(cursor *stgui.Cursor, macro string) (*stgui.Screen, bool)
	traversalControls BoardControls
}

func (b *Board) GetDimensions() (height int, width int) {
	return b.grid.Size()
}

func (b *Board) Print(cursor *stgui.Cursor) {
	fmt.Printf(b.Render(cursor))
}

func (b *Board) Render(cursor *stgui.Cursor) string {
	return b.grid.Render()
}

func (b *Board) Select(cursor *stgui.Cursor, macro string) (screen *stgui.Screen, exit bool) {
	switch macro {
	case string(b.traversalControls.Up):
		cursor.Up()
		return
	case string(b.traversalControls.Down):
		cursor.Down()
		return
	case string(b.traversalControls.Left):
		cursor.Left()
		return
	case string(b.traversalControls.Right):
		cursor.Right()
		return
	}

	return b.selector(cursor, macro)
}
