package widgets_test

import (
	"testing"

	"github.com/jamieyoung5/stgui"
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

	// Plain button: runs the callback, stays put.
	if got := btn.OnActivate(); got != stgui.Stay() {
		t.Errorf("Expected a plain button to stay put, got %v", got)
	}
	if !clicked {
		t.Error("Callback not triggered on activate")
	}

	// A button with no callback at all mustn't fall over.
	if got := widgets.NewButton("Safe", nil).OnActivate(); got != stgui.Stay() {
		t.Errorf("Expected a callback-less button to stay put, got %v", got)
	}
}

func TestQuitButton(t *testing.T) {
	if got := widgets.NewQuitButton("Quit").OnActivate(); got != stgui.Close() {
		t.Errorf("Expected the quit button to close the screen, got %v", got)
	}
}

func TestNavButton(t *testing.T) {
	next := stgui.NewScreen(nil)
	btn := widgets.NewNavButton("Go", next)

	if got := btn.OnActivate(); got != stgui.Push(next) {
		t.Errorf("Expected the button to push its screen, got %v", got)
	}
}

// An action button picks where to go when it's pressed.
func TestActionButton(t *testing.T) {
	next := stgui.NewScreen(nil)
	allow := false

	btn := widgets.NewActionButton("Login", func() stgui.Action {
		if !allow {
			return stgui.Stay()
		}
		return stgui.Push(next)
	})

	if got := btn.OnActivate(); got != stgui.Stay() {
		t.Errorf("Expected the button to stay put while denied, got %v", got)
	}

	allow = true
	if got := btn.OnActivate(); got != stgui.Push(next) {
		t.Errorf("Expected the button to push once allowed, got %v", got)
	}
}
