package stgui_test

import (
	"strings"
	"testing"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/testutil"
)

func TestScreen_Cursor(t *testing.T) {
	cursor := stgui.NewCursor(nil, nil, 0, 0)

	if got := stgui.NewScreen(cursor).Cursor; got != cursor {
		t.Errorf("Expected the screen to hold the cursor it was built with, got %v", got)
	}

	if got := stgui.NewScreen(nil).Cursor; got != nil {
		t.Errorf("Expected no cursor on a screen built without one, got %v", got)
	}
}

func TestScreen_Render(t *testing.T) {
	left := &testutil.MockWidget{RenderOutput: "left"}
	right := &testutil.MockWidget{RenderOutput: "right"}
	below := &testutil.MockWidget{RenderOutput: "below"}

	screen := stgui.NewScreen(nil, stgui.Row(left, right), stgui.Row(below))

	lines := strings.Split(screen.Render(), "\n")

	// Side by side within a row, blank line between rows, nothing trailing off
	// the end.
	want := []string{"left    right", "", "below"}
	if len(lines) != len(want) {
		t.Fatalf("Expected %d lines, got %d: %q", len(want), len(lines), lines)
	}
	for i, line := range want {
		if lines[i] != line {
			t.Errorf("Line %d: expected %q, got %q", i, line, lines[i])
		}
	}
}

func TestScreen_Render_Empty(t *testing.T) {
	if got := stgui.NewScreen(nil).Render(); got != "" {
		t.Errorf("Expected an empty screen to render nothing, got %q", got)
	}
}
