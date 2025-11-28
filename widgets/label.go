package widgets

import "strings"

type Label struct {
	Text string
}

func NewLabel(text string) *Label {
	return &Label{Text: text}
}

func (l *Label) RenderLines() []string {
	return strings.Split(l.Text, "\n")
}
