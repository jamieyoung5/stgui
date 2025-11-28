package widgets

import (
	"strings"
)

type Input struct {
	Value       string
	Placeholder string
	Width       int
}

func NewInput(placeholder string, width int) *Input {
	return &Input{
		Value:       "",
		Placeholder: placeholder,
		Width:       width,
	}
}

func (i *Input) RenderLines() []string {
	display := i.Value
	if display == "" {
		display = i.Placeholder
	}

	if len(display) > i.Width {
		display = display[len(display)-i.Width:]
	} else {
		display = display + strings.Repeat("_", i.Width-len(display))
	}

	return []string{display}
}

func (i *Input) HandleInput(input string) {
	if input == "\x7f" || input == "\b" {
		if len(i.Value) > 0 {
			i.Value = i.Value[:len(i.Value)-1]
		}
		return
	}

	if strings.HasPrefix(input, "\x1b") {
		return
	}

	if len(input) == 1 {
		i.Value += input
	}
}
