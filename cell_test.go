package stgui_test

import (
	"testing"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/testutil"
)

func TestCell_RenderLines_Value(t *testing.T) {
	c := &stgui.Cell{Value: "test"}
	lines := c.RenderLines()
	if len(lines) != 1 || lines[0] != "test" {
		t.Errorf("Expected ['test'], got %v", lines)
	}
}

func TestCell_RenderLines_Child(t *testing.T) {
	child := &testutil.MockRenderable{Lines: []string{"line1", "line2"}}
	c := &stgui.Cell{Child: child}
	lines := c.RenderLines()
	if len(lines) != 2 || lines[0] != "line1" || lines[1] != "line2" {
		t.Errorf("Expected ['line1', 'line2'], got %v", lines)
	}
}

func TestCell_RenderLines_Selected(t *testing.T) {
	c := &stgui.Cell{Value: "sel", Selected: true}
	lines := c.RenderLines()
	// Expected: whiteBGBlack + "sel" + resetStyle
	// whiteBGBlack = "\033[47m\033[30m", resetStyle = "\033[0m"
	expected := "\033[47m\033[30m" + "sel" + "\033[0m"
	if len(lines) != 1 || lines[0] != expected {
		t.Errorf("Expected [%q], got %v", expected, lines)
	}
}

// Newlines have to be split out, or the grid counts the value as one line and the
// layout stops matching what's on screen.
func TestCell_RenderLines_MultiLineValue(t *testing.T) {
	c := &stgui.Cell{Value: "one\ntwo"}
	lines := c.RenderLines()
	if len(lines) != 2 || lines[0] != "one" || lines[1] != "two" {
		t.Errorf("Expected ['one', 'two'], got %q", lines)
	}
}

func TestCell_RenderLines_EmptyChild(t *testing.T) {
	c := &stgui.Cell{Child: &testutil.MockRenderable{}}
	lines := c.RenderLines()
	if len(lines) != 1 || lines[0] != "." {
		t.Errorf("Expected ['.'] for a child with no lines, got %q", lines)
	}
}

func TestCell_RenderLines_Empty(t *testing.T) {
	c := &stgui.Cell{}
	lines := c.RenderLines()
	if len(lines) != 1 || lines[0] != "." {
		t.Errorf("Expected ['.'], got %v", lines)
	}
}

func TestCell_Style(t *testing.T) {
	c := &stgui.Cell{Value: "x", Style: stgui.StyleRed}

	if got := c.RenderLines()[0]; got != stgui.StyleRed+"x\033[0m" {
		t.Errorf("Expected the cell's style to be applied, got %q", got)
	}

	// Selection wins.
	c.Selected = true
	if got := c.StyleSeq(); got == stgui.StyleRed {
		t.Errorf("Expected selection to override the cell's style, got %q", got)
	}
}

func TestCell_IsEmpty(t *testing.T) {
	if !(&stgui.Cell{}).IsEmpty() {
		t.Error("Expected a cell with no contents to be empty")
	}
	if (&stgui.Cell{Value: 0}).IsEmpty() {
		t.Error("Expected a cell holding a value not to be empty")
	}
}
