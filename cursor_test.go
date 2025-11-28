package stgui_test

import (
	"testing"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/testutil"
)

func TestCursor_Movement(t *testing.T) {
	data := [][]any{
		{"1", "2"},
		{"3", "4"},
	}
	g, _ := stgui.NewGrid(data, nil)

	c := stgui.NewCursor(nil, g, 0, 0, nil)
	g.Cells[0][0].Selected = true

	c.Right()
	if c.Row != 0 || c.Col != 1 {
		t.Errorf("Expected 0,1 after Right, got %d,%d", c.Row, c.Col)
	}
	if g.Cells[0][0].Selected {
		t.Error("Previous cell (0,0) still selected")
	}
	if !g.Cells[0][1].Selected {
		t.Error("New cell (0,1) not selected")
	}

	c.Right()
	if c.Col != 1 {
		t.Error("Cursor moved past right boundary")
	}

	c.Down()
	if c.Row != 1 || c.Col != 1 {
		t.Errorf("Expected 1,1 after Down, got %d,%d", c.Row, c.Col)
	}

	c.Down()
	if c.Row != 1 {
		t.Error("Cursor moved past bottom boundary")
	}

	c.Left()
	if c.Row != 1 || c.Col != 0 {
		t.Errorf("Expected 1,0 after Left, got %d,%d", c.Row, c.Col)
	}

	c.Up()
	if c.Row != 0 || c.Col != 0 {
		t.Errorf("Expected 0,0 after Up, got %d,%d", c.Row, c.Col)
	}
}

func TestCursor_Select(t *testing.T) {
	mw := &testutil.MockWidget{}
	c := stgui.NewCursor(mw, nil, 0, 0, nil)

	c.Select("ENTER")
	if mw.LastInput != "ENTER" {
		t.Errorf("Select input not passed to widget")
	}
}
