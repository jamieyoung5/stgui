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
	Screens *gostrc.Stack[*View]
}

func NewApp(screen *View) *App {
	screenStack := gostrc.NewStack[*View]()
	screenStack.Push(screen)

	return &App{screenStack}
}

func (ic *App) Display() {
	term.Clear()
	screen := ic.Screens.Peek()

	serializedScreen := screen.Render()
	fmt.Println(serializedScreen)
}

func (ic *App) Run() {
	quit := make(chan bool)
	go renderWorker(ic, quit)

	for !ic.Screens.IsEmpty() {
		ic.Display()

		reader := bufio.NewReader(os.Stdin)

		ic.listen(reader)
	}

	quit <- true
}

func (ic *App) listen(reader *bufio.Reader) {
	for {
		sequence, err := keyboard.ReadInput(reader)
		if err != nil {
			fmt.Printf("error reading from input (%d)", err)
			continue
		}

		screen := ic.Screens.Peek()
		for _, cursor := range screen.Cursors {
			if cursor == nil {
				continue
			}

			if macro, ok := cursor.controls[sequence]; ok {
				ic.handleSelection(macro, cursor, screen)
			}

			return
		}
		ic.Display()
	}
}

func (ic *App) handleSelection(macro string, cursor *Cursor, view *View) {
	nextScreen, exit := view.Components[cursor.gridY][cursor.gridX].Component.Select(cursor, macro)
	ic.Display()

	if exit {
		if nextScreen != nil && !view.Persist {
			ic.Screens.Pop()
		}

		ic.Screens.Push(nextScreen)
	} else {
		ic.Screens.Pop()
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
