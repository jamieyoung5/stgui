package widgets_test

import (
	"testing"

	"github.com/jamieyoung5/stgui/widgets"
)

func TestLabel(t *testing.T) {
	lbl := widgets.NewLabel("Hello\nWorld")
	lines := lbl.RenderLines()
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "Hello" || lines[1] != "World" {
		t.Errorf("Unexpected content: %v", lines)
	}
}
