package testutil

import "github.com/jamieyoung5/stgui"

type MockRenderable struct {
	Lines []string
}

func (m *MockRenderable) RenderLines() []string {
	return m.Lines
}

type MockWidget struct {
	RenderOutput string
	LastInput    stgui.Input
	ReturnAction stgui.Action
}

func (m *MockWidget) Render() string {
	return m.RenderOutput
}

func (m *MockWidget) Select(cursor *stgui.Cursor, input stgui.Input) stgui.Action {
	m.LastInput = input
	return m.ReturnAction
}
