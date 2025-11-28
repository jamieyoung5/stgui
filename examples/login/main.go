package main

import (
	"fmt"
	"os"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/widgets"
)

func main() {
	titleLabel := widgets.NewLabel(" --- SYSTEM LOGIN --- ")

	userLabel := widgets.NewLabel("Username:")
	userInput := widgets.NewInput("guest", 20)

	passLabel := widgets.NewLabel("Password:")
	passInput := widgets.NewInput("", 20)

	statusLabel := widgets.NewLabel("Status: Idle")

	loginBtn := widgets.NewButton("Login", func() {
		if userInput.Value == "admin" && passInput.Value == "secret" {
			statusLabel.Text = "Status: ACCESS GRANTED"
		} else {
			statusLabel.Text = fmt.Sprintf("Status: Denied (%s)", userInput.Value)
		}
	})

	quitBtn := widgets.NewButton("Exit", func() {
		fmt.Print("\033[?25h")
		os.Exit(0)
	})

	gridData := [][]any{
		{nil, titleLabel, nil},
		{userLabel, userInput, nil},
		{passLabel, passInput, nil},
		{nil, statusLabel, nil},
		{loginBtn, nil, quitBtn},
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
		panic(err)
	}
}
