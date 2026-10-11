package events

import (
	"bufio"
	"os"

	"github.com/jamieyoung5/stgui/keyboard"
)

// Event is anything the app loop can be handed. Embed Msg to make your own.
type Event interface {
	IsMsg()
}

// Msg is the marker to embed in your own events.
type Msg struct{}

func (m Msg) IsMsg() {}

type ErrorEvent struct {
	Msg
	Err error
}

type KeyPressEvent struct {
	Msg
	Key keyboard.Key
}

type ResizeEvent struct {
	Msg
	Width  int
	Height int
}

func listenForInput(eventChan chan Event, r *bufio.Reader) {
	for {
		key, err := keyboard.ReadInput(r, os.Stdin)
		if err != nil {
			// Stdin is gone. Say so and stop, or we spin forever on a reader
			// that will never come back.
			eventChan <- ErrorEvent{Err: err}
			return
		}

		eventChan <- KeyPressEvent{Key: key}
	}
}

// Listen starts the keyboard and resize listeners. Returns straight away, both
// run in the background.
func Listen(eventChan chan Event, r *bufio.Reader) {
	go listenForInput(eventChan, r)
	go listenForResize(eventChan)
}
