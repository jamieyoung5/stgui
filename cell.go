package stgui

import (
	"fmt"
	"strconv"
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

	// Style wraps the whole cell in an ANSI sequence - Bg256(137) for a dark
	// square, say.
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

// StyleSeq is the sequence the cell draws with. The selection highlight goes on
// the end of the cell's own style.
func (c *Cell) StyleSeq() string {
	if c.Selected {
		return c.Style + selectedStyle
	}
	return c.Style
}

// content is the cell's lines, unstyled. empty stands in for an empty cell.
func (c *Cell) content(empty string) []string {
	return c.appendContent(nil, empty)
}

// appendContent appends the cell's unstyled lines to dst. Grids pass their frame
// buffer, so a cell costs no allocation of its own.
func (c *Cell) appendContent(dst []string, empty string) []string {
	if c.Child != nil {
		lines := c.Child.RenderLines()
		if len(lines) == 0 {
			return append(dst, empty)
		}

		return append(dst, lines...)
	}

	if c.Value == nil {
		return append(dst, empty)
	}

	return appendLines(dst, valueString(c.Value))
}

// valueString is fmt.Sprint for the types a cell usually holds, skipping the
// reflection for them.
func valueString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		return strconv.FormatBool(v)
	}

	return fmt.Sprint(value)
}

// appendLines splits s on newlines into dst, with nothing allocated for the
// usual single-line value.
func appendLines(dst []string, s string) []string {
	for {
		i := strings.IndexByte(s, '\n')
		if i < 0 {
			return append(dst, s)
		}

		dst = append(dst, s[:i])
		s = s[i+1:]
	}
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
