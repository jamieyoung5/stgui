package stgui

import (
	"strings"
)

// GridStyle says how a grid is drawn. Dividers go in verbatim, so " | " spaces
// the columns out while "|" butts them together. A zero GridStyle puts nothing
// between cells at all, which is what a board of coloured squares wants.
type GridStyle struct {
	VerticalDivider   string // between columns
	HorizontalDivider string // fills the line between rows
	Intersection      string // where the two cross
	NoValue           string // stands in for an empty cell

	// MinCellWidth keeps cells an even size when their contents aren't.
	MinCellWidth int
	// Align defaults to AlignLeft.
	Align Align
	// SelectedStyle replaces the cursor highlight, which is hard to make out
	// against coloured cells. It goes on the end of the cell's own style, so a
	// background alone tints the cell.
	SelectedStyle string

	// Which boundaries get a divider. Nil means all of them. A sudoku marks its
	// boxes with:
	//
	//	DrawColDivider: func(col int) bool { return col%3 == 2 },
	DrawColDivider func(col int) bool
	DrawRowDivider func(row int) bool
}

type Grid struct {
	Cells [][]*Cell
	Style *GridStyle

	width, height int

	// Scratch reused between frames, so drawing is single-goroutine.
	blocks     blockBuf
	rowHeights []int
	colWidths  []int
	gaps       []string
	lineEnds   []int
}

// NewGrid builds a grid from a slice of rows. Ragged input is fine - short rows
// are padded out with empty cells, so the grid is always rectangular. Empty input
// gives an empty grid. Pass a nil style for no dividers and a single space
// between columns.
func NewGrid(grid [][]any, style *GridStyle) *Grid {
	if style == nil {
		style = defaultStyle()
	}

	rows := len(grid)

	cols := 0
	for _, row := range grid {
		if len(row) > cols {
			cols = len(row)
		}
	}

	if rows == 0 || cols == 0 {
		return &Grid{Style: style}
	}

	return &Grid{
		Cells:  createGridCells(rows, cols, grid),
		Style:  style,
		width:  cols,
		height: rows,
	}
}

// NewEmptyGrid builds a rows by cols grid with nothing in it, to be filled in
// with Set as the app runs.
func NewEmptyGrid(rows, cols int, style *GridStyle) *Grid {
	if style == nil {
		style = defaultStyle()
	}

	if rows <= 0 || cols <= 0 {
		return &Grid{Style: style}
	}

	cells := make([][]*Cell, rows)
	for r := range rows {
		cells[r] = make([]*Cell, cols)
		for c := range cols {
			cells[r][c] = &Cell{}
		}
	}

	return &Grid{
		Cells:  cells,
		Style:  style,
		width:  cols,
		height: rows,
	}
}

// WithGridSymbols is the plain ASCII look: "|" and "-" between cells, "+" where
// they meet, "." for an empty cell.
func WithGridSymbols() *GridStyle {
	return &GridStyle{
		VerticalDivider:   " | ",
		HorizontalDivider: "-",
		Intersection:      "+",
		NoValue:           ".",
	}
}

// WithBoxSymbols is the same thing in Unicode box-drawing characters, with cells
// of a fixed width sitting flush against the dividers. Good for a sudoku.
func WithBoxSymbols(cellWidth int) *GridStyle {
	return &GridStyle{
		VerticalDivider:   "│",
		HorizontalDivider: "─",
		Intersection:      "┼",
		NoValue:           " ",
		MinCellWidth:      cellWidth,
		Align:             AlignCenter,
	}
}

// Render draws the grid as one string, lines separated by "\n".
//
// A grid keeps scratch space between frames, so one grid can't be drawn from two
// goroutines at once.
func (g *Grid) Render() string {
	return strings.Join(g.RenderLines(), "\n")
}

// RenderLines draws the grid, one string per terminal line. A cell taller than one
// line stretches its whole row; columns are as wide as their widest cell.
func (g *Grid) RenderLines() []string {
	style := g.style()
	if len(g.Cells) == 0 {
		return []string{style.NoValue}
	}

	numRows, numCols := g.height, g.width

	// Measure every cell unstyled. Padding is written inside the cell's style
	// below, so the colour covers it.
	g.blocks.reset(numRows * numCols)
	rowHeights := resetInts(&g.rowHeights, numRows)
	colWidths := resetInts(&g.colWidths, numCols)

	styleWidth := 0
	for r := range numRows {
		for c := range numCols {
			cell := g.Cells[r][c]

			h, w := g.blocks.add(cell, style.NoValue)
			rowHeights[r] = max(rowHeights[r], h)
			colWidths[c] = max(colWidths[c], w)
			styleWidth = max(styleWidth, len(cell.Style))
		}
	}

	// What a styled cell adds to a line, for sizing the buffer below. Nothing
	// at all when no cell carries a style.
	if styleWidth > 0 {
		styleWidth += len(resetStyle)
	}

	for c := range numCols {
		colWidths[c] = max(colWidths[c], style.MinCellWidth)
	}

	// Work out each boundary separately: a style can mark some and skip others.
	gaps := resetStrings(&g.gaps, max(numCols-1, 0))
	rowWidth := 0
	for c := range numCols {
		rowWidth += colWidths[c]
		if c < numCols-1 {
			gaps[c] = style.columnDivider(c)
			rowWidth += len(gaps[c])
		}
	}

	totalLines := 0
	for r := range numRows {
		totalLines += rowHeights[r]
		if r < numRows-1 && style.drawsRowDivider(r) {
			totalLines++
		}
	}

	// One buffer for the whole frame, sliced into lines at the end. Go strings
	// share their backing array, so the frame costs one allocation and not one
	// per line.
	var sb strings.Builder
	sb.Grow(totalLines*(rowWidth+numCols*styleWidth) + len(selectedStyle) + len(resetStyle))

	ends := g.lineEnds[:0]
	divider := ""

	for r := range numRows {
		for line := range rowHeights[r] {
			for c := range numCols {
				text, width := g.blocks.at(r*numCols+c, line)
				writeCell(&sb, text, width, colWidths[c], style.Align, style.cellStyle(g.Cells[r][c]))

				if c < numCols-1 {
					sb.WriteString(gaps[c])
				}
			}

			ends = append(ends, sb.Len())
		}

		if r < numRows-1 && style.drawsRowDivider(r) {
			// Every divider is the same, so build one and point at it again.
			if divider == "" {
				divider = style.dividerLine(colWidths, gaps)
			}

			// A negative end marks the divider rather than a slice of the
			// frame, which keeps the one copy of it.
			ends = append(ends, -1)
		}
	}

	g.lineEnds = ends

	frame := sb.String()
	output := make([]string, 0, len(ends))

	start := 0
	for _, end := range ends {
		if end < 0 {
			output = append(output, divider)
			continue
		}

		output = append(output, frame[start:end])
		start = end
	}

	return output
}

// Refresh swaps in new contents, keeping the style. The grid may end up a
// different size, so check any cursor pointing into it afterwards.
func (g *Grid) Refresh(grid [][]any) {
	rows := len(grid)
	cols := 0
	for _, row := range grid {
		if len(row) > cols {
			cols = len(row)
		}
	}

	if rows == 0 || cols == 0 {
		g.Cells = nil
		g.width, g.height = 0, 0
		return
	}

	g.Cells = createGridCells(rows, cols, grid)
	g.width, g.height = cols, rows
}

// Size is the grid's dimensions in cells, as (columns, rows).
func (g *Grid) Size() (width int, height int) {
	return g.width, g.height
}

// At is the cell at (row, col), or nil if that's off the grid. Where you go to
// change a Style.
func (g *Grid) At(row, col int) *Cell {
	if row < 0 || row >= len(g.Cells) {
		return nil
	}

	cells := g.Cells[row]
	if col < 0 || col >= len(cells) {
		return nil
	}

	return cells[col]
}

// Set puts a value in a cell, leaving its style and selection alone. Renderable
// values become child widgets, everything else is drawn with fmt.Sprint. Writing
// off the grid does nothing.
func (g *Grid) Set(row, col int, value any) {
	if cell := g.At(row, col); cell != nil {
		setCellContent(cell, value)
	}
}

// Get is the cell's child widget, or its value, or nil if there's neither.
func (g *Grid) Get(row, col int) any {
	cell := g.At(row, col)
	if cell == nil {
		return nil
	}
	if cell.Child != nil {
		return cell.Child
	}

	return cell.Value
}

// Each walks the cells in reading order. Saves a nested loop when setting up or
// restyling a whole board.
func (g *Grid) Each(fn func(row, col int, cell *Cell)) {
	for r, cells := range g.Cells {
		for c, cell := range cells {
			fn(r, c, cell)
		}
	}
}

// style falls back to a plain one, so a Grid built by hand still draws.
func (g *Grid) style() *GridStyle {
	if g.Style == nil {
		return &plainStyle
	}
	return g.Style
}

// plainStyle: no dividers, but a column of space so values don't run together.
var plainStyle = GridStyle{VerticalDivider: " "}

func defaultStyle() *GridStyle {
	return &GridStyle{VerticalDivider: " "}
}

// columnDivider is what goes between col and the next column along.
func (s *GridStyle) columnDivider(col int) string {
	if s.DrawColDivider != nil && !s.DrawColDivider(col) {
		return ""
	}
	return s.VerticalDivider
}

// drawsRowDivider reports whether a divider line follows row. Without a
// horizontal character there's nothing to draw it with.
func (s *GridStyle) drawsRowDivider(row int) bool {
	if s.HorizontalDivider == "" && s.Intersection == "" {
		return false
	}
	if s.DrawRowDivider != nil {
		return s.DrawRowDivider(row)
	}

	return true
}

// filler is what a row divider is made of. A space when the style has no
// horizontal character, so an intersection-only style still lines up.
func (s *GridStyle) filler() string {
	if s.HorizontalDivider == "" {
		return " "
	}
	return s.HorizontalDivider
}

// writeDividerJoin writes the bit of row divider that crosses a column gap.
// Sized to the gap, or the two come out misaligned.
func (s *GridStyle) writeDividerJoin(sb *strings.Builder, gapWidth int) {
	if gapWidth <= 0 {
		return
	}
	if s.Intersection == "" {
		writeRepeat(sb, s.filler(), gapWidth)
		return
	}

	left := (gapWidth - 1) / 2

	writeRepeat(sb, s.filler(), left)
	sb.WriteString(s.Intersection)
	writeRepeat(sb, s.filler(), gapWidth-left-1)
}

// cellStyle lets the grid's style override the cursor highlight.
func (s *GridStyle) cellStyle(cell *Cell) string {
	if cell.Selected && s.SelectedStyle != "" {
		return cell.Style + s.SelectedStyle
	}
	return cell.StyleSeq()
}

func createGridCells(rows int, cols int, grid [][]any) [][]*Cell {
	cellGrid := make([][]*Cell, rows)
	for r := range rows {
		cellGrid[r] = make([]*Cell, cols)
		for c := range cols {
			cell := &Cell{}

			if c < len(grid[r]) {
				setCellContent(cell, grid[r][c])
			}

			cellGrid[r][c] = cell
		}
	}

	return cellGrid
}

func setCellContent(cell *Cell, element any) {
	if renderable, ok := element.(Renderable); ok {
		cell.Child, cell.Value = renderable, nil
		return
	}

	cell.Child, cell.Value = nil, element
}

// dividerLine is the line that goes between two rows. It depends only on the
// column widths, so a grid builds one and reuses it.
func (s *GridStyle) dividerLine(colWidths []int, gaps []string) string {
	filler := s.filler()

	width := 0
	for _, w := range colWidths {
		width += w
	}
	for _, gap := range gaps {
		width += VisibleWidth(gap)
	}

	var sb strings.Builder
	sb.Grow(width * len(filler))

	for c, w := range colWidths {
		writeRepeat(&sb, filler, w)
		if c < len(colWidths)-1 {
			s.writeDividerJoin(&sb, VisibleWidth(gaps[c]))
		}
	}

	return sb.String()
}

// writeCell writes one line of a cell: its style, the text padded out to width,
// then a reset. Padding goes inside the style so the colour covers it.
func writeCell(sb *strings.Builder, text string, textWidth, width int, align Align, style string) {
	left, right := 0, 0
	if gap := width - textWidth; gap > 0 {
		switch align {
		case AlignRight:
			left = gap
		case AlignCenter:
			left = gap / 2
			right = gap - left
		default:
			right = gap
		}
	}

	if style != "" {
		sb.WriteString(style)
	}

	writeSpaces(sb, left)
	sb.WriteString(text)
	writeSpaces(sb, right)

	if style != "" {
		sb.WriteString(resetStyle)
	}
}

func writeRepeat(sb *strings.Builder, s string, n int) {
	for range n {
		sb.WriteString(s)
	}
}

const spaceRun = "                                                                "

func writeSpaces(sb *strings.Builder, n int) {
	for n > len(spaceRun) {
		sb.WriteString(spaceRun)
		n -= len(spaceRun)
	}
	if n > 0 {
		sb.WriteString(spaceRun[:n])
	}
}
