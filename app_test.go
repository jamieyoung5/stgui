package stgui_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/testutil"
	"golang.org/x/term"
)

func TestNewApp(t *testing.T) {
	screen := stgui.NewScreen(nil, nil)
	app := stgui.NewApp(screen)

	if app == nil {
		t.Fatal("NewApp returned nil")
	}
	if app.Screens == nil {
		t.Error("App.Screens stack not initialized")
	}
	if app.Screens.Peek() != screen {
		t.Error("App did not push initial screen")
	}
}

func TestApp_Display(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	mw := &testutil.MockWidget{RenderOutput: "WidgetOutput"}
	screen := stgui.NewScreen(nil, [][]stgui.Widget{{mw}})
	app := stgui.NewApp(screen)

	app.Display()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "WidgetOutput") {
		t.Errorf("Output missing widget content. Got: %q", output)
	}
	if !strings.Contains(output, "\033[?25l") {
		t.Error("Output missing cursor hide sequence")
	}
}

// Only what changed gets redrawn, and what the last frame left behind is wiped.
func TestApp_Display_RedrawsChangedLinesOnly(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	mw := &testutil.MockWidget{RenderOutput: "first\nsecond"}
	app := stgui.NewApp(stgui.NewScreen(nil, [][]stgui.Widget{{mw}}))
	app.Display()

	mw.RenderOutput = "first"
	app.Display()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if strings.Count(output, "first") != 1 {
		t.Errorf("Unchanged line was redrawn: %q", output)
	}
	// Second line is gone, so its row gets cleared.
	if !strings.Contains(output, "\033[2;1H\033[K") {
		t.Errorf("Removed line was not erased: %q", output)
	}
}

// No terminal, nothing to put into raw mode. This used to panic.
func TestApp_Run_WithoutTerminal(t *testing.T) {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		t.Skip("stdin is a terminal, Run would take over the session")
	}

	err := stgui.NewApp(stgui.NewScreen(nil, nil)).Run()
	if err == nil {
		t.Error("Expected an error when stdin is not a terminal")
	}
}
