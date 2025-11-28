package stgui_test

import (
	"strings"
	"testing"

	"github.com/jamieyoung5/stgui"
)

func TestGrid_NewGrid(t *testing.T) {
	data := [][]any{
		{"A", "B"},
		{"C", "D"},
	}
	g, err := stgui.NewGrid(data, nil)
	if err != nil {
		t.Fatalf("NewGrid failed: %v", err)
	}
	h, w := g.Size()
	if h != 2 || w != 2 {
		t.Errorf("Expected size 2x2, got %dx%d", h, w)
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

func TestGrid_NewGrid_Ragged(t *testing.T) {
	data := [][]any{
		{"A"},
		{"B", "C"},
	}
	g, err := stgui.NewGrid(data, nil)
	if err != nil {
		t.Fatalf("NewGrid failed on ragged: %v", err)
	}

	h, w := g.Size()
	if h != 2 || w != 2 {
		t.Errorf("Expected size 2x2 (max width), got %dx%d", h, w)
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
}
