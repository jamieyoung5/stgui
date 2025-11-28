//go:build !windows

package events

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

func listenForResize(eventChan chan Event) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)

	for range sigChan {
		w, h, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			eventChan <- ResizeEvent{Width: w, Height: h}
		}
	}
}
