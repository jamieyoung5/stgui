package stgui

import "github.com/jamieyoung5/stgui/keyboard"

// Key is a key press, aliased from the keyboard package so binding keys doesn't
// mean importing it.
type Key = keyboard.Key

// The named keys, re-exported. Anything else is the characters typed, so the
// "q" key is Key("q").
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

// Control is what a key means, as opposed to which key it was. Widgets act on
// controls, not keys, which is what keeps bindings swappable.
//
// Controls are lowercase and keys are uppercase, so the two can't be mistaken
// for each other.
type Control string

const (
	ControlUp    Control = "up"
	ControlDown  Control = "down"
	ControlLeft  Control = "left"
	ControlRight Control = "right"

	// ControlSelect activates whatever the cursor is on. Enter by default.
	ControlSelect Control = "select"
)

// Bindings maps keys onto controls. A key that isn't in here still reaches
// widgets, just with no control attached, which is how typing works.
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

// Input is one key press as a widget sees it. Control is empty when the key
// isn't bound to anything.
type Input struct {
	Key     Key
	Control Control
}

// Is reports whether the input is bound to control. An unbound key is never Is
// anything.
func (i Input) Is(control Control) bool {
	return i.Control != "" && i.Control == control
}

// Rune is the character typed, and whether there was one. Named keys, escape
// sequences and bound keys report false, so a key doing a job in the interface
// doesn't get typed into a text field as well.
func (i Input) Rune() (rune, bool) {
	if i.Control != "" {
		return 0, false
	}
	return i.Key.Rune()
}
