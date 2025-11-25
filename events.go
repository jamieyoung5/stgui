package stgui

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

func listenForInput(eventChan chan Event, r *bufio.Reader) {
	sequence, err := keyboard.ReadInput(r, os.Stdin)
	if err != nil {
		eventChan <- ErrorEvent{
			Err: err,
		}
	}

	eventChan <- KeyPressEvent{
		Input: sequence,
	}
}

func listen(eventChan chan Event, r *bufio.Reader) {
	go listenForInput(eventChan, r)
}
