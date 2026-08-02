package widgets_test

import (
	"testing"

	"github.com/jamieyoung5/stgui/keyboard"
	"github.com/jamieyoung5/stgui/widgets"
)

func TestInput_Typing(t *testing.T) {
	in := widgets.NewInput("name", 6)

	if lines := in.RenderLines(); lines[0] != "name__" {
		t.Errorf("Expected the placeholder padded to the width, got %q", lines[0])
	}

	for _, key := range []string{"h", "i"} {
		in.HandleInput(key)
	}
	if in.Value != "hi" {
		t.Fatalf("Expected 'hi', got %q", in.Value)
	}
	if lines := in.RenderLines(); lines[0] != "hi____" {
		t.Errorf("Expected 'hi____', got %q", lines[0])
	}

	in.HandleInput(keyboard.BackspaceKey)
	if in.Value != "h" {
		t.Errorf("Expected backspace to remove the last character, got %q", in.Value)
	}

	in.HandleInput(keyboard.BackspaceKey)
	in.HandleInput(keyboard.BackspaceKey)
	if in.Value != "" {
		t.Errorf("Expected backspace on an empty value to be safe, got %q", in.Value)
	}
}

func TestInput_IgnoresNamedKeys(t *testing.T) {
	in := widgets.NewInput("", 6)

	for _, key := range []string{keyboard.EnterKey, keyboard.TabKey, keyboard.UpArrowKey, "\x1b[Z"} {
		in.HandleInput(key)
	}

	if in.Value != "" {
		t.Errorf("Expected named keys and escape sequences to be ignored, got %q", in.Value)
	}
}

// Runes, not bytes: multi-byte characters used to get dropped on the way in and
// sliced in half on the way out.
func TestInput_MultiByteCharacters(t *testing.T) {
	in := widgets.NewInput("", 3)

	for _, key := range []string{"é", "ü", "ø"} {
		in.HandleInput(key)
	}
	if in.Value != "éüø" {
		t.Fatalf("Expected 'éüø', got %q", in.Value)
	}

	if lines := in.RenderLines(); lines[0] != "éüø" {
		t.Errorf("Expected 'éüø' to fill the width, got %q", lines[0])
	}

	in.HandleInput(keyboard.BackspaceKey)
	if in.Value != "éü" {
		t.Errorf("Expected backspace to remove one whole rune, got %q", in.Value)
	}
}

func TestInput_ScrollsToEnd(t *testing.T) {
	in := widgets.NewInput("", 3)
	in.Value = "abcdef"

	if lines := in.RenderLines(); lines[0] != "def" {
		t.Errorf("Expected the end of the value to be shown, got %q", lines[0])
	}
}

func TestInput_Masked(t *testing.T) {
	in := widgets.NewMaskedInput("secret", 6)

	if lines := in.RenderLines(); lines[0] != "secret" {
		t.Errorf("Expected the placeholder to show unmasked, got %q", lines[0])
	}

	in.HandleInput("a")
	in.HandleInput("b")

	if lines := in.RenderLines(); lines[0] != "**____" {
		t.Errorf("Expected the value to be masked, got %q", lines[0])
	}
	if in.Value != "ab" {
		t.Errorf("Expected the underlying value to be readable, got %q", in.Value)
	}
}
