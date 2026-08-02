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

// The constructor selects the starting cell, so nobody has to reach into
// Grid.Cells to highlight it by hand.
func TestCursor_SelectsStartingCell(t *testing.T) {
	g, _ := stgui.NewGrid([][]any{{"1", "2"}, {"3", "4"}}, nil)

	stgui.NewCursor(nil, g, 1, 0, nil)

	if !g.Cells[1][0].Selected {
		t.Error("Starting cell not selected")
	}
	if g.Cells[0][0].Selected {
		t.Error("Cell other than the starting cell is selected")
	}
}

func TestCursor_ClampsStartingPosition(t *testing.T) {
	g, _ := stgui.NewGrid([][]any{{"1", "2"}}, nil)

	c := stgui.NewCursor(nil, g, 9, -3, nil)
	if c.Row != 0 || c.Col != 0 {
		t.Errorf("Expected out-of-range start to clamp to 0,0, got %d,%d", c.Row, c.Col)
	}
	if !g.Cells[0][0].Selected {
		t.Error("Clamped cell not selected")
	}
}

func TestCursor_DefaultsToDirectionalControls(t *testing.T) {
	c := stgui.NewCursor(nil, nil, 0, 0, nil)

	if c.Controls["UP"] != stgui.ControlUp {
		t.Errorf("Expected default controls, got %v", c.Controls)
	}
}

// No grid, but the app will still call these on whatever the screen holds.
func TestCursor_NoGrid(t *testing.T) {
	c := stgui.NewCursor(nil, nil, 0, 0, nil)

	c.Up()
	c.Down()
	c.Left()
	c.Right()

	if cell := c.Cell(); cell != nil {
		t.Errorf("Expected no cell without a grid, got %v", cell)
	}
	if screen, exit := c.Select("ENTER"); screen != nil || exit {
		t.Error("Expected Select on a widgetless cursor to do nothing")
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
