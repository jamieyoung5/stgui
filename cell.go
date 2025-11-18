package stgui

import "fmt"

type Layout interface {
	RenderLines() []string
}

type Cell struct {
	Value   any
	SubGrid Layout

	Right   *Cell
	Down    *Cell
	Symbols *Symbols
}

func (c *Cell) RenderLines() []string {
	if c.SubGrid != nil {
		return c.SubGrid.RenderLines()
	}
	if c.Value != nil {
		return []string{fmt.Sprint(c.Value)}
	}

	return []string{c.Symbols.NoValue} // empty cell
}
