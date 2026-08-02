package stgui_test

import (
	"strings"
	"testing"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/testutil"
)

func TestGrid_NewGrid(t *testing.T) {
	data := [][]any{
		{"A", "B", "C"},
		{"D", "E", "F"},
	}
	g, err := stgui.NewGrid(data, nil)
	if err != nil {
		t.Fatalf("NewGrid failed: %v", err)
	}
	w, h := g.Size()
	if w != 3 || h != 2 {
		t.Errorf("Expected size 3x2 (cols x rows), got %dx%d", w, h)
	}
	if g.Cells[0][0].Value != "A" {
		t.Errorf("Expected cell 0,0 to be 'A', got %v", g.Cells[0][0].Value)
	}
}

func TestGrid_NewGrid_Empty(t *testing.T) {
	g, err := stgui.NewGrid([][]any{}, nil)
	if err != nil {
		t.Fatalf("NewGrid failed on empty: %v", err)
	}
	if g.Cells != nil {
		t.Error("Expected nil Cells for empty grid")
	}
}

// Used to panic: NewGrid left Style nil on an empty grid.
func TestGrid_RenderLines_Empty(t *testing.T) {
	g, _ := stgui.NewGrid([][]any{}, stgui.WithGridSymbols())

	lines := g.RenderLines()
	if len(lines) != 1 || lines[0] != "." {
		t.Errorf("Expected ['.'] for an empty grid, got %q", lines)
	}

	bare := &stgui.Grid{}
	if lines := bare.RenderLines(); len(lines) != 1 {
		t.Errorf("Expected a single line for a zero grid, got %q", lines)
	}
}

func TestGrid_NewGrid_Ragged(t *testing.T) {
	data := [][]any{
		{"A"},
		{"B", "C"},
	}
	g, err := stgui.NewGrid(data, nil)
	if err != nil {
		t.Fatalf("NewGrid failed on ragged: %v", err)
	}

	w, h := g.Size()
	if w != 2 || h != 2 {
		t.Errorf("Expected size 2x2 (max width), got %dx%d", w, h)
	}

	if g.Cells[0][1].Value != nil {
		t.Errorf("Expected cell 0,1 (missing in input) to be nil, got %v", g.Cells[0][1].Value)
	}
	if g.Cells[1][1].Value != "C" {
		t.Errorf("Expected cell 1,1 to be 'C', got %v", g.Cells[1][1].Value)
	}
}

func TestGrid_Render(t *testing.T) {
	data := [][]any{
		{"A", "B"},
	}
	style := &stgui.GridStyle{
		VerticalDivider:   "|",
		HorizontalDivider: "-",
		Intersection:      "+",
		NoValue:           ".",
	}
	g, _ := stgui.NewGrid(data, style)

	output := g.RenderLines()
	if len(output) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(output))
	}

	if !strings.Contains(output[0], "A") || !strings.Contains(output[0], "|") || !strings.Contains(output[0], "B") {
		t.Errorf("Render output unexpected: %q", output[0])
	}
}

func TestGrid_Render_MultiRow(t *testing.T) {
	data := [][]any{
		{"A", "B"},
		{"C", "D"},
	}
	style := stgui.WithGridSymbols()
	g, _ := stgui.NewGrid(data, style)

	output := g.RenderLines()
	if len(output) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(output))
	}

	if !strings.Contains(output[1], "-") {
		t.Error("Divider line missing between rows")
	}
}

// Used to panic. The output was sized as if every row were one line tall, so any
// widget taller than that ran off the end.
func TestGrid_Render_TallCell(t *testing.T) {
	tall := &testutil.MockRenderable{Lines: []string{"1", "2", "3"}}
	g, _ := stgui.NewGrid([][]any{
		{tall, "x"},
		{"C", "D"},
	}, stgui.WithGridSymbols())

	output := g.RenderLines()

	// 3 for the tall row, 1 divider, 1 for the second row.
	if len(output) != 5 {
		t.Fatalf("Expected 5 lines, got %d: %q", len(output), output)
	}
	if !strings.HasPrefix(output[0], "1") || !strings.HasPrefix(output[2], "3") {
		t.Errorf("Tall cell not laid out over its row: %q", output)
	}
	if !strings.Contains(output[3], "-") {
		t.Errorf("Divider not drawn after the tall row: %q", output)
	}

	// All the same width, or the dividers won't line up.
	for i, line := range output {
		if len(line) != len(output[0]) {
			t.Errorf("Line %d width %d, expected %d: %q", i, len(line), len(output[0]), line)
		}
	}
}

// No divider characters, no divider lines - not blank ones.
func TestGrid_Render_NoStyleHasNoBlankRows(t *testing.T) {
	g, _ := stgui.NewGrid([][]any{
		{"A", "B"},
		{"C", "D"},
	}, nil)

	output := g.RenderLines()
	if len(output) != 2 {
		t.Fatalf("Expected 2 lines for an unstyled grid, got %d: %q", len(output), output)
	}
}

func TestGrid_Refresh(t *testing.T) {
	g, _ := stgui.NewGrid([][]any{{"A"}}, nil)
	err := g.Refresh([][]any{{"B", "C"}})
	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}

	if len(g.Cells) != 1 || len(g.Cells[0]) != 2 {
		t.Errorf("Cells dimensions not updated correctly. Got %dx%d", len(g.Cells), len(g.Cells[0]))
	}
	if g.Cells[0][0].Value != "B" {
		t.Errorf("Expected first cell to be 'B', got %v", g.Cells[0][0].Value)
	}

	if err := g.Refresh(nil); err != nil {
		t.Fatalf("Refresh failed on empty: %v", err)
	}
	if w, h := g.Size(); w != 0 || h != 0 {
		t.Errorf("Expected an empty grid after refreshing with no rows, got %dx%d", w, h)
	}
}

func TestNewEmptyGrid(t *testing.T) {
	g := stgui.NewEmptyGrid(3, 2, nil)

	w, h := g.Size()
	if w != 2 || h != 3 {
		t.Fatalf("Expected a 2x3 grid, got %dx%d", w, h)
	}
	g.Each(func(row, col int, cell *stgui.Cell) {
		if !cell.IsEmpty() {
			t.Errorf("Cell %d,%d is not empty", row, col)
		}
	})

	if empty := stgui.NewEmptyGrid(0, 4, nil); empty.Cells != nil {
		t.Error("Expected no cells when a dimension is zero")
	}
}

func TestGrid_SetGetAt(t *testing.T) {
	g := stgui.NewEmptyGrid(2, 2, nil)

	g.Set(1, 0, "x")
	if got := g.Get(1, 0); got != "x" {
		t.Errorf("Expected 'x', got %v", got)
	}

	// Keeps whatever styling and selection the cell already had.
	cell := g.At(1, 0)
	cell.Style = stgui.StyleRed
	cell.Selected = true
	g.Set(1, 0, "y")
	if cell.Style != stgui.StyleRed || !cell.Selected {
		t.Error("Set cleared the cell's style or selection")
	}

	child := &testutil.MockRenderable{Lines: []string{"widget"}}
	g.Set(0, 0, child)
	if g.At(0, 0).Child != child {
		t.Error("Expected a Renderable value to become the cell's child")
	}
	if g.At(0, 0).Value != nil {
		t.Error("Expected the previous value to be cleared")
	}

	// Off the grid: ignored, not a panic.
	g.Set(9, 9, "z")
	if g.At(9, 9) != nil || g.Get(-1, 0) != nil {
		t.Error("Expected nil outside the grid")
	}
}

func TestGrid_Each(t *testing.T) {
	g := stgui.NewEmptyGrid(2, 3, nil)

	seen := 0
	g.Each(func(row, col int, cell *stgui.Cell) {
		seen++
		g.Set(row, col, row*3+col)
	})

	if seen != 6 {
		t.Errorf("Expected 6 cells, visited %d", seen)
	}
	if got := g.Get(1, 2); got != 5 {
		t.Errorf("Expected cell 1,2 to be 5, got %v", got)
	}
}

// Style goes on after padding, so a cell is coloured edge to edge and not just
// behind its text.
func TestGrid_StyleFillsCell(t *testing.T) {
	g := stgui.NewEmptyGrid(1, 2, &stgui.GridStyle{MinCellWidth: 3})
	g.At(0, 0).Style = stgui.StyleRed

	line := g.RenderLines()[0]

	if !strings.Contains(line, stgui.StyleRed+"   \033[0m") {
		t.Errorf("Expected the style to cover the whole cell, got %q", line)
	}
}

func TestGrid_MinCellWidthAndAlign(t *testing.T) {
	style := &stgui.GridStyle{MinCellWidth: 5, Align: stgui.AlignCenter}
	g, _ := stgui.NewGrid([][]any{{"A", "B"}}, style)

	// Two five-column cells, touching: the style sets no divider.
	if line := g.RenderLines()[0]; line != "  A    B  " {
		t.Errorf("Expected centred cells 5 wide, got %q", line)
	}

	style.Align = stgui.AlignRight
	if line := g.RenderLines()[0]; line != "    A    B" {
		t.Errorf("Expected right aligned cells, got %q", line)
	}
}

// Box-drawing characters are multi-byte. Measured by byte length, the dividers
// came out wider than the rows they sit between.
func TestGrid_BoxSymbolsAlign(t *testing.T) {
	g := stgui.NewEmptyGrid(2, 3, stgui.WithBoxSymbols(3))
	g.Set(0, 0, "♜")

	lines := g.RenderLines()
	want := stgui.VisibleWidth(lines[0])
	for i, line := range lines {
		if got := stgui.VisibleWidth(line); got != want {
			t.Errorf("Line %d is %d columns wide, expected %d: %q", i, got, want, line)
		}
	}
}

// Empty cells use the style's NoValue, so a board can leave squares blank.
func TestGrid_EmptyCellsUseStyleNoValue(t *testing.T) {
	g := stgui.NewEmptyGrid(1, 2, stgui.WithBoxSymbols(3))

	if line := g.RenderLines()[0]; strings.Contains(line, ".") {
		t.Errorf("Expected blank cells, got %q", line)
	}
}

// Dividers go in verbatim, which is how a board gets its cells to touch.
func TestGrid_DividerDrawnAsWritten(t *testing.T) {
	data := [][]any{{"A", "B"}}

	tight, _ := stgui.NewGrid(data, &stgui.GridStyle{VerticalDivider: "|"})
	if line := tight.RenderLines()[0]; line != "A|B" {
		t.Errorf("Expected 'A|B', got %q", line)
	}

	spaced, _ := stgui.NewGrid(data, &stgui.GridStyle{VerticalDivider: " | "})
	if line := spaced.RenderLines()[0]; line != "A | B" {
		t.Errorf("Expected 'A | B', got %q", line)
	}

	none, _ := stgui.NewGrid(data, &stgui.GridStyle{})
	if line := none.RenderLines()[0]; line != "AB" {
		t.Errorf("Expected touching cells, got %q", line)
	}

	// Nil style: readable without asking for a divider.
	plain, _ := stgui.NewGrid(data, nil)
	if line := plain.RenderLines()[0]; line != "A B" {
		t.Errorf("Expected 'A B' by default, got %q", line)
	}
}

// Some boundaries and not others, the way a sudoku marks its boxes.
func TestGrid_DividerPredicates(t *testing.T) {
	style := &stgui.GridStyle{
		VerticalDivider:   "│",
		HorizontalDivider: "─",
		Intersection:      "┼",
		DrawColDivider:    func(col int) bool { return col%3 == 2 },
		DrawRowDivider:    func(row int) bool { return row%3 == 2 },
	}

	g := stgui.NewEmptyGrid(6, 6, style)
	g.Each(func(row, col int, cell *stgui.Cell) { g.Set(row, col, "x") })

	lines := g.RenderLines()

	// Six rows, plus the divider after row three.
	if len(lines) != 7 {
		t.Fatalf("Expected 7 lines, got %d: %q", len(lines), lines)
	}
	if lines[0] != "xxx│xxx" {
		t.Errorf("Expected a divider only after the third column, got %q", lines[0])
	}
	if lines[3] != "───┼───" {
		t.Errorf("Expected the row divider to cross at the box corner, got %q", lines[3])
	}

	for i, line := range lines {
		if got := stgui.VisibleWidth(line); got != 7 {
			t.Errorf("Line %d is %d columns wide, expected 7: %q", i, got, line)
		}
	}
}

// The white highlight is hard to see on a coloured cell, so boards can swap it.
func TestGrid_SelectedStyleOverride(t *testing.T) {
	style := &stgui.GridStyle{SelectedStyle: stgui.StyleBGGreen}
	g := stgui.NewEmptyGrid(1, 1, style)
	g.Set(0, 0, "x")
	g.At(0, 0).Selected = true

	if line := g.RenderLines()[0]; line != stgui.StyleBGGreen+"x\033[0m" {
		t.Errorf("Expected the style's highlight, got %q", line)
	}
}
