package stgui_test

import (
	"regexp"
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

// Styled widgets used to be padded by byte length, so escape sequences pushed
// the next column out of line.
func TestScreen_Render_StyledColumnsAlign(t *testing.T) {
	styled := &testutil.MockWidget{
		RenderOutput: stgui.StyleRed + "ab" + "\033[0m" + "\n" + "cd",
	}
	right := &testutil.MockWidget{RenderOutput: "X\nY"}

	screen := stgui.NewScreen(nil, stgui.Row(styled, right))

	lines := strings.Split(screen.Render(), "\n")
	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines, got %d: %q", len(lines), lines)
	}

	first := strings.Index(stripANSI(lines[0]), "X")
	second := strings.Index(stripANSI(lines[1]), "Y")
	if first != second {
		t.Errorf("Expected the second column to start at the same place on both lines, got %d and %d in %q", first, second, lines)
	}
}

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

// A row with nothing in it draws nothing, rather than spacing the rows either
// side of it further apart.
func TestScreen_Render_EmptyRow(t *testing.T) {
	screen := stgui.NewScreen(nil,
		stgui.Row(&testutil.MockWidget{RenderOutput: "top"}),
		stgui.Row(nil),
		stgui.Row(&testutil.MockWidget{RenderOutput: "bottom"}),
	)

	if got, want := screen.Render(), "top\n\nbottom"; got != want {
		t.Errorf("Expected %q, got %q", want, got)
	}
}
