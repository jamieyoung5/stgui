package stgui

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jamieyoung5/gostrc"
	"golang.org/x/term"
)

type App struct {
	Screens *gostrc.Stack[*Screen]
}

func NewApp(screen *Screen) *App {
	screenStack := gostrc.NewStack[*Screen]()
	screenStack.Push(screen)

	return &App{screenStack}
}

func (a *App) Display() {
	clearTerm()
	screen := a.Screens.Peek()

	serializedScreen := screen.Render()
	fmt.Print(strings.ReplaceAll(serializedScreen, "\n", "\r\n"))
}

func (a *App) Run() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	reader := bufio.NewReader(os.Stdin)
	eventChan := make(chan Event)

	a.Display()
	for !a.Screens.IsEmpty() {
		listen(eventChan, reader)
		event := <-eventChan

		switch v := event.(type) {
		case KeyPressEvent:
			a.handleInput(v.Input)
		case ErrorEvent:
			panic(v.Err)
		default:
			panic(errors.New("recieved unknown event"))
		}
	}

}

func (a *App) handleInput(sequence string) {
	screen := a.Screens.Peek()
	for _, cursor := range screen.Cursors {
		if cursor == nil {
			continue
		}

		if input, ok := cursor.controls[sequence]; ok {
			a.handleSelection(input, cursor, screen)
			break
		}
	}
}

func (a *App) handleSelection(input string, cursor *Cursor, screen *Screen) {
	nextScreen, exit := cursor.Select(input)
	a.Display()

	if exit {
		if nextScreen != nil && !screen.Persist {
			a.Screens.Pop()
		}
		if nextScreen != nil {
			a.Screens.Push(nextScreen)
		} else {
			a.Screens.Pop()
		}
	}
}

func clearTerm() {
	fmt.Print("\033[H\033[J")
}
