package stgui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jamieyoung5/gostrc"
	"github.com/jamieyoung5/stgui/term"
)

type Component interface {
	GetDimensions() (height int, width int)
	Print(cursor *Cursor)
	Serialize(cursor *Cursor) string
	Select(cursor *Cursor, macro string) (screen *Screen, exit bool)
}

type ViewComponent struct {
	Cursor    *Cursor
	Component Component
}

func (vc *ViewComponent) Serialize() string {
	return vc.Component.Serialize(vc.Cursor)
}

type Screen struct {
	Views   [][]*ViewComponent
	Cursors []*Cursor
	Persist bool
}

func (s *Screen) Serialize() string {
	var builder strings.Builder

	for _, row := range s.Views {
		var items []string
		for _, componentNode := range row {
			items = append(items, componentNode.Serialize())
		}
		// Draw each row of Components side by side, with a specified number of spaces in between
		builder.WriteString(sideBySide(items, 4))
		builder.WriteString("\n\n") // Add spacing between rows
	}

	return builder.String()
}

type Canvas struct {
	Screens *gostrc.Stack[*Screen]
}

func NewCanvas(screen *Screen) *Canvas {
	screenStack := gostrc.NewStack[*Screen]()
	screenStack.Push(screen)

	return &Canvas{screenStack}
}

func (ic *Canvas) Draw() {
	term.Clear()
	screen := ic.Screens.Peek()

	serializedScreen := screen.Serialize()
	fmt.Println(serializedScreen)
}

func (ic *Canvas) Initiate() {
	quit := make(chan bool)
	go render(ic, quit)

	for !ic.Screens.IsEmpty() {
		ic.Draw()

		reader := bufio.NewReader(os.Stdin)

		ic.listenForInput(reader)
	}

	quit <- true
}

func (ic *Canvas) listenForInput(reader *bufio.Reader) {
	for {
		sequence, err := term.ReadKeySequence(reader)
		if err != nil {
			fmt.Printf("error reading from input (%d)", err)
			continue
		}

		screen := ic.Screens.Peek()
		for _, cursor := range screen.Cursors {
			if cursor == nil {
				continue
			}

			if macro, ok := cursor.controls[term.Encode(sequence)]; ok {
				nextScreen, exit := screen.Views[cursor.gridY][cursor.gridX].Component.Select(cursor, macro)
				ic.Draw()

				if exit {
					if nextScreen != nil {
						if !screen.Persist {
							ic.Screens.Pop()
						}
					}

					ic.Screens.Push(nextScreen)
				} else {
					ic.Screens.Pop()
				}

				return
			}
		}
		ic.Draw()
	}
}
