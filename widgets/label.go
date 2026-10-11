package widgets

import (
	"strings"

	"github.com/jamieyoung5/stgui"
)

// Label is a piece of text. It works in a grid cell and as a screen row, and
// ignores input.
type Label struct {
	Text string
}

func NewLabel(text string) *Label {
	return &Label{Text: text}
}

func (l *Label) RenderLines() []string {
	return strings.Split(l.Text, "\n")
}

func (l *Label) Render() string {
	return l.Text
}

func (l *Label) Select(cursor *stgui.Cursor, input stgui.Input) stgui.Action {
	return stgui.Stay()
}
