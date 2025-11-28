package widgets

// Clickable interface for widgets that respond to an activation event (like Enter).
type Clickable interface {
	OnClick()
}

// InputHandler interface for widgets that accept text input.
type InputHandler interface {
	HandleInput(input string)
}
