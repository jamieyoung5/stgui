package widgets_test

import (
	"testing"

	"github.com/jamieyoung5/stgui"
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

// A label doubles as a top-level widget, so a title or status line doesn't need
// wrapping in anything.
func TestLabel_AsWidget(t *testing.T) {
	lbl := widgets.NewLabel("status\nline two")

	if got := lbl.Render(); got != "status\nline two" {
		t.Errorf("Unexpected render output: %q", got)
	}

	// Not interactive, so input passes straight through.
	if got := lbl.Select(nil, typed("x")); got != stgui.Stay() {
		t.Errorf("Expected a label to ignore input, got %v", got)
	}

	if got := widgets.NewLabel("").Render(); got != "" {
		t.Errorf("Expected an empty label to render nothing, got %q", got)
	}
}
