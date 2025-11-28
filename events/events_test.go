package events_test

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/jamieyoung5/stgui/events"
)

func TestEventImplementation(t *testing.T) {
	var _ events.Event = events.Msg{}
	var _ events.Event = events.ErrorEvent{}
	var _ events.Event = events.KeyPressEvent{}
	var _ events.Event = events.ResizeEvent{}
}

func TestMsgMethods(t *testing.T) {
	m := events.Msg{}
	m.IsMsg()

	e := events.ErrorEvent{}
	e.IsMsg()

	k := events.KeyPressEvent{}
	k.IsMsg()

	r := events.ResizeEvent{}
	r.IsMsg()
}

func TestListenDoesNotBlock(t *testing.T) {
	eventChan := make(chan events.Event)
	r := bufio.NewReader(strings.NewReader(""))

	done := make(chan struct{})
	go func() {
		events.Listen(eventChan, r)
		close(done)
	}()

	select {
	case <-done:
		// Pass: Listen returned immediately.
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Listen failed to return immediately")
	}
}
