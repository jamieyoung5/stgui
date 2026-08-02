package widgets_test

import (
	"strings"
	"testing"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/keyboard"
	"github.com/jamieyoung5/stgui/widgets"
)

func TestContainer(t *testing.T) {
	data := [][]any{{"A"}}
	g, _ := stgui.NewGrid(data, nil)
	c := widgets.NewContainer(g)

	w, h := c.Size()
	if w != 1 || h != 1 {
		t.Errorf("Expected size 1x1, got %dx%d", w, h)
	}

	out := c.Render()
	if !strings.Contains(out, "A") {
		t.Error("Render output missing grid content")
	}

	cursor := stgui.NewCursor(c, g, 0, 0, stgui.DefaultDirectionalControls)
	c.Select(cursor, stgui.ControlDown)

	clicked := false
	btn := widgets.NewButton("Btn", func() { clicked = true })
	gClick, _ := stgui.NewGrid([][]any{{btn}}, nil)
	cClick := widgets.NewContainer(gClick)
	cursorClick := stgui.NewCursor(cClick, gClick, 0, 0, nil)

	cClick.Select(cursorClick, keyboard.EnterKey)
	if !clicked {
		t.Error("Container did not delegate ENTER to Clickable child")
	}

	inp := widgets.NewInput("txt", 10)
	gInp, _ := stgui.NewGrid([][]any{{inp}}, nil)
	cInp := widgets.NewContainer(gInp)
	cursorInp := stgui.NewCursor(cInp, gInp, 0, 0, nil)

	cInp.Select(cursorInp, "X")
	if inp.Value != "X" {
		t.Error("Container did not delegate input to InputHandler child")
	}
}

func TestContainer_Navigation(t *testing.T) {
	g, _ := stgui.NewGrid([][]any{{"A", "B"}, {"C", "D"}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0, nil)

	c.Select(cursor, stgui.ControlRight)
	c.Select(cursor, stgui.ControlDown)
	if cursor.Row != 1 || cursor.Col != 1 {
		t.Fatalf("Expected the cursor at 1,1, got %d,%d", cursor.Row, cursor.Col)
	}

	c.Select(cursor, stgui.ControlLeft)
	c.Select(cursor, stgui.ControlUp)
	if cursor.Row != 0 || cursor.Col != 0 {
		t.Fatalf("Expected the cursor back at 0,0, got %d,%d", cursor.Row, cursor.Col)
	}
}

// A Navigator child decides where the app goes next, and the container passes
// that decision straight back to the app.
func TestContainer_DelegatesToNavigator(t *testing.T) {
	next := stgui.NewScreen(nil, nil)
	navBtn := widgets.NewNavButton("Go", next)
	quitBtn := widgets.NewQuitButton("Quit")

	g, _ := stgui.NewGrid([][]any{{navBtn, quitBtn}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0, nil)

	screen, exit := c.Select(cursor, keyboard.EnterKey)
	if screen != next || !exit {
		t.Errorf("Expected the nav button's screen, got %v exit=%v", screen, exit)
	}

	cursor.Right()
	screen, exit = c.Select(cursor, keyboard.EnterKey)
	if screen != nil || !exit {
		t.Errorf("Expected the quit button to close the screen, got %v exit=%v", screen, exit)
	}
}

// You can put the cursor on an empty cell or a plain string, so typing at one has
// to do nothing rather than fall over.
func TestContainer_SelectOnNonWidgetCell(t *testing.T) {
	g, _ := stgui.NewGrid([][]any{{"just text", nil}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0, nil)

	if screen, exit := c.Select(cursor, keyboard.EnterKey); screen != nil || exit {
		t.Error("Expected ENTER on a plain value to do nothing")
	}

	cursor.Right()
	if screen, exit := c.Select(cursor, "x"); screen != nil || exit {
		t.Error("Expected input on an empty cell to do nothing")
	}
}

func TestNewScreen(t *testing.T) {
	g, _ := stgui.NewGrid([][]any{{"A", "B"}}, stgui.WithGridSymbols())

	screen := widgets.NewScreen(g, 0, 1)

	cursor := screen.ActiveCursor()
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

// OnKey is the thing that makes boards work: it sees the key and the cursor, so
// one handler covers every cell.
func TestContainer_OnKey(t *testing.T) {
	g, _ := stgui.NewGrid([][]any{{"a", "b"}, {"c", "d"}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0, nil)

	type press struct {
		input    string
		row, col int
	}
	var seen []press
	c.OnKey = func(cursor *stgui.Cursor, input string) (*stgui.Screen, bool) {
		seen = append(seen, press{input, cursor.Row, cursor.Col})
		return nil, false
	}

	// The container keeps navigation for itself.
	c.Select(cursor, stgui.ControlRight)
	c.Select(cursor, "5")
	c.Select(cursor, keyboard.EnterKey)

	if len(seen) != 2 {
		t.Fatalf("Expected 2 keys to reach OnKey, got %d: %v", len(seen), seen)
	}
	if seen[0] != (press{"5", 0, 1}) {
		t.Errorf("Expected '5' at the moved-to cell, got %v", seen[0])
	}
	if seen[1].input != keyboard.EnterKey {
		t.Errorf("Expected ENTER to reach OnKey on a plain cell, got %v", seen[1])
	}
}

// Whatever is in the focused cell gets first go. OnKey only picks up leftovers.
func TestContainer_OnKeyIsAFallback(t *testing.T) {
	input := widgets.NewInput("", 5)
	btn := widgets.NewButton("Btn", nil)
	g, _ := stgui.NewGrid([][]any{{input, btn}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0, nil)

	called := 0
	c.OnKey = func(cursor *stgui.Cursor, key string) (*stgui.Screen, bool) {
		called++
		return nil, false
	}

	c.Select(cursor, "x")
	if input.Value != "x" || called != 0 {
		t.Errorf("Expected the input to take the key, value %q, OnKey calls %d", input.Value, called)
	}

	cursor.Right()
	c.Select(cursor, keyboard.EnterKey)
	if called != 0 {
		t.Error("Expected the button to take ENTER")
	}
}

// A board quits or moves on by returning from OnKey.
func TestContainer_OnKeyCanExit(t *testing.T) {
	g, _ := stgui.NewGrid([][]any{{"a"}}, nil)
	c := widgets.NewContainer(g)
	cursor := stgui.NewCursor(c, g, 0, 0, nil)

	next := stgui.NewScreen(nil, nil)
	c.OnKey = func(cursor *stgui.Cursor, input string) (*stgui.Screen, bool) {
		return next, true
	}

	if screen, exit := c.Select(cursor, "x"); screen != next || !exit {
		t.Errorf("Expected OnKey's screen to be returned, got %v exit=%v", screen, exit)
	}
}
