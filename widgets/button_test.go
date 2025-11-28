package widgets_test

import (
	"testing"

	"github.com/jamieyoung5/stgui/widgets"
)

func TestButton(t *testing.T) {
	clicked := false
	btn := widgets.NewButton("Click Me", func() {
		clicked = true
	})

	lines := btn.RenderLines()
	if len(lines) != 1 || lines[0] != "[ Click Me ]" {
		t.Errorf("Unexpected RenderLines output: %v", lines)
	}

	btn.OnClick()
	if !clicked {
		t.Error("Callback not triggered on OnClick")
	}

	safeBtn := widgets.NewButton("Safe", nil)
	safeBtn.OnClick()
}
