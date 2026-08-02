package stgui

import (
	"fmt"
	"strings"
)

const noValue = "."
const selectedStyle = StyleBGWhite + StyleBlack
const resetStyle = "\033[0m"

// Renderable is anything that can draw itself into a cell, one string per line.
type Renderable interface {
	RenderLines() []string
}

type Cell struct {
	Value any
	Child Renderable

	// Style wraps the whole cell in an ANSI sequence - Bg256(94) for a dark
	// square, say. Ignored while the cell is selected.
	Style string

	Selected bool
}

// RenderLines draws whatever the cell holds, styled, one string per line. None of
// them contain a newline; the grid needs them separate to lay rows out.
//
// Grids don't call this. They pad first and style after, so that the style covers
// the padding as well - this is the unpadded version, for a lone cell.
func (c *Cell) RenderLines() []string {
	return applyStyle(c.content(noValue), c.StyleSeq())
}

// IsEmpty reports whether the cell has nothing to draw.
func (c *Cell) IsEmpty() bool {
	return c.Child == nil && c.Value == nil
}

// StyleSeq is the sequence the cell draws with. Selection beats Style.
func (c *Cell) StyleSeq() string {
	if c.Selected {
		return selectedStyle
	}
	return c.Style
}

// content is the cell's lines, unstyled. empty stands in for an empty cell.
func (c *Cell) content(empty string) []string {
	var lines []string

	if c.Child != nil {
		lines = c.Child.RenderLines()
	} else if c.Value != nil {
		lines = strings.Split(fmt.Sprint(c.Value), "\n")
	}

	if len(lines) == 0 {
		lines = []string{empty}
	}

	return lines
}

func applyStyle(lines []string, style string) []string {
	if style == "" {
		return lines
	}

	styled := make([]string, len(lines))
	for i, line := range lines {
		styled[i] = style + line + resetStyle
	}

	return styled
}
