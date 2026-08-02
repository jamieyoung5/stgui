package main

import (
	"fmt"
	"os"

	"github.com/jamieyoung5/stgui"
	"github.com/jamieyoung5/stgui/widgets"
)

func main() {
	loginScreen, err := buildLoginScreen()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if err := stgui.NewApp(loginScreen).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func buildLoginScreen() (*stgui.Screen, error) {
	userInput := widgets.NewInput("guest", 20)
	passInput := widgets.NewMaskedInput("", 20)
	statusLabel := widgets.NewLabel("Status: Idle")

	// Build the welcome screen up front; the login button points at it once the
	// details check out.
	welcome, err := buildWelcomeScreen()
	if err != nil {
		return nil, err
	}

	loginBtn := widgets.NewButton("Login", nil)
	loginBtn.Callback = func() {
		if userInput.Value == "admin" && passInput.Value == "secret" {
			loginBtn.Screen = welcome
			return
		}

		loginBtn.Screen = nil
		statusLabel.Text = fmt.Sprintf("Status: Denied (%s)", userInput.Value)
	}

	gridData := [][]any{
		{nil, widgets.NewLabel(" --- SYSTEM LOGIN --- "), nil},
		{widgets.NewLabel("Username:"), userInput, nil},
		{widgets.NewLabel("Password:"), passInput, nil},
		{nil, statusLabel, nil},
		{loginBtn, nil, widgets.NewQuitButton("Exit")},
	}

	grid, err := stgui.NewGrid(gridData, stgui.WithGridSymbols())
	if err != nil {
		return nil, err
	}

	// Start focused on the username field.
	return widgets.NewScreen(grid, 1, 1), nil
}

func buildWelcomeScreen() (*stgui.Screen, error) {
	gridData := [][]any{
		{widgets.NewLabel("Access granted.\nWelcome, admin.")},
		{widgets.NewQuitButton("Log out")},
	}

	grid, err := stgui.NewGrid(gridData, stgui.WithGridSymbols())
	if err != nil {
		return nil, err
	}

	return widgets.NewScreen(grid, 1, 0), nil
}
