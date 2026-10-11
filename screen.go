package stgui

import (
	"strings"
)

// Widget is a top-level piece of a screen. It draws itself and decides what each
// input means. widgets.Container is the usual one, backed by a grid.
type Widget interface {
	// Render draws the widget, lines separated by "\n".
	Render() string
	// Select handles one input and returns what to do next. Stay() carries on
	// with the current screen.
	Select(cursor *Cursor, input Input) Action
}

// Screen is one terminal-full: widgets laid out in rows, plus the cursor moving
// through them.
type Screen struct {
	// Cursor is where input goes. A screen without one draws but ignores
	// every key.
	Cursor *Cursor

	widgets [][]Widget

	// Scratch reused between frames, like Grid's.
	blocks   [][]string
	measured blockBuf
	widths   []int
	lineEnds []int
	height   int
}

// NewScreen lays rows of widgets out top to bottom, with cursor taking the
// input.
func NewScreen(cursor *Cursor, rows ...[]Widget) *Screen {
	return &Screen{
		Cursor:  cursor,
		widgets: rows,
	}
}

// Row is one row of a screen: widgets sitting side by side.
func Row(widgets ...Widget) []Widget {
	return widgets
}

// Render draws the lot: widgets sharing a row go side by side, rows get a blank
// line between them.
//
// Like Grid.Render, this keeps scratch space between frames and so is for one
// goroutine at a time.
func (s *Screen) Render() string {
	return strings.Join(s.renderLines(), "\n")
}

// rowSpacing is the gap between two widgets sharing a row.
const rowSpacing = 4

// renderLines is Render, one string per terminal line. What the app draws from:
// widgets that can hand over lines directly do, so a frame never joins and
// re-splits them.
func (s *Screen) renderLines() []string {
	// Frames are much the same height, so last one's is a good guess.
	out := make([]string, 0, s.height)

	for _, widgetRow := range s.widgets {
		blocks, height := s.renderWidgets(widgetRow)
		if height == 0 {
			continue
		}

		// Blank line between rows.
		if len(out) > 0 {
			out = append(out, "")
		}

		out = s.appendSideBySide(out, blocks, height)
	}

	s.height = len(out)

	return out
}

// renderWidgets draws each widget in a row to its lines, and reports the tallest.
func (s *Screen) renderWidgets(widgets []Widget) (blocks [][]string, height int) {
	blocks = s.blocks[:0]

	for _, elem := range widgets {
		if elem == nil {
			continue
		}

		var lines []string
		if renderable, ok := elem.(Renderable); ok {
			lines = renderable.RenderLines()
		} else {
			lines = strings.Split(elem.Render(), "\n")
		}

		blocks = append(blocks, lines)
		height = max(height, len(lines))
	}

	s.blocks = blocks

	return blocks, height
}

// appendSideBySide lays blocks out as columns, each padded to its widest line so
// the next one starts in the same place on every row.
func (s *Screen) appendSideBySide(out []string, blocks [][]string, height int) []string {
	// A row with one widget in it is already the lines we want, and there is
	// nothing to its right to line up with.
	if len(blocks) == 1 {
		return append(out, blocks[0]...)
	}

	s.measured.reset(len(blocks))
	widths := resetInts(&s.widths, len(blocks))

	content, padding := 0, 0
	for i, block := range blocks {
		_, widths[i] = s.measured.addLines(block)
		padding += widths[i]

		for _, line := range block {
			content += len(line)
		}
	}

	last := len(blocks) - 1

	// One buffer for the row, sliced into lines after, as a grid does.
	var sb strings.Builder
	sb.Grow(content + height*(padding+rowSpacing*last))

	ends := s.lineEnds[:0]

	for line := range height {
		for i := range blocks {
			// Past the end of a shorter widget this is empty, which pads it
			// out to the tallest in the row.
			text, width := s.measured.at(i, line)

			// Nothing follows the last column, so it needs no padding.
			if i == last {
				sb.WriteString(text)
				break
			}

			writeCell(&sb, text, width, widths[i], AlignLeft, "")
			writeSpaces(&sb, rowSpacing)
		}

		ends = append(ends, sb.Len())
	}

	s.lineEnds = ends

	frame := sb.String()
	start := 0
	for _, end := range ends {
		out = append(out, frame[start:end])
		start = end
	}

	return out
}
