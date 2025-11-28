package stgui

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jamieyoung5/gostrc"
	"github.com/jamieyoung5/stgui/events"
	"golang.org/x/term"
)

type App struct {
	Screens    *gostrc.Stack[*Screen]
	lastBuffer []string
}

func NewApp(screen *Screen) *App {
	screenStack := gostrc.NewStack[*Screen]()
	screenStack.Push(screen)

	return &App{Screens: screenStack}
}

func (a *App) Display() {
	screen := a.Screens.Peek()
	output := screen.Render()
	lines := strings.Split(output, "\n")

	if a.lastBuffer == nil {
		fmt.Print("\033[H\033[J")
		fmt.Print("\033[?25l")
	}

	bw := bufio.NewWriter(os.Stdout)
	for i, line := range lines {
		if a.lastBuffer != nil && i < len(a.lastBuffer) && a.lastBuffer[i] == line {
			continue
		}

		fmt.Fprintf(bw, "\033[%d;1H\033[K%s", i+1, line)
	}

	if a.lastBuffer != nil && len(lines) < len(a.lastBuffer) {
		for i := len(lines); i < len(a.lastBuffer); i++ {
			fmt.Printf("\033[%d;1H\033[K", i+1)
		}
	}

	bw.Flush()
	a.lastBuffer = lines
}

func (a *App) Run() error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		fmt.Print("\033[?25h")
		term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	reader := bufio.NewReader(os.Stdin)
	eventChan := make(chan events.Event)

	events.Listen(eventChan, reader)

	a.Display()
	for !a.Screens.IsEmpty() {
		event := <-eventChan

		switch v := event.(type) {
		case events.KeyPressEvent:
			a.handleInput(v.Input)
		case events.ResizeEvent:
			a.handleResize()
		case events.ErrorEvent:
			return v.Err
		default:
			return errors.New("recieved unknown event")
		}
	}

	return nil
}

func (a *App) handleResize() {
	a.lastBuffer = nil
	a.Display()
}

func (a *App) handleInput(sequence string) {
	screen := a.Screens.Peek()
	for _, cursor := range screen.Cursors {
		if cursor == nil {
			continue
		}

		if input, ok := cursor.Controls[sequence]; ok {
			a.handleSelection(input, cursor, screen)
			break
		}

		a.handleSelection(sequence, cursor, screen)
		break
	}
}

func (a *App) handleSelection(input string, cursor *Cursor, screen *Screen) {
	nextScreen, exit := cursor.Select(input)

	if exit {
		if nextScreen != nil && !screen.Persist {
			a.Screens.Pop()
		}
		if nextScreen != nil {
			a.Screens.Push(nextScreen)
		} else {
			a.Screens.Pop()
		}

		if !a.Screens.IsEmpty() {
			a.Display()
		}
	} else {
		a.Display()
	}
}

func clearTerm() {
	fmt.Print("\033[H\033[J")
}
