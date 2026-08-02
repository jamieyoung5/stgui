package widgets_test

import (
	"testing"

	"github.com/jamieyoung5/stgui/widgets"
)

func TestText(t *testing.T) {
	txt := widgets.NewText(widgets.NewLabel("status\nline two"))

	if got := txt.Render(); got != "status\nline two" {
		t.Errorf("Unexpected render output: %q", got)
	}

	w, h := txt.Size()
	if w != 8 || h != 2 {
		t.Errorf("Expected size 8x2, got %dx%d", w, h)
	}

	// Not interactive, so input goes straight through.
	if screen, exit := txt.Select(nil, "x"); screen != nil || exit {
		t.Error("Expected text to ignore input")
	}
}

func TestText_Empty(t *testing.T) {
	txt := widgets.NewText(nil)

	if got := txt.Render(); got != "" {
		t.Errorf("Expected empty output, got %q", got)
	}
	if w, h := txt.Size(); w != 0 || h != 0 {
		t.Errorf("Expected size 0x0, got %dx%d", w, h)
	}
}
