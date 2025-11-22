package stgui

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/jamieyoung5/gostrc"
	"github.com/jamieyoung5/stgui/keyboard"
	"github.com/jamieyoung5/stgui/term"
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
	term.Clear()
	screen := a.Screens.Peek()

	serializedScreen := screen.Render()
	fmt.Println(serializedScreen)
}

func (a *App) Run() {
	quit := make(chan bool)
	go renderWorker(a, quit)

	for !a.Screens.IsEmpty() {
		a.Display()

		reader := bufio.NewReader(os.Stdin)

		a.listen(reader)
	}

	quit <- true
}

func (a *App) listen(reader *bufio.Reader) {
	for {
		sequence, err := keyboard.ReadInput(reader)
		if err != nil {
			fmt.Printf("error reading from input (%d)", err)
			continue
		}

		screen := a.Screens.Peek()
		for _, cursor := range screen.Cursors {
			if cursor == nil {
				continue
			}

			if input, ok := cursor.controls[sequence]; ok {
				a.handleSelection(input, cursor, screen)
			}

			return
		}
		a.Display()
	}
}

func (a *App) handleSelection(input string, cursor *Cursor, screen *Screen) {
	nextScreen, exit := screen.SelectElement(cursor, input)
	a.Display()

	if exit {
		if nextScreen != nil && !screen.Persist {
			a.Screens.Pop()
		}

		a.Screens.Push(nextScreen)
	} else {
		a.Screens.Pop()
	}
}

func renderWorker(app *App, quit chan bool) {
	screen := app.Screens.Peek().Render()
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Millisecond * 500)
			render(app, screen)
		}
	}
}

func render(app *App, currentScreen string) {
	newScreen := app.Screens.Peek().Render()
	if newScreen != currentScreen {
		term.Clear()
		app.Display()
	}
}
