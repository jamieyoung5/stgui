package main

import (
	"fmt"
	"os"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/widgets"
)

func main() {
	count := 0
	countLabel := widgets.NewLabel(fmt.Sprintf("Current Count: %d", count))

	updateDisplay := func() {
		countLabel.Text = fmt.Sprintf("Current Count: %d", count)
	}

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

	quitBtn := widgets.NewButton("Quit", func() {
		fmt.Print("\033[?25h")
		os.Exit(0)
	})

	gridData := [][]any{
		{decBtn, countLabel, incBtn},
		{resetBtn, nil, quitBtn},
	}

	grid, err := stgui.NewGrid(gridData, stgui.WithGridSymbols())
	if err != nil {
		panic(err)
	}

	container := widgets.NewContainer(grid)

	startRow, startCol := 1, 1

	cursor := stgui.NewCursor(container, grid, startRow, startCol, stgui.DefaultDirectionalControls)

	grid.Cells[startRow][startCol].Selected = true

	screen := stgui.NewScreen(
		[]*stgui.Cursor{cursor},
		[][]stgui.Widget{{container}},
	)

	app := stgui.NewApp(screen)

	if err := app.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
