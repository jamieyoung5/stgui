package stgui

import (
	"io"
	"os"
	"strings"
	"testing"
)

// Not testutil: that imports stgui, which an internal test can't do.
type stubWidget struct {
	next *Screen
	exit bool
	last string
}

func (s *stubWidget) Size() (int, int) { return 0, 0 }

func (s *stubWidget) Render() string { return "stub" }

func (s *stubWidget) Select(cursor *Cursor, input string) (*Screen, bool) {
	s.last = input
	return s.next, s.exit
}

// Swap stdout for a pipe we drain and throw away, or drawing scribbles all over
// the test output.
func discardStdout(t *testing.T) {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	done := make(chan struct{})
	go func() {
		io.Copy(io.Discard, r)
		close(done)
	}()

	t.Cleanup(func() {
		os.Stdout = old
		w.Close()
		<-done
		r.Close()
	})
}

func screenWith(widget Widget) *Screen {
	return NewScreen([]*Cursor{NewCursor(widget, nil, 0, 0, nil)}, [][]Widget{{widget}})
}

func TestApp_HandleSelection_NavigatesToScreen(t *testing.T) {
	discardStdout(t)

	next := NewScreen(nil, nil)
	widget := &stubWidget{next: next, exit: true}
	current := screenWith(widget)
	app := NewApp(current)

	app.handleInput("ENTER")

	if app.Screens.Size() != 1 {
		t.Errorf("Expected the current screen to be replaced, stack size is %d", app.Screens.Size())
	}
	if app.Screens.Peek() != next {
		t.Error("Expected the new screen to be on top of the stack")
	}
}

func TestApp_HandleSelection_PersistKeepsScreen(t *testing.T) {
	discardStdout(t)

	next := NewScreen(nil, nil)
	widget := &stubWidget{next: next, exit: true}
	current := screenWith(widget)
	current.Persist = true
	app := NewApp(current)

	app.handleInput("ENTER")

	if app.Screens.Size() != 2 {
		t.Errorf("Expected the persisted screen to stay on the stack, size is %d", app.Screens.Size())
	}
	if app.Screens.Peek() != next {
		t.Error("Expected the new screen to be on top of the stack")
	}
}

func TestApp_HandleSelection_ExitPopsScreen(t *testing.T) {
	discardStdout(t)

	app := NewApp(screenWith(&stubWidget{exit: true}))

	app.handleInput("ENTER")

	if !app.Screens.IsEmpty() {
		t.Error("Expected exiting the last screen to empty the stack, which ends Run")
	}
}

func TestApp_HandleInput_MapsControls(t *testing.T) {
	discardStdout(t)

	widget := &stubWidget{}
	screen := screenWith(widget)
	screen.Cursors[0].Controls = map[string]string{"w": ControlUp}
	app := NewApp(screen)

	app.handleInput("w")
	if widget.last != ControlUp {
		t.Errorf("Expected 'w' to map to %q, widget got %q", ControlUp, widget.last)
	}

	app.handleInput("q")
	if widget.last != "q" {
		t.Errorf("Expected unmapped input to pass through, widget got %q", widget.last)
	}
}

func TestApp_IsQuitKey(t *testing.T) {
	app := NewApp(NewScreen(nil, nil))

	if !app.isQuitKey("CTRL+C") {
		t.Error("Expected Ctrl+C to quit by default")
	}
	if app.isQuitKey("a") {
		t.Error("Expected an ordinary key not to quit")
	}

	app.QuitKeys = nil
	if app.isQuitKey("CTRL+C") {
		t.Error("Expected no quit keys once QuitKeys is cleared")
	}
}

func TestApp_Fit(t *testing.T) {
	app := &App{width: 4, height: 2}

	got := app.fit([]string{"abcdefg", "ab", "cut"})

	if len(got) != 2 {
		t.Fatalf("Expected the frame to be clipped to 2 lines, got %d: %q", len(got), got)
	}
	if got[0] != "abcd" {
		t.Errorf("Expected the first line clipped to 4 columns, got %q", got[0])
	}
	if got[1] != "ab" {
		t.Errorf("Expected a short line to be left alone, got %q", got[1])
	}
}

func TestTruncateVisible(t *testing.T) {
	if got := truncateVisible("abc", 10); got != "abc" {
		t.Errorf("Expected a short line unchanged, got %q", got)
	}

	// Escapes cost no columns, so it's the text that gets cut. The style is
	// closed off so it doesn't bleed into the next line.
	styled := selectedStyle + "abcdef" + resetStyle
	got := truncateVisible(styled, 3)
	if !strings.HasPrefix(got, selectedStyle) {
		t.Errorf("Expected the style to be kept, got %q", got)
	}
	if !strings.HasSuffix(got, resetStyle) {
		t.Errorf("Expected the style to be closed, got %q", got)
	}
	if !strings.Contains(got, "abc") || strings.Contains(got, "abcd") {
		t.Errorf("Expected 3 visible characters, got %q", got)
	}
}

func TestApp_NoScreen(t *testing.T) {
	discardStdout(t)

	app := NewApp(nil)
	if !app.Screens.IsEmpty() {
		t.Error("Expected no screens when constructed with nil")
	}

	// Empty stack: mustn't panic.
	app.Display()
}
