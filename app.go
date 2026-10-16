package stgui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"

	"github.com/jamieyoung5/gostrc"
	"github.com/jamieyoung5/stgui/events"
	"golang.org/x/term"
)

// EventHandler takes an event the app doesn't recognise.
type EventHandler func(event events.Event) Action

const (
	clearScreen = "\033[H\033[J"
	hideCursor  = "\033[?25l"
	showCursor  = "\033[?25h"

	// Button reporting, with coordinates that survive a wide terminal.
	enableMouse  = "\033[?1000h\033[?1006h"
	disableMouse = "\033[?1006l\033[?1000l"

	// Room for a burst of events before Send blocks.
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

	mouse     bool
	onEvent   EventHandler
	eventChan chan events.Event
	done      chan struct{}
	quitting  bool

	out           *bufio.Writer
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

// WithMouse asks the terminal to report clicks, which select whatever was
// clicked on. Off by default: it is the terminal's own text selection that gets
// taken over, so it has to be worth it.
func (a *App) WithMouse() *App {
	a.mouse = true
	return a
}

// OnEvent sets the handler for events the app doesn't recognise. The screen is
// redrawn afterwards.
func (a *App) OnEvent(handler EventHandler) *App {
	a.onEvent = handler
	return a
}

// Bind points a key at a control, for every screen. A bound key can't be typed
// into an Input.
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
// the app has the event, or immediately if the app has quit.
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

	lines := a.fit(screen.renderLines())

	if a.out == nil {
		a.out = bufio.NewWriter(os.Stdout)
	}
	bw := a.out

	if a.lastBuffer == nil {
		bw.WriteString(clearScreen)
		bw.WriteString(hideCursor)
	}

	var esc [16]byte
	for i, line := range lines {
		if i < len(a.lastBuffer) && a.lastBuffer[i] == line {
			continue
		}

		bw.Write(appendMoveTo(esc[:0], i+1))
		bw.WriteString(line)
	}

	// Wipe anything the last frame left behind.
	for i := len(lines); i < len(a.lastBuffer); i++ {
		bw.Write(appendMoveTo(esc[:0], i+1))
	}

	bw.Flush()
	a.lastBuffer = lines
}

// appendMoveTo is the escape that parks the cursor at the start of row and clears
// the line. Built by hand: it runs for every changed line of every frame.
func appendMoveTo(dst []byte, row int) []byte {
	dst = append(dst, '\033', '[')
	dst = strconv.AppendInt(dst, int64(row), 10)

	return append(dst, ';', '1', 'H', '\033', '[', 'K')
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

	if a.mouse {
		fmt.Print(enableMouse)
	}

	// Unblocks anything waiting in Send.
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
		case events.MouseEvent:
			a.handleMouse(event)
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
	if a.mouse {
		fmt.Print(disableMouse)
	}

	fmt.Printf("\033[%d;1H%s\n", len(a.lastBuffer)+1, showCursor)
	term.Restore(fd, oldState)
}

func (a *App) isQuitKey(key Key) bool {
	return slices.Contains(a.QuitKeys, key)
}

// fit clips the frame to the terminal. A line that wraps would shove everything
// below it down a row and throw the whole diff out.
//
// Clipped in place: the frame is freshly rendered and nothing else holds it.
func (a *App) fit(lines []string) []string {
	if a.height > 0 && len(lines) > a.height {
		lines = lines[:a.height]
	}

	if a.width <= 0 {
		return lines
	}

	for i, line := range lines {
		if VisibleWidth(line) > a.width {
			lines[i] = truncateVisible(line, a.width)
		}
	}

	return lines
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

// handleMouse selects whatever was clicked on, as long as the cursor can reach
// it. A click on another widget is ignored: one screen has one cursor.
func (a *App) handleMouse(event events.MouseEvent) {
	if !event.Pressed {
		return
	}

	screen := a.Screens.Peek()
	if screen == nil || screen.Cursor == nil {
		return
	}

	widget, row, col, ok := screen.Hit(event.Col, event.Row)
	if !ok || widget != screen.Cursor.widget {
		return
	}

	screen.Cursor.MoveTo(row, col)
	a.apply(screen.Cursor.Select(Input{Control: ControlSelect}))
	a.Display()
}

// handleEvent passes an unrecognised event to OnEvent, if a handler is set.
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

	// Leave the last frame alone for restore to park the cursor after.
	if a.quitting || a.Screens.IsEmpty() {
		return
	}

	a.Display()
}
