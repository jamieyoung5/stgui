package stgui

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jamieyoung5/stgui/events"
)

// Not testutil: that imports stgui, which an internal test can't do.
type stubWidget struct {
	action Action
	last   Input
}

func (s *stubWidget) Render() string { return "stub" }

func (s *stubWidget) Select(cursor *Cursor, input Input) Action {
	s.last = input
	return s.action
}

// A custom event, the kind you'd Send yourself.
type stubEvent struct {
	events.Msg
	value string
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
	return NewScreen(NewCursor(widget, nil, 0, 0), Row(widget))
}

func TestApp_Push_KeepsCurrentScreen(t *testing.T) {
	discardStdout(t)

	next := NewScreen(nil)
	widget := &stubWidget{action: Push(next)}
	app := NewApp(screenWith(widget))

	app.handleInput(KeyEnter)

	if app.Screens.Size() != 2 {
		t.Errorf("Expected Push to stack a screen on top, size is %d", app.Screens.Size())
	}
	if app.Screens.Peek() != next {
		t.Error("Expected the new screen to be on top of the stack")
	}
}

func TestApp_Replace_SwapsScreen(t *testing.T) {
	discardStdout(t)

	next := NewScreen(nil)
	widget := &stubWidget{action: Replace(next)}
	app := NewApp(screenWith(widget))

	app.handleInput(KeyEnter)

	if app.Screens.Size() != 1 {
		t.Errorf("Expected Replace to swap the screen out, size is %d", app.Screens.Size())
	}
	if app.Screens.Peek() != next {
		t.Error("Expected the new screen to be on top of the stack")
	}
}

// Push then Close comes back to where you were.
func TestApp_Close_ReturnsToPushedScreen(t *testing.T) {
	discardStdout(t)

	closer := &stubWidget{action: Close()}
	second := screenWith(closer)
	first := screenWith(&stubWidget{action: Push(second)})
	app := NewApp(first)

	app.handleInput(KeyEnter)
	app.handleInput(KeyEnter)

	if app.Screens.Size() != 1 {
		t.Fatalf("Expected to be back to one screen, size is %d", app.Screens.Size())
	}
	if app.Screens.Peek() != first {
		t.Error("Expected Close to land back on the screen underneath")
	}
}

func TestApp_Close_EmptiesStackOnLastScreen(t *testing.T) {
	discardStdout(t)

	app := NewApp(screenWith(&stubWidget{action: Close()}))

	app.handleInput(KeyEnter)

	if !app.Screens.IsEmpty() {
		t.Error("Expected closing the last screen to empty the stack, which ends Run")
	}
}

func TestApp_Quit_StopsWithScreensLeft(t *testing.T) {
	discardStdout(t)

	quitter := &stubWidget{action: Quit()}
	second := screenWith(quitter)
	app := NewApp(screenWith(&stubWidget{action: Push(second)}))

	app.handleInput(KeyEnter)
	app.handleInput(KeyEnter)

	if !app.quitting {
		t.Error("Expected Quit to stop the loop")
	}
	if app.Screens.IsEmpty() {
		t.Error("Expected Quit to leave the stack alone, not unwind it")
	}
}

// A bound key reaches the widget as a control, an unbound one as a plain key.
func TestApp_HandleInput_AppliesBindings(t *testing.T) {
	discardStdout(t)

	widget := &stubWidget{}
	app := NewApp(screenWith(widget)).Bind("w", ControlUp)

	app.handleInput("w")
	if !widget.last.Is(ControlUp) {
		t.Errorf("Expected 'w' bound to %q, widget got %+v", ControlUp, widget.last)
	}
	if widget.last.Key != "w" {
		t.Errorf("Expected the original key kept alongside the control, got %q", widget.last.Key)
	}

	app.handleInput("q")
	if widget.last.Control != "" {
		t.Errorf("Expected an unbound key to carry no control, got %q", widget.last.Control)
	}
	if widget.last.Key != "q" {
		t.Errorf("Expected an unbound key to pass through, got %q", widget.last.Key)
	}
	if r, ok := widget.last.Rune(); !ok || r != 'q' {
		t.Errorf("Expected an unbound key to be typeable, got %q ok=%v", r, ok)
	}
}

// A bound key must not also be typed into a text field.
func TestInput_BoundKeyIsNotText(t *testing.T) {
	bound := Input{Key: "w", Control: ControlUp}
	if _, ok := bound.Rune(); ok {
		t.Error("Expected a bound key not to count as typed text")
	}

	if _, ok := (Input{Key: KeyEnter}).Rune(); ok {
		t.Error("Expected a named key not to count as typed text")
	}
}

func TestApp_OnEvent_HandlesCustomEvents(t *testing.T) {
	discardStdout(t)

	app := NewApp(screenWith(&stubWidget{}))

	var got string
	app.OnEvent(func(event events.Event) Action {
		if e, ok := event.(stubEvent); ok {
			got = e.value
		}
		return Stay()
	})

	app.handleEvent(stubEvent{value: "tick"})

	if got != "tick" {
		t.Errorf("Expected OnEvent to see the custom event, got %q", got)
	}
}

// An event nobody handles must be ignored, not fatal.
func TestApp_OnEvent_UnsetIgnoresEvent(t *testing.T) {
	discardStdout(t)

	app := NewApp(screenWith(&stubWidget{}))
	app.handleEvent(stubEvent{value: "ignored"})

	if app.Screens.Size() != 1 {
		t.Error("Expected an unhandled event to leave the stack alone")
	}
}

// Send must not block when the app never runs, or a sending goroutine leaks.
func TestApp_Send_DoesNotBlock(t *testing.T) {
	app := NewApp(screenWith(&stubWidget{}))

	done := make(chan struct{})
	go func() {
		app.Send(stubEvent{value: "a"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Send blocked with room in the buffer")
	}

	if got := <-app.eventChan; got.(stubEvent).value != "a" {
		t.Error("Expected the sent event on the channel")
	}
}

// Once Run has returned, Send gives up rather than blocking forever.
func TestApp_Send_GivesUpAfterQuit(t *testing.T) {
	app := NewApp(screenWith(&stubWidget{}))
	close(app.done)

	// Fill the buffer so the send has nowhere to go.
	for range eventBuffer {
		app.eventChan <- stubEvent{}
	}

	done := make(chan struct{})
	go func() {
		app.Send(stubEvent{value: "dropped"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Send blocked after the app finished")
	}
}

func TestApp_IsQuitKey(t *testing.T) {
	app := NewApp(NewScreen(nil))

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
