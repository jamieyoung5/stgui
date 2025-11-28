package testutil

import "github.com/jamieyoung5/stgui"

type MockRenderable struct {
	Lines []string
}

func (m *MockRenderable) RenderLines() []string {
	return m.Lines
}

type MockWidget struct {
	W, H         int
	RenderOutput string
	LastInput    string
	ReturnScreen *stgui.Screen
	ReturnExit   bool
}

func (m *MockWidget) Size() (int, int) {
	return m.W, m.H
}

func (m *MockWidget) Render() string {
	return m.RenderOutput
}

func (m *MockWidget) Select(cursor *stgui.Cursor, input string) (*stgui.Screen, bool) {
	m.LastInput = input
	return m.ReturnScreen, m.ReturnExit
}
