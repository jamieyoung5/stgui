package stgui

import "testing"

// A 3x3 grid of 3-wide cells, so column c starts at 3c.
func hitGrid(t *testing.T) *gridWidget {
	t.Helper()

	g := NewEmptyGrid(3, 3, &GridStyle{MinCellWidth: 3, NoValue: "   "})
	for r := range 3 {
		for c := range 3 {
			g.Set(r, c, " x ")
		}
	}
	g.RenderLines() // CellAt reads the last frame

	return &gridWidget{Grid: g}
}

func TestGridCellAt(t *testing.T) {
	g := hitGrid(t)

	cases := []struct {
		col, row     int
		wantR, wantC int
		wantOK       bool
	}{
		{0, 0, 0, 0, true},
		{2, 0, 0, 0, true},
		{3, 0, 0, 1, true},
		{8, 2, 2, 2, true},
		{9, 0, 0, 0, false},  // past the right edge
		{0, 3, 0, 0, false},  // past the bottom
		{-1, 0, 0, 0, false}, // left of the grid
		{0, -1, 0, 0, false},
	}

	for _, c := range cases {
		r, col, ok := g.CellAt(c.col, c.row)
		if ok != c.wantOK || (ok && (r != c.wantR || col != c.wantC)) {
			t.Errorf("CellAt(%d,%d) = %d,%d,%v want %d,%d,%v",
				c.col, c.row, r, col, ok, c.wantR, c.wantC, c.wantOK)
		}
	}
}

func TestGridCellAtBeforeRender(t *testing.T) {
	g := NewEmptyGrid(2, 2, &GridStyle{MinCellWidth: 3})
	if _, _, ok := g.CellAt(0, 0); ok {
		t.Error("CellAt should report false before the grid has been drawn")
	}
}

func TestScreenHitSkipsRowGaps(t *testing.T) {
	g := hitGrid(t)
	top := &testWidget{lines: []string{"title"}}

	screen := NewScreen(nil, Row(top), Row(g))
	screen.renderLines()

	// Line 0 is the title, line 1 the blank between rows, lines 2-4 the grid.
	if _, _, _, ok := screen.Hit(0, 1); ok {
		t.Error("the blank line between rows should hit nothing")
	}

	widget, row, col, ok := screen.Hit(4, 3)
	if !ok {
		t.Fatal("expected a hit inside the grid")
	}
	if widget != Widget(g) {
		t.Error("hit the wrong widget")
	}
	if row != 1 || col != 1 {
		t.Errorf("got cell %d,%d want 1,1", row, col)
	}
}

func TestScreenHitSideBySide(t *testing.T) {
	left := hitGrid(t)
	right := hitGrid(t)

	screen := NewScreen(nil, Row(left, right))
	screen.renderLines()

	// The left grid is 9 wide, then rowSpacing before the right one starts.
	if _, _, _, ok := screen.Hit(9, 0); ok {
		t.Error("the gap between widgets should hit nothing")
	}

	widget, _, col, ok := screen.Hit(9+rowSpacing+3, 0)
	if !ok {
		t.Fatal("expected a hit in the right grid")
	}
	if widget != Widget(right) {
		t.Error("hit the left grid, expected the right one")
	}
	if col != 1 {
		t.Errorf("got column %d want 1", col)
	}
}

type testWidget struct{ lines []string }

func (w *testWidget) RenderLines() []string { return w.lines }
func (w *testWidget) Render() string        { return w.lines[0] }
func (w *testWidget) Select(*Cursor, Input) Action {
	return Stay()
}
