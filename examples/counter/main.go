package main

import (
	"fmt"
	"os"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/widgets"
)

func main() {
	count := 0
	countLabel := widgets.NewLabel("")

	updateDisplay := func() {
		countLabel.Text = fmt.Sprintf("Current Count: %d", count)
	}
	updateDisplay()

	incBtn := widgets.NewButton("Increment (+)", func() {
		count++
		updateDisplay()
	})

	decBtn := widgets.NewButton("Decrement (-)", func() {
		count--
		updateDisplay()
	})

	resetBtn := widgets.NewButton("Reset", func() {
		count = 0
		updateDisplay()
	})

	quitBtn := widgets.NewQuitButton("Quit")

	gridData := [][]any{
		{decBtn, countLabel, incBtn},
		{resetBtn, nil, quitBtn},
	}

	grid := stgui.NewGrid(gridData, stgui.WithGridSymbols())

	// Start on the label in the middle of the top row.
	screen := widgets.NewGridScreen(grid, 0, 1)

	if err := stgui.NewApp(screen).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
