package stgui

import (
	"fmt"
	"strings"
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

	Up    *Cell
	Down  *Cell
	Left  *Cell
	Right *Cell
}

func (c *Cell) RenderLines() []string {
	if c.Child != nil {
		return c.Child.RenderLines()
	}
	if c.Value != nil {
		if c.Selected {
			strValue := fmt.Sprint(c.Value)
			trimmed := strings.TrimSpace(fmt.Sprint(c.Value))
			coloured := fmt.Sprint(whiteBGBlack, trimmed, resetStyle)
			return []string{strings.ReplaceAll(strValue, trimmed, coloured)}
		}
		return []string{fmt.Sprint(c.Value)}
	}

	return []string{noValue} // empty cell
}

func (c *Cell) Width() int {
	cell := c.Right

	var count int
	for cell != nil {
		count++
		cell = cell.Right
	}

	return count
}

func (c *Cell) Height() int {
	cell := c.Down

	var count int
	for cell != nil {
		count++
		cell = cell.Down
	}

	return count
}
