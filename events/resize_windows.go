//go:build windows

package events

import (
	"os"
	"time"

	"golang.org/x/term"
)

func listenForResize(eventChan chan Event) {
	// Initial size
	lastW, lastH, _ := term.GetSize(int(os.Stdout.Fd()))

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		w, h, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			continue
		}

		if w != lastW || h != lastH {
			lastW = w
			lastH = h
			eventChan <- ResizeEvent{Width: w, Height: h}
		}
	}
}
