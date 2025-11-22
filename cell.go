package stgui

import "fmt"

const noValue = "."

type Renderable interface {
	RenderLines() []string
}

type Cell struct {
	Value any
	Child Renderable

	Right *Cell
	Down  *Cell
}

func (c *Cell) RenderLines() []string {
	if c.Child != nil {
		return c.Child.RenderLines()
	}
	if c.Value != nil {
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
