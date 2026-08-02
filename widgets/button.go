package widgets

import (
	"fmt"

	"github.com/jamieyoung5/stgui"
)

// Button runs a callback on Enter, and can move to another screen or close the
// current one on the way.
type Button struct {
	Label    string
	Callback func()

	// Screen to go to when pressed. It replaces the current screen, or stacks
	// on top of it if that one has Persist set.
	Screen *stgui.Screen
	// Exit closes the current screen, quitting if it's the last one open.
	Exit bool
}

func NewButton(label string, callback func()) *Button {
	return &Button{
		Label:    label,
		Callback: callback,
	}
}

// NewNavButton makes a button that takes you to next.
func NewNavButton(label string, next *stgui.Screen) *Button {
	return &Button{
		Label:  label,
		Screen: next,
	}
}

// NewQuitButton makes a button that closes the current screen. Quit this way
// rather than with os.Exit, or the terminal is left in raw mode.
func NewQuitButton(label string) *Button {
	return &Button{
		Label: label,
		Exit:  true,
	}
}

func (b *Button) RenderLines() []string {
	return []string{fmt.Sprintf("[ %s ]", b.Label)}
}

func (b *Button) OnClick() {
	if b.Callback != nil {
		b.Callback()
	}
}

// OnActivate runs the callback first, then reads Screen and Exit - so a callback
// can decide where the button goes.
func (b *Button) OnActivate() (*stgui.Screen, bool) {
	b.OnClick()
	return b.Screen, b.Exit || b.Screen != nil
}
