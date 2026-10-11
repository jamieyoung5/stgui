package main

import (
	"fmt"
	"os"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/widgets"
)

func main() {
	if err := stgui.NewApp(buildLoginScreen()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func buildLoginScreen() *stgui.Screen {
	userInput := widgets.NewInput("guest", 20)
	passInput := widgets.NewMaskedInput("", 20)
	statusLabel := widgets.NewLabel("Status: Idle")

	// Build the welcome screen up front; the login button points at it once the
	// details check out.
	welcome := buildWelcomeScreen()

	// Checks the inputs when pressed, and decides where to go from there.
	loginBtn := widgets.NewActionButton("Login", func() stgui.Action {
		if userInput.Value == "admin" && passInput.Value == "secret" {
			return stgui.Push(welcome)
		}

		statusLabel.Text = fmt.Sprintf("Status: Denied (%s)", userInput.Value)
		return stgui.Stay()
	})

	gridData := [][]any{
		{nil, widgets.NewLabel(" --- SYSTEM LOGIN --- "), nil},
		{widgets.NewLabel("Username:"), userInput, nil},
		{widgets.NewLabel("Password:"), passInput, nil},
		{nil, statusLabel, nil},
		{loginBtn, nil, widgets.NewQuitButton("Exit")},
	}

	grid := stgui.NewGrid(gridData, stgui.WithGridSymbols())

	// Start focused on the username field.
	return widgets.NewGridScreen(grid, 1, 1)
}

func buildWelcomeScreen() *stgui.Screen {
	gridData := [][]any{
		{widgets.NewLabel("Access granted.\nWelcome, admin.")},
		{widgets.NewQuitButton("Log out")},
	}

	grid := stgui.NewGrid(gridData, stgui.WithGridSymbols())

	return widgets.NewGridScreen(grid, 1, 0)
}
