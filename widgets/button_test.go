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

	btn.OnClick()
	if !clicked {
		t.Error("Callback not triggered on OnClick")
	}

	safeBtn := widgets.NewButton("Safe", nil)
	safeBtn.OnClick()

	// Plain button: runs the callback, stays put.
	if screen, exit := btn.OnActivate(); screen != nil || exit {
		t.Errorf("Expected a plain button not to change screens, got %v exit=%v", screen, exit)
	}
}

func TestQuitButton(t *testing.T) {
	screen, exit := widgets.NewQuitButton("Quit").OnActivate()
	if screen != nil || !exit {
		t.Errorf("Expected the quit button to close the screen, got %v exit=%v", screen, exit)
	}
}

func TestNavButton(t *testing.T) {
	next := stgui.NewScreen(nil, nil)
	btn := widgets.NewNavButton("Go", next)

	screen, exit := btn.OnActivate()
	if screen != next {
		t.Errorf("Expected the button's screen, got %v", screen)
	}
	if !exit {
		t.Error("Expected navigating to close the current screen")
	}
}

// Callback first, so it can pick where to go based on what the press just did.
func TestButton_CallbackSetsScreen(t *testing.T) {
	next := stgui.NewScreen(nil, nil)
	btn := widgets.NewButton("Login", nil)
	btn.Callback = func() { btn.Screen = next }

	if screen, exit := btn.OnActivate(); screen != next || !exit {
		t.Errorf("Expected the callback's screen, got %v exit=%v", screen, exit)
	}
}
