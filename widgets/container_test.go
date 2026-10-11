package widgets_test

import (
	"strings"
	"testing"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/widgets"
)

// typed is an unbound key press.
func typed(key stgui.Key) stgui.Input {
	return stgui.Input{Key: key}
}

// bound is a key press that resolved to a control.
func bound(key stgui.Key, control stgui.Control) stgui.Input {
	return stgui.Input{Key: key, Control: control}
}

// selected is Enter doing its default job.
var selected = bound(stgui.KeyEnter, stgui.ControlSelect)

func TestContainer(t *testing.T) {
	data := [][]any{{"A"}}
	g := stgui.NewGrid(data, nil)
	c := widgets.NewContainer(g)

	w, h := c.Size()
	if w != 1 || h != 1 {
		t.Errorf("Expected size 1x1, got %dx%d", w, h)
	}

	out := c.Render()
	if !strings.Contains(out, "A") {
		t.Error("Render output missing grid content")
	}

	clicked := false
	btn := widgets.NewButton("Btn", func() { clicked = true })
	gClick := stgui.NewGrid([][]any{{btn}}, nil)
	cClick := widgets.NewContainer(gClick)
	cursorClick := stgui.NewCursor(cClick, gClick, 0, 0)

	cClick.Select(cursorClick, selected)
	if !clicked {
		t.Error("Container did not delegate ControlSelect to Clickable child")
	}

	inp := widgets.NewInput("txt", 10)
	gInp := stgui.NewGrid([][]any{{inp}}, nil)
	cInp := widgets.NewContainer(gInp)
	cursorInp := stgui.NewCursor(cInp, gInp, 0, 0)

	cInp.Select(cursorInp, typed("X"))
	if inp.Value != "X" {
		t.Error("Container did not delegate input to InputHandler child")
	}
}

func TestContainer_Navigation(t *testing.T) {
	g := stgui.NewGrid([][]any{{"A", "B"}, {"C", "D"}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0)

	c.Select(cursor, bound(stgui.KeyRight, stgui.ControlRight))
	c.Select(cursor, bound(stgui.KeyDown, stgui.ControlDown))
	if cursor.Row != 1 || cursor.Col != 1 {
		t.Fatalf("Expected the cursor at 1,1, got %d,%d", cursor.Row, cursor.Col)
	}

	c.Select(cursor, bound(stgui.KeyLeft, stgui.ControlLeft))
	c.Select(cursor, bound(stgui.KeyUp, stgui.ControlUp))
	if cursor.Row != 0 || cursor.Col != 0 {
		t.Fatalf("Expected the cursor back at 0,0, got %d,%d", cursor.Row, cursor.Col)
	}
}

// Controls are what move the cursor, not keys: an unbound arrow key does
// nothing, and a letter bound to a direction moves.
func TestContainer_MovesOnControlNotKey(t *testing.T) {
	g := stgui.NewGrid([][]any{{"A", "B"}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0)

	c.Select(cursor, typed(stgui.KeyRight))
	if cursor.Col != 0 {
		t.Error("Expected an unbound right arrow not to move the cursor")
	}

	c.Select(cursor, bound("d", stgui.ControlRight))
	if cursor.Col != 1 {
		t.Errorf("Expected 'd' bound to right to move the cursor, col is %d", cursor.Col)
	}
}

// A Navigator child decides where the app goes next, and the container passes
// that decision straight back to the app.
func TestContainer_DelegatesToNavigator(t *testing.T) {
	next := stgui.NewScreen(nil)
	navBtn := widgets.NewNavButton("Go", next)
	quitBtn := widgets.NewQuitButton("Quit")

	g := stgui.NewGrid([][]any{{navBtn, quitBtn}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0)

	if got := c.Select(cursor, selected); got != stgui.Push(next) {
		t.Errorf("Expected the nav button's action, got %v", got)
	}

	cursor.Right()
	if got := c.Select(cursor, selected); got != stgui.Close() {
		t.Errorf("Expected the quit button to close the screen, got %v", got)
	}
}

// You can put the cursor on an empty cell or a plain string, so typing at one has
// to do nothing rather than fall over.
func TestContainer_SelectOnNonWidgetCell(t *testing.T) {
	g := stgui.NewGrid([][]any{{"just text", nil}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0)

	if got := c.Select(cursor, selected); got != stgui.Stay() {
		t.Errorf("Expected selecting a plain value to do nothing, got %v", got)
	}

	cursor.Right()
	if got := c.Select(cursor, typed("x")); got != stgui.Stay() {
		t.Errorf("Expected input on an empty cell to do nothing, got %v", got)
	}
}

func TestNewGridScreen(t *testing.T) {
	g := stgui.NewGrid([][]any{{"A", "B"}}, stgui.WithGridSymbols())

	screen := widgets.NewGridScreen(g, 0, 1)

	cursor := screen.Cursor
	if cursor == nil {
		t.Fatal("Expected the screen to have a cursor")
	}
	if cursor.Row != 0 || cursor.Col != 1 {
		t.Errorf("Expected the cursor at 0,1, got %d,%d", cursor.Row, cursor.Col)
	}
	if !strings.Contains(screen.Render(), "A") {
		t.Errorf("Expected the grid to be rendered, got %q", screen.Render())
	}
}

// OnKey sees the input and the cursor, so one handler covers every cell.
func TestContainer_OnKey(t *testing.T) {
	g := stgui.NewGrid([][]any{{"a", "b"}, {"c", "d"}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0)

	type press struct {
		input    stgui.Input
		row, col int
	}
	var seen []press
	c.OnKey(func(cursor *stgui.Cursor, input stgui.Input) stgui.Action {
		seen = append(seen, press{input, cursor.Row, cursor.Col})
		return stgui.Stay()
	})

	// The container keeps navigation for itself.
	c.Select(cursor, bound(stgui.KeyRight, stgui.ControlRight))
	c.Select(cursor, typed("5"))
	c.Select(cursor, selected)

	if len(seen) != 2 {
		t.Fatalf("Expected 2 inputs to reach OnKey, got %d: %v", len(seen), seen)
	}
	if seen[0] != (press{typed("5"), 0, 1}) {
		t.Errorf("Expected '5' at the moved-to cell, got %v", seen[0])
	}
	if !seen[1].input.Is(stgui.ControlSelect) {
		t.Errorf("Expected select to reach OnKey on a plain cell, got %v", seen[1])
	}
}

// Whatever is in the focused cell gets first go. OnKey only picks up leftovers.
func TestContainer_OnKeyIsAFallback(t *testing.T) {
	input := widgets.NewInput("", 5)
	btn := widgets.NewButton("Btn", nil)
	g := stgui.NewGrid([][]any{{input, btn}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0)

	called := 0
	c.OnKey(func(cursor *stgui.Cursor, in stgui.Input) stgui.Action {
		called++
		return stgui.Stay()
	})

	c.Select(cursor, typed("x"))
	if input.Value != "x" || called != 0 {
		t.Errorf("Expected the input to take the key, value %q, OnKey calls %d", input.Value, called)
	}

	cursor.Right()
	c.Select(cursor, selected)
	if called != 0 {
		t.Error("Expected the button to take the selection")
	}
}

// A board quits or moves on by returning an action from OnKey.
func TestContainer_OnKeyCanNavigate(t *testing.T) {
	g := stgui.NewGrid([][]any{{"a"}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0)

	next := stgui.NewScreen(nil)
	c.OnKey(func(cursor *stgui.Cursor, input stgui.Input) stgui.Action {
		return stgui.Replace(next)
	})

	if got := c.Select(cursor, typed("x")); got != stgui.Replace(next) {
		t.Errorf("Expected OnKey's action to be returned, got %v", got)
	}
}
