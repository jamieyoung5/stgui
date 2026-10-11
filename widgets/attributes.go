package widgets

import "github.com/jamieyoung5/stgui"

// Clickable is a widget that does something when the cursor selects it.
type Clickable interface {
	OnClick()
}

// Navigator is a widget that decides where the app goes when the cursor selects
// it - Push a screen, Close this one, or Stay put.
//
// Navigator wins over Clickable, so if a widget is both, OnActivate has to do
// the click as well.
type Navigator interface {
	OnActivate() stgui.Action
}

// InputHandler is a widget you can type into. It gets every input the cursor
// didn't use for navigation.
type InputHandler interface {
	HandleInput(input stgui.Input)
}
