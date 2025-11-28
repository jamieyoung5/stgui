package stgui

import (
	"fmt"
)

const noValue = "."
const whiteBGBlack = "\033[47m\033[30m" // White background, Black text
const resetStyle = "\033[0m"

type Renderable interface {
	RenderLines() []string
}

type Cell struct {
	Value any
	Child Renderable

	Selected bool
}

func (c *Cell) RenderLines() []string {
	var lines []string

	if c.Child != nil {
		lines = c.Child.RenderLines()
	} else if c.Value != nil {
		lines = []string{fmt.Sprint(c.Value)}
	} else {
		lines = []string{noValue}
	}

	if c.Selected {
		styledLines := make([]string, len(lines))
		for i, line := range lines {
			styledLines[i] = whiteBGBlack + line + resetStyle
		}
		return styledLines
	}

	return lines
}
