package stgui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/jamieyoung5/gostrc"
	"github.com/jamieyoung5/stgui/events"
	"golang.org/x/term"
)

// EventHandler takes an event the app doesn't recognise, one you sent yourself.
type EventHandler func(event events.Event) Action

const (
	clearScreen = "\033[H\033[J"
	hideCursor  = "\033[?25l"
	showCursor  = "\033[?25h"

	// Room for a burst from a goroutine before Send starts waiting.
	eventBuffer = 64
)

// App owns the terminal and runs the event loop over a stack of screens. Only the
// top one is drawn and fed input. Run returns once the stack empties.
type App struct {
	Screens *gostrc.Stack[*Screen]

	// QuitKeys get you out whatever the focused widget would have done with
	// them. Ctrl+C by default; nil forces an in-app quit.
	QuitKeys []Key

	// Bindings map keys onto controls for every screen. Arrow keys and Enter
	// by default.
	Bindings Bindings

	onEvent   EventHandler
	eventChan chan events.Event
	done      chan struct{}
	quitting  bool

	lastBuffer    []string
	width, height int
}

func NewApp(screen *Screen) *App {
	screenStack := gostrc.NewStack[*Screen]()
	if screen != nil {
		screenStack.Push(screen)
	}

	return &App{
		Screens:   screenStack,
		QuitKeys:  []Key{KeyCtrlC},
		Bindings:  DefaultBindings(),
		eventChan: make(chan events.Event, eventBuffer),
		done:      make(chan struct{}),
	}
}

// OnEvent sets the handler for events the app doesn't know about, the ones you
// Send yourself. The screen is redrawn afterwards, so writing to a label is
// enough to animate it.
func (a *App) OnEvent(handler EventHandler) *App {
	a.onEvent = handler
	return a
}

// Bind points a key at a control, for every screen. Binding a key takes it
// away from typing, so don't bind letters on a screen with an Input on it.
func (a *App) Bind(key Key, control Control) *App {
	if a.Bindings == nil {
		a.Bindings = Bindings{}
	}
	a.Bindings[key] = control

	return a
}

// QuitOn adds keys that get you out of the app, on top of the default Ctrl+C.
func (a *App) QuitOn(keys ...Key) *App {
	a.QuitKeys = append(a.QuitKeys, keys...)
	return a
}

// Send hands an event to the running app, from any goroutine. It returns once
// the app has the event, or straight away if the app has already quit, so a
// sending goroutine never leaks.
func (a *App) Send(event events.Event) {
	select {
	case a.eventChan <- event:
	case <-a.done:
	}
}

// Display draws the top screen, writing out only the lines that changed since the
// last frame.
func (a *App) Display() {
	screen := a.Screens.Peek()
	if screen == nil {
		return
	}

	lines := a.fit(strings.Split(screen.Render(), "\n"))

	bw := bufio.NewWriter(os.Stdout)
	if a.lastBuffer == nil {
		bw.WriteString(clearScreen)
		bw.WriteString(hideCursor)
	}

	for i, line := range lines {
		if i < len(a.lastBuffer) && a.lastBuffer[i] == line {
			continue
		}

		fmt.Fprintf(bw, "\033[%d;1H\033[K%s", i+1, line)
	}

	// Wipe anything the last frame left behind.
	for i := len(lines); i < len(a.lastBuffer); i++ {
		fmt.Fprintf(bw, "\033[%d;1H\033[K", i+1)
	}

	bw.Flush()
	a.lastBuffer = lines
}

// Run takes over the terminal and handles events until the screen stack empties,
// a quit key comes in, or something goes wrong. The terminal is put back either
// way.
func (a *App) Run() error {
	fd := int(os.Stdin.Fd())

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("stgui: could not put the terminal into raw mode: %w", err)
	}
	defer a.restore(fd, oldState)

	// Lets anything blocked in Send go, once we stop reading.
	defer close(a.done)

	a.readTerminalSize()
	events.Listen(a.eventChan, bufio.NewReader(os.Stdin))

	a.Display()
	for !a.Screens.IsEmpty() && !a.quitting {
		switch event := (<-a.eventChan).(type) {
		case events.KeyPressEvent:
			if a.isQuitKey(event.Key) {
				return nil
			}
			a.handleInput(event.Key)
		case events.ResizeEvent:
			a.handleResize(event)
		case events.ErrorEvent:
			// No stdin, no app. Piped or closed sessions end here, which is
			// not a failure.
			if errors.Is(event.Err, io.EOF) {
				return nil
			}
			return event.Err
		default:
			a.handleEvent(event)
		}
	}

	return nil
}

// restore hands the terminal back: cursor visible, raw mode off, and parked below
// the last frame so the shell prompt doesn't land on top of it.
func (a *App) restore(fd int, oldState *term.State) {
	fmt.Printf("\033[%d;1H%s\n", len(a.lastBuffer)+1, showCursor)
	term.Restore(fd, oldState)
}

func (a *App) isQuitKey(key Key) bool {
	return slices.Contains(a.QuitKeys, key)
}

// fit clips the frame to the terminal. A line that wraps would shove everything
// below it down a row and throw the whole diff out.
func (a *App) fit(lines []string) []string {
	if a.height > 0 && len(lines) > a.height {
		lines = lines[:a.height]
	}

	if a.width <= 0 {
		return lines
	}

	fitted := make([]string, len(lines))
	for i, line := range lines {
		fitted[i] = truncateVisible(line, a.width)
	}

	return fitted
}

func (a *App) handleResize(event events.ResizeEvent) {
	a.width, a.height = event.Width, event.Height
	// The terminal threw its contents away, so the diff is worthless. Redraw
	// the lot.
	a.lastBuffer = nil
	a.Display()
}

func (a *App) readTerminalSize() {
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		a.width, a.height = w, h
	}
}

// handleEvent passes an event we don't recognise to OnEvent. An app without
// one ignores it.
func (a *App) handleEvent(event events.Event) {
	if a.onEvent == nil {
		return
	}

	a.apply(a.onEvent(event))
}

func (a *App) handleInput(key Key) {
	screen := a.Screens.Peek()
	if screen == nil || screen.Cursor == nil {
		return
	}

	a.apply(screen.Cursor.Select(Input{Key: key, Control: a.Bindings[key]}))
}

// apply carries out what a widget asked for and redraws.
func (a *App) apply(action Action) {
	switch action.kind {
	case actionPush:
		a.Screens.Push(action.screen)
	case actionReplace:
		a.Screens.Pop()
		a.Screens.Push(action.screen)
	case actionClose:
		a.Screens.Pop()
	case actionQuit:
		a.quitting = true
	}

	if action.kind != actionStay {
		// Nothing on the new screen lines up with the old one. Start over.
		a.lastBuffer = nil
	}

	// On the way out, leave the frame alone for restore to park the cursor
	// after.
	if a.quitting || a.Screens.IsEmpty() {
		return
	}

	a.Display()
}
