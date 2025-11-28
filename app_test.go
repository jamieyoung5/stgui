package stgui_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/testutil"
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
