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

func TestCell_RenderLines_Empty(t *testing.T) {
	c := &stgui.Cell{}
	lines := c.RenderLines()
	if len(lines) != 1 || lines[0] != "." {
		t.Errorf("Expected ['.'], got %v", lines)
	}
}
