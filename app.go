package stgui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jamieyoung5/gostrc"
	"github.com/jamieyoung5/stgui/keyboard"
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

	a.Display()
	for !a.Screens.IsEmpty() {
		if !a.listen(reader) {
			break
		}
	}
}

func (a *App) listen(reader *bufio.Reader) bool {
	sequence, err := keyboard.ReadInput(reader, os.Stdin)
	if err != nil {
		return false
	}

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

	return true
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
