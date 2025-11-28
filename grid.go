package stgui

import (
	"regexp"
	"strings"

	"github.com/jamieyoung5/gostrc/strutil"
)

type GridStyle struct {
	VerticalDivider   string
	HorizontalDivider string
	Intersection      string
	NoValue           string
}

type Grid struct {
	Cells [][]*Cell
	Style *GridStyle

	width, height int
}

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleLen(str string) int {
	return len(ansiEscape.ReplaceAllString(str, ""))
}

// TODO: add guarantees that this will be evenly sized (i.e it is a square/rectangle)
func NewGrid(grid [][]any, style *GridStyle) (*Grid, error) {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return &Grid{}, nil
	}
	if style == nil {
		style = &GridStyle{}
	}

	rows := len(grid)

	cols := 0
	for _, row := range grid {
		if len(row) > cols {
			cols = len(row)
		}
	}

	if cols == 0 {
		return &Grid{}, nil
	}

	cellGrid := createGridCells(rows, cols, grid)

	return &Grid{
		Cells:  cellGrid,
		Style:  style,
		width:  cols,
		height: rows,
	}, nil
}

func WithGridSymbols() *GridStyle {
	return &GridStyle{
		VerticalDivider:   "|",
		HorizontalDivider: "-",
		Intersection:      "+",
		NoValue:           ".",
	}
}

func (g *Grid) Render() string {
	return strings.Join(g.RenderLines(), "\n")
}

func (g *Grid) RenderLines() []string {
	if len(g.Cells) == 0 {
		return []string{g.Style.NoValue}
	}

	numRows := g.height
	numCols := g.width

	cellBlocks := make([][][]string, numRows)
	for r := 0; r < numRows; r++ {
		cellBlocks[r] = make([][]string, numCols)
		for c := 0; c < numCols; c++ {
			cellBlocks[r][c] = g.Cells[r][c].RenderLines()
		}
	}

	maxHeights := make([]int, numRows)
	maxWidths := make([]int, numCols)

	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			if c >= len(cellBlocks[r]) {
				continue
			}

			block := cellBlocks[r][c]
			h := len(block)
			w := strutil.MaxVisibleLen(block)

			if h > maxHeights[r] {
				maxHeights[r] = h
			}
			if w > maxWidths[c] {
				maxWidths[c] = w
			}
		}
	}

	vDivider := " " + g.Style.VerticalDivider + " "
	hDivider := g.Style.HorizontalDivider
	cDivider := g.Style.Intersection

	hJoin := strings.Repeat(hDivider, len(vDivider)/2) +
		cDivider +
		strings.Repeat(hDivider, len(vDivider)-(len(vDivider)/2)-1)

	paddedBlocks := make([][][]string, numRows)
	for r := range numRows {
		paddedBlocks[r] = make([][]string, numCols)
		for c := range numCols {
			var rows []string
			if c < len(cellBlocks[r]) {
				rows = cellBlocks[r][c]
			}
			paddedBlocks[r][c] = strutil.PadVisibleRows(rows, maxHeights[r], maxWidths[c], g.Style.NoValue)
		}
	}

	totalOutputLines := 0
	if numRows > 0 {
		totalOutputLines = numRows + (numRows - 1) // Rows + Dividers
	}
	output := make([]string, totalOutputLines)

	currentLineIndex := 0

	for r := range numRows {
		for h := 0; h < maxHeights[r]; h++ {
			var sb strings.Builder
			for c := range numCols {
				sb.WriteString(paddedBlocks[r][c][h])
				if c < numCols-1 {
					sb.WriteString(vDivider)
				}
			}
			output[currentLineIndex] = sb.String()
			currentLineIndex++
		}

		if r < numRows-1 {
			var sb strings.Builder
			for c := range numCols {
				sb.WriteString(strings.Repeat(hDivider, maxWidths[c]))
				if c < numCols-1 {
					sb.WriteString(hJoin)
				}
			}
			output[currentLineIndex] = sb.String()
			currentLineIndex++
		}
	}

	return output
}

func (g *Grid) Refresh(grid [][]any) error {
	rows := len(grid)
	cols := 0
	if rows > 0 {
		for _, row := range grid {
			if len(row) > cols {
				cols = len(row)
			}
		}
	}

	if cols == 0 {
		g.Cells = [][]*Cell{}
		g.width = 0
		g.height = 0
		return nil
	}

	newGrid := createGridCells(rows, cols, grid)

	g.Cells = newGrid
	g.width = cols
	g.height = rows
	return nil
}

func (g *Grid) Size() (height int, width int) {
	return g.height, g.width
}

func createGridCells(rows int, cols int, grid [][]any) [][]*Cell {
	cellGrid := make([][]*Cell, rows)
	for r := range rows {
		cellGrid[r] = make([]*Cell, cols)
		for c := range cols {
			cell := &Cell{}

			var element any
			if c < len(grid[r]) {
				element = grid[r][c]
			} else {
				element = nil
			}

			if renderableElement, ok := element.(Renderable); ok {
				cell.Child = renderableElement
			} else {
				cell.Value = element
			}

			cellGrid[r][c] = cell
		}
	}

	return cellGrid
}
