package widgets

import (
	"strings"
	"unicode/utf8"

	"github.com/jamieyoung5/stgui/keyboard"
)

// Input is a one-line text field. It collects what you type while its cell has
// focus, and shows the placeholder until then.
type Input struct {
	Value       string
	Placeholder string
	Width       int

	// Mask replaces every character of Value on screen. For passwords.
	Mask rune
}

func NewInput(placeholder string, width int) *Input {
	return &Input{
		Value:       "",
		Placeholder: placeholder,
		Width:       width,
	}
}

// NewMaskedInput is an Input that shows asterisks.
func NewMaskedInput(placeholder string, width int) *Input {
	input := NewInput(placeholder, width)
	input.Mask = '*'
	return input
}

func (i *Input) RenderLines() []string {
	display := i.Value
	if display != "" && i.Mask != 0 {
		display = strings.Repeat(string(i.Mask), utf8.RuneCountInString(display))
	}
	if display == "" {
		display = i.Placeholder
	}

	runes := []rune(display)
	if len(runes) > i.Width {
		// Show the end - that's where the typing is.
		display = string(runes[len(runes)-i.Width:])
	} else {
		display += strings.Repeat("_", i.Width-len(runes))
	}

	return []string{display}
}

// HandleInput adds a character, or takes one off on backspace. Named keys and
// escape sequences are ignored.
func (i *Input) HandleInput(input string) {
	if input == keyboard.BackspaceKey {
		if runes := []rune(i.Value); len(runes) > 0 {
			i.Value = string(runes[:len(runes)-1])
		}
		return
	}

	if strings.HasPrefix(input, "\x1b") {
		return
	}

	if utf8.RuneCountInString(input) == 1 {
		i.Value += input
	}
}
