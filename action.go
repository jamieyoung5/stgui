package stgui

// Action is what a widget asks the app to do once it has handled an input.
// Screens live on a stack: Push puts one on top of the current one, Replace
// swaps it out. The zero value is Stay.
type Action struct {
	kind   actionKind
	screen *Screen
}

type actionKind uint8

const (
	actionStay actionKind = iota
	actionPush
	actionReplace
	actionClose
	actionQuit
)

// Stay redraws the current screen without changing it.
func Stay() Action {
	return Action{kind: actionStay}
}

// Push opens screen on top of the current one, which stays on the stack for
// Close to come back to. A nil screen stays put.
func Push(screen *Screen) Action {
	if screen == nil {
		return Stay()
	}
	return Action{kind: actionPush, screen: screen}
}

// Replace swaps the current screen for screen, with no coming back. A nil
// screen stays put.
func Replace(screen *Screen) Action {
	if screen == nil {
		return Stay()
	}
	return Action{kind: actionReplace, screen: screen}
}

// Close closes the current screen, quitting if it was the last one.
func Close() Action {
	return Action{kind: actionClose}
}

// Quit tears the app down, whatever is on the stack.
func Quit() Action {
	return Action{kind: actionQuit}
}

func (a Action) String() string {
	switch a.kind {
	case actionPush:
		return "Push"
	case actionReplace:
		return "Replace"
	case actionClose:
		return "Close"
	case actionQuit:
		return "Quit"
	default:
		return "Stay"
	}
}
