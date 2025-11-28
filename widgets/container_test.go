package widgets_test

import (
	"strings"
	"testing"

	"github.com/jamieyoung5/stgui"
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
	c.Select(cursor, "DOWN")

	clicked := false
	btn := widgets.NewButton("Btn", func() { clicked = true })
	gClick, _ := stgui.NewGrid([][]any{{btn}}, nil)
	cClick := widgets.NewContainer(gClick)
	cursorClick := stgui.NewCursor(cClick, gClick, 0, 0, nil)

	cClick.Select(cursorClick, "ENTER")
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
