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
	// against coloured cells.
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
}

// NewGrid builds a grid from a slice of rows. Ragged input is fine - short rows
// are padded out with empty cells, so the grid is always rectangular. Pass a nil
// style for no dividers and a single space between columns.
func NewGrid(grid [][]any, style *GridStyle) (*Grid, error) {
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
		return &Grid{Style: style}, nil
	}

	return &Grid{
		Cells:  createGridCells(rows, cols, grid),
		Style:  style,
		width:  cols,
		height: rows,
	}, nil
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

	cellBlocks := make([][][]string, numRows)
	rowHeights := make([]int, numRows)
	colWidths := make([]int, numCols)

	for r := range numRows {
		cellBlocks[r] = make([][]string, numCols)
		for c := range numCols {
			// Measure and pad while unstyled, then style at the end, so the
			// colour covers the padding and not just the text.
			block := g.Cells[r][c].content(style.NoValue)
			cellBlocks[r][c] = block

			if h := len(block); h > rowHeights[r] {
				rowHeights[r] = h
			}
			if w := maxVisibleWidth(block); w > colWidths[c] {
				colWidths[c] = w
			}
		}
	}

	for c := range numCols {
		colWidths[c] = max(colWidths[c], style.MinCellWidth)
	}

	// Work out each boundary separately: a style can mark some and skip others.
	gaps := make([]string, max(numCols-1, 0))
	for c := range gaps {
		gaps[c] = style.columnDivider(c)
	}

	output := make([]string, 0, numRows)

	for r := range numRows {
		padded := make([][]string, numCols)
		for c := range numCols {
			block := padBlock(cellBlocks[r][c], rowHeights[r], colWidths[c], style.Align)
			padded[c] = applyStyle(block, style.cellStyle(g.Cells[r][c]))
		}

		for line := range rowHeights[r] {
			var sb strings.Builder
			for c := range numCols {
				sb.WriteString(padded[c][line])
				if c < numCols-1 {
					sb.WriteString(gaps[c])
				}
			}
			output = append(output, sb.String())
		}

		if r < numRows-1 && style.drawsRowDivider(r) {
			var sb strings.Builder
			for c := range numCols {
				sb.WriteString(strings.Repeat(style.filler(), colWidths[c]))
				if c < numCols-1 {
					sb.WriteString(style.dividerJoin(VisibleWidth(gaps[c])))
				}
			}
			output = append(output, sb.String())
		}
	}

	return output
}

// Refresh swaps in new contents, keeping the style. The grid may end up a
// different size, so check any cursor pointing into it afterwards.
func (g *Grid) Refresh(grid [][]any) error {
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
		return nil
	}

	g.Cells = createGridCells(rows, cols, grid)
	g.width, g.height = cols, rows
	return nil
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
		return defaultStyle()
	}
	return g.Style
}

// defaultStyle: no dividers, but a column of space so values don't run together.
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

// dividerJoin is the bit of row divider that crosses a column gap. Sized to the
// gap, or the two come out misaligned.
func (s *GridStyle) dividerJoin(gapWidth int) string {
	if gapWidth <= 0 {
		return ""
	}
	if s.Intersection == "" {
		return strings.Repeat(s.filler(), gapWidth)
	}

	left := (gapWidth - 1) / 2

	return strings.Repeat(s.filler(), left) +
		s.Intersection +
		strings.Repeat(s.filler(), gapWidth-left-1)
}

// cellStyle lets the grid's style override the cursor highlight.
func (s *GridStyle) cellStyle(cell *Cell) string {
	if cell.Selected && s.SelectedStyle != "" {
		return s.SelectedStyle
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

// padBlock fills a cell out to height rows of width columns.
func padBlock(lines []string, height, width int, align Align) []string {
	padded := make([]string, height)
	for i := range height {
		var line string
		if i < len(lines) {
			line = lines[i]
		}
		padded[i] = padLine(line, width, align)
	}

	return padded
}

func padLine(line string, width int, align Align) string {
	gap := width - VisibleWidth(line)
	if gap <= 0 {
		return line
	}

	switch align {
	case AlignRight:
		return strings.Repeat(" ", gap) + line
	case AlignCenter:
		left := gap / 2
		return strings.Repeat(" ", left) + line + strings.Repeat(" ", gap-left)
	default:
		return line + strings.Repeat(" ", gap)
	}
}
