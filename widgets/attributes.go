package widgets

import "github.com/jamieyoung5/stgui"

// Clickable is a widget that does something when you press Enter on it.
type Clickable interface {
	OnClick()
}

// Navigator is a widget that changes screens when you press Enter on it. Return
// a screen to go there; return exit to close the current one, which quits the app
// if it's the last.
//
// Navigator wins over Clickable, so if a widget is both, OnActivate has to do the
// click as well.
type Navigator interface {
	OnActivate() (next *stgui.Screen, exit bool)
}

// InputHandler is a widget you can type into.
type InputHandler interface {
	HandleInput(input string)
}
