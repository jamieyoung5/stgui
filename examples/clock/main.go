// Clock: a screen that redraws without anyone touching the keyboard. A
// goroutine sends a tick and OnEvent writes it to a label.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/events"
	"github.com/jamieyoung5/stgui/widgets"
)

// tickEvent is our own event. Embedding events.Msg is all it takes.
type tickEvent struct {
	events.Msg
	at time.Time
}

func main() {
	started := time.Now()
	clock := widgets.NewLabel(started.Format(time.TimeOnly))
	uptime := widgets.NewLabel("up 0s")

	quit := widgets.NewContainer(stgui.NewGrid([][]any{
		{widgets.NewQuitButton("Quit")},
	}, nil))

	screen := stgui.NewScreen(quit.Focus(0, 0),
		stgui.Row(widgets.NewLabel("Clock - redraws on its own, q to quit")),
		stgui.Row(clock),
		stgui.Row(uptime),
		stgui.Row(quit),
	)

	// Anything the app doesn't recognise lands in OnEvent.
	app := stgui.NewApp(screen).QuitOn("q").OnEvent(func(event events.Event) stgui.Action {
		tick, ok := event.(tickEvent)
		if !ok {
			return stgui.Stay()
		}

		clock.Text = tick.at.Format(time.TimeOnly)
		uptime.Text = fmt.Sprintf("up %s", tick.at.Sub(started).Round(time.Second))

		return stgui.Stay()
	})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		for t := range ticker.C {
			app.Send(tickEvent{at: t})
		}
	}()

	if err := app.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
