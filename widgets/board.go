package widgets

import (
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
	selector          func(cursor *stgui.Cursor, input string) (*stgui.Screen, bool)
	traversalControls BoardControls
}

func (b *Board) Size() (height int, width int) {
	return b.grid.Size()
}

func (b *Board) Render() string {
	return b.grid.Render()
}

func (b *Board) Select(cursor *stgui.Cursor, input string) (screen *stgui.Screen, exit bool) {
	switch input {
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

	return b.selector(cursor, input)
}
