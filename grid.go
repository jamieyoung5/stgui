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
	Root  *Cell
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
	cols := len(grid[0])
	cellGrid, err := createGridCells(rows, cols, grid)
	if err != nil {
		return nil, err
	}

	return &Grid{
		Root:  cellGrid,
		Style: style,
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
	if g.Root == nil {
		return []string{g.Style.NoValue}
	}

	cellBlocks := [][][]string{}
	rowStart := g.Root
	numRows := 0
	numCols := 0

	for rowStart != nil {
		numRows++
		cellBlocks = append(cellBlocks, [][]string{})
		rowBlocks := &cellBlocks[len(cellBlocks)-1]

		cell := rowStart
		currentNumCols := 0
		for cell != nil {
			currentNumCols++
			*rowBlocks = append(*rowBlocks, cell.RenderLines())
			cell = cell.Right
		}

		if currentNumCols > numCols {
			numCols = currentNumCols
		}
		rowStart = rowStart.Down
	}

	if numRows == 0 || numCols == 0 {
		return []string{g.Style.NoValue}
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
	for r := 0; r < numRows; r++ {
		paddedBlocks[r] = make([][]string, numCols)
		for c := 0; c < numCols; c++ {
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

	for r := 0; r < numRows; r++ {
		for h := 0; h < maxHeights[r]; h++ {
			var sb strings.Builder
			for c := 0; c < numCols; c++ {
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
			for c := 0; c < numCols; c++ {
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
	newGrid, err := createGridCells(
		len(grid),
		len(grid[0]),
		grid,
	)
	if err != nil {
		return err
	}

	g.Root = newGrid

	return nil
}

func (g *Grid) Size() (height int, width int) {
	return g.Root.Width(), g.Root.Height()
}

func createGridCells(rows int, cols int, grid [][]any) (*Cell, error) {
	cellGrid := make([][]*Cell, rows)
	for r := range rows {
		cellGrid[r] = make([]*Cell, cols)
		for c := range cols {
			cell := &Cell{}

			element := grid[r][c]
			if renderableElement, ok := element.(Renderable); ok {
				cell.Child = renderableElement
			} else {
				cell.Value = element
			}

			cellGrid[r][c] = cell
		}
	}

	for r := range rows {
		for c := range cols {
			cell := cellGrid[r][c]

			if r > 0 {
				cell.Up = cellGrid[r-1][c]
			}
			if r < rows-1 {
				cell.Down = cellGrid[r+1][c]
			}
			if c > 0 {
				cell.Left = cellGrid[r][c-1]
			}
			if c < cols-1 {
				cell.Right = cellGrid[r][c+1]
			}
		}
	}

	return cellGrid[0][0], nil
}
