package events

import (
	"bufio"
	"os"

	"github.com/jamieyoung5/stgui/keyboard"
)

type Event interface {
	IsMsg()
}

type Msg struct{}

func (m Msg) IsMsg() {}

type ErrorEvent struct {
	Err error
	Msg
}

func (e ErrorEvent) IsMsg() {}

type KeyPressEvent struct {
	Input string
	Msg
}

type ResizeEvent struct {
	Width  int
	Height int
	Msg
}

func listenForInput(eventChan chan Event, r *bufio.Reader) {
	for {
		sequence, err := keyboard.ReadInput(r, os.Stdin)
		if err != nil {
			eventChan <- ErrorEvent{Err: err}
		}

		eventChan <- KeyPressEvent{Input: sequence}
	}
}

func Listen(eventChan chan Event, r *bufio.Reader) {
	go listenForInput(eventChan, r)
	go listenForResize(eventChan)
}
