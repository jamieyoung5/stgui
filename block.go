package stgui

// blockBuf holds every cell's lines for one frame in a single flat slice, with
// each line's visible width alongside. Cells are measured once and written once,
// and the backing arrays are reused between frames.
type blockBuf struct {
	lines  []string
	widths []int
	starts []int32
	ends   []int32
}

func (b *blockBuf) reset(cells int) {
	b.lines = b.lines[:0]
	b.widths = b.widths[:0]
	b.starts = resetInt32s(&b.starts, cells)[:0]
	b.ends = resetInt32s(&b.ends, cells)[:0]
}

// add stores a cell's lines and reports the space they need.
func (b *blockBuf) add(cell *Cell, empty string) (height, width int) {
	start := b.open()
	b.lines = cell.appendContent(b.lines, empty)

	return b.measure(start)
}

// addLines is add for a block that is already a slice of lines.
func (b *blockBuf) addLines(lines []string) (height, width int) {
	start := b.open()
	b.lines = append(b.lines, lines...)

	return b.measure(start)
}

func (b *blockBuf) open() int {
	start := len(b.lines)
	b.starts = append(b.starts, int32(start))

	return start
}

// measure records the visible width of everything added since start. Scanning for
// escape sequences is the expensive part of a frame, so it happens once here.
func (b *blockBuf) measure(start int) (height, width int) {
	for _, line := range b.lines[start:] {
		w := VisibleWidth(line)
		b.widths = append(b.widths, w)
		if w > width {
			width = w
		}
	}

	b.ends = append(b.ends, int32(len(b.lines)))

	return len(b.lines) - start, width
}

// at is cell i's given line and its visible width. Lines past the end of a
// shorter cell come back empty, which pads it out to its row.
func (b *blockBuf) at(cell, line int) (string, int) {
	i := int(b.starts[cell]) + line
	if i >= int(b.ends[cell]) {
		return "", 0
	}

	return b.lines[i], b.widths[i]
}

// resetInts hands back a zeroed slice of n ints, reusing buf's memory.
func resetInts(buf *[]int, n int) []int {
	if cap(*buf) < n {
		*buf = make([]int, n)
	}

	s := (*buf)[:n]
	clear(s)

	return s
}

func resetStrings(buf *[]string, n int) []string {
	if cap(*buf) < n {
		*buf = make([]string, n)
	}

	return (*buf)[:n]
}

func resetInt32s(buf *[]int32, n int) []int32 {
	if cap(*buf) < n {
		*buf = make([]int32, n)
	}

	return (*buf)[:n]
}
