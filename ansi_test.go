package stgui

import "testing"

func TestVisibleWidth(t *testing.T) {
	cases := map[string]int{
		"":                                 0,
		"abc":                              3,
		selectedStyle + "abc" + resetStyle: 3,
		"♜":                                1,
		"♜♞♝":                              3,
		"éüø":                              3,
		StyleRed + "♜" + resetStyle:        1,
	}

	for input, want := range cases {
		if got := VisibleWidth(input); got != want {
			t.Errorf("VisibleWidth(%q) = %d, want %d", input, got, want)
		}
	}
}

// One column per rune, so a line of chess pieces gets cut between pieces and not
// through one.
func TestTruncateVisible_Unicode(t *testing.T) {
	if got := truncateVisible("♜♞♝♛♚", 3); got != "♜♞♝" {
		t.Errorf("Expected '♜♞♝', got %q", got)
	}
	if got := truncateVisible("♜♞♝", 5); got != "♜♞♝" {
		t.Errorf("Expected a short line unchanged, got %q", got)
	}
}

func TestAnsiSequenceAt(t *testing.T) {
	if got := ansiSequenceAt("\033[31mred"); got != "\033[31m" {
		t.Errorf("Expected the escape sequence, got %q", got)
	}
	if got := ansiSequenceAt("red"); got != "" {
		t.Errorf("Expected no sequence for plain text, got %q", got)
	}
}
