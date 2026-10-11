package stgui

import "github.com/jamieyoung5/stgui/keyboard"

// Key is a key press, aliased so binding keys doesn't mean importing keyboard.
type Key = keyboard.Key

// The named keys, re-exported. Anything else is the characters typed: Key("q").
const (
	KeyEnter     = keyboard.KeyEnter
	KeyEscape    = keyboard.KeyEscape
	KeyBackspace = keyboard.KeyBackspace
	KeyTab       = keyboard.KeyTab
	KeyCtrlC     = keyboard.KeyCtrlC
	KeyUp        = keyboard.KeyUp
	KeyDown      = keyboard.KeyDown
	KeyLeft      = keyboard.KeyLeft
	KeyRight     = keyboard.KeyRight
)

// Control is what a key means. Widgets act on controls, not keys, which keeps
// bindings swappable. Controls are lowercase, keys uppercase.
type Control string

const (
	ControlUp    Control = "up"
	ControlDown  Control = "down"
	ControlLeft  Control = "left"
	ControlRight Control = "right"

	// ControlSelect activates whatever the cursor is on. Enter by default.
	ControlSelect Control = "select"
)

// Bindings maps keys onto controls. An unbound key still reaches widgets, with
// no control attached.
type Bindings map[Key]Control

// DefaultBindings is the arrow keys plus Enter to activate.
func DefaultBindings() Bindings {
	return Bindings{
		KeyUp:    ControlUp,
		KeyDown:  ControlDown,
		KeyLeft:  ControlLeft,
		KeyRight: ControlRight,
		KeyEnter: ControlSelect,
	}
}

// Input is one key press as a widget sees it. Control is empty for an unbound
// key.
type Input struct {
	Key     Key
	Control Control
}

// Is reports whether the input is bound to control.
func (i Input) Is(control Control) bool {
	return i.Control != "" && i.Control == control
}

// Rune is the character typed, if there was one. Named keys, escape sequences
// and bound keys report false.
func (i Input) Rune() (rune, bool) {
	if i.Control != "" {
		return 0, false
	}
	return i.Key.Rune()
}
