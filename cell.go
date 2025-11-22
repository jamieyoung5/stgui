package stgui

import "fmt"

const noValue = "."

type Layout interface {
	RenderLines() []string
}

type Cell struct {
	Value   any
	SubGrid Layout

	Right *Cell
	Down  *Cell
}

func (c *Cell) RenderLines() []string {
	if c.SubGrid != nil {
		return c.SubGrid.RenderLines()
	}
	if c.Value != nil {
		return []string{fmt.Sprint(c.Value)}
	}

	return []string{noValue} // empty cell
}

func (c *Cell) Length() int {
	cell := c.Right

	var count int
	for cell != nil {
		count++
		cell = cell.Right
	}

	return count
}

func (c *Cell) Depth() int {
	cell := c.Down

	var count int
	for cell != nil {
		count++
		cell = cell.Down
	}

	return count
}
