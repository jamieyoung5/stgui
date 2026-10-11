package widgets

import (
	"fmt"

	"github.com/jamieyoung5/stgui"
)

// Button does something when selected. What it does is fixed when you make it.
type Button struct {
	// Label is what's drawn.
	Label string

	activate func() stgui.Action
}

// NewButton makes a button that runs callback and stays on this screen.
func NewButton(label string, callback func()) *Button {
	return NewActionButton(label, func() stgui.Action {
		if callback != nil {
			callback()
		}
		return stgui.Stay()
	})
}

// NewActionButton makes a button that decides where the app goes when pressed.
func NewActionButton(label string, activate func() stgui.Action) *Button {
	return &Button{Label: label, activate: activate}
}

// NewNavButton makes a button that opens next on top of the current screen.
func NewNavButton(label string, next *stgui.Screen) *Button {
	return NewActionButton(label, func() stgui.Action {
		return stgui.Push(next)
	})
}

// NewQuitButton makes a button that closes the current screen. Quit this way
// rather than os.Exit, or the terminal is left in raw mode.
func NewQuitButton(label string) *Button {
	return NewActionButton(label, stgui.Close)
}

func (b *Button) RenderLines() []string {
	return []string{fmt.Sprintf("[ %s ]", b.Label)}
}

func (b *Button) OnActivate() stgui.Action {
	if b.activate == nil {
		return stgui.Stay()
	}
	return b.activate()
}
