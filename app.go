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
	"github.com/jamieyoung5/stgui/keyboard"
	"golang.org/x/term"
)

const (
	clearScreen = "\033[H\033[J"
	hideCursor  = "\033[?25l"
	showCursor  = "\033[?25h"
)

// App owns the terminal and runs the event loop over a stack of screens. Only the
// top one is drawn and fed input. Run returns once the stack empties.
type App struct {
	Screens *gostrc.Stack[*Screen]

	// QuitKeys get you out whatever the focused widget would have done with
	// them. Ctrl+C by default; nil forces an in-app quit action.
	QuitKeys []string

	lastBuffer    []string
	width, height int
}

func NewApp(screen *Screen) *App {
	screenStack := gostrc.NewStack[*Screen]()
	if screen != nil {
		screenStack.Push(screen)
	}

	return &App{
		Screens:  screenStack,
		QuitKeys: []string{keyboard.CtrlCKey},
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

	a.readTerminalSize()

	eventChan := make(chan events.Event)
	events.Listen(eventChan, bufio.NewReader(os.Stdin))

	a.Display()
	for !a.Screens.IsEmpty() {
		switch event := (<-eventChan).(type) {
		case events.KeyPressEvent:
			if a.isQuitKey(event.Input) {
				return nil
			}
			a.handleInput(event.Input)
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
			return fmt.Errorf("stgui: received unknown event %T", event)
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

func (a *App) isQuitKey(input string) bool {
	return slices.Contains(a.QuitKeys, input)
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

func (a *App) handleInput(sequence string) {
	screen := a.Screens.Peek()
	cursor := screen.ActiveCursor()
	if cursor == nil {
		return
	}

	input := sequence
	if mapped, ok := cursor.Controls[sequence]; ok {
		input = mapped
	}

	a.handleSelection(input, cursor, screen)
}

func (a *App) handleSelection(input string, cursor *Cursor, screen *Screen) {
	nextScreen, exit := cursor.Select(input)

	if exit {
		if nextScreen == nil || !screen.Persist {
			a.Screens.Pop()
		}
		if nextScreen != nil {
			a.Screens.Push(nextScreen)
		}

		// That was the last screen, so we're on the way out. Leave the frame
		// alone for restore to park the cursor after.
		if a.Screens.IsEmpty() {
			return
		}

		// Nothing on the new screen lines up with the old one. Start over.
		a.lastBuffer = nil
	}

	a.Display()
}
