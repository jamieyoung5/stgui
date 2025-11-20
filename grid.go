package stgui

import (
	"strings"
)

type Symbols struct {
	VerticalDivider     string
	HorizontalDivider   string
	CrossSectionDivider string
	NoValue             string
}

type Grid struct {
	Content *Cell
	Symbols *Symbols
}

func NewGrid(grid [][]any, symbols *Symbols) (*Grid, error) {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return &Grid{}, nil
	}
	if symbols == nil {
		symbols = &Symbols{}
	}

	rows := len(grid)
	cols := len(grid[0])
	cellGrid, err := createGridCells(rows, cols, grid)
	if err != nil {
		return nil, err
	}

	return &Grid{
		Content: cellGrid,
		Symbols: symbols,
	}, nil
}

func WithGridSymbols() *Symbols {
	return &Symbols{
		VerticalDivider:     "|",
		HorizontalDivider:   "-",
		CrossSectionDivider: "+",
		NoValue:             ".",
	}
}

func (g *Grid) Render() string { return Render(g) }

func (g *Grid) RenderLines() []string {
	if g.Content == nil {
		return []string{g.Symbols.NoValue}
	}

	cellBlocks := [][][]string{}
	rowStart := g.Content
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
		return []string{g.Symbols.NoValue}
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
			w := MaxLength(block)

			if h > maxHeights[r] {
				maxHeights[r] = h
			}
			if w > maxWidths[c] {
				maxWidths[c] = w
			}
		}
	}

	vDivider := " " + g.Symbols.VerticalDivider + " "
	hDivider := g.Symbols.HorizontalDivider
	cDivider := g.Symbols.CrossSectionDivider

	hJoin := strings.Repeat(hDivider, len(vDivider)/2) +
		cDivider +
		strings.Repeat(hDivider, len(vDivider)-(len(vDivider)/2)-1)

	paddedBlocks := make([][][]string, numRows)
	for r := 0; r < numRows; r++ {
		paddedBlocks[r] = make([][]string, numCols)
		for c := 0; c < numCols; c++ {
			var block []string
			if c < len(cellBlocks[r]) {
				block = cellBlocks[r][c]
			}
			paddedBlocks[r][c] = PadBlock(block, maxHeights[r], maxWidths[c], g.Symbols.NoValue)
		}
	}

	output := make([]string, numRows)
	for r := 0; r < numRows; r++ {
		for h := 0; h < maxHeights[r]; h++ {
			var sb strings.Builder
			for c := 0; c < numCols; c++ {
				sb.WriteString(paddedBlocks[r][c][h])
				if c < numCols-1 {
					sb.WriteString(vDivider)
				}
			}
			output[r] = sb.String()
		}

		if r < numRows-1 {
			var sb strings.Builder
			for c := 0; c < numCols; c++ {
				sb.WriteString(strings.Repeat(hDivider, maxWidths[c]))
				if c < numCols-1 {
					sb.WriteString(hJoin)
				}
			}
			output[r] = sb.String()
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

	g.Content = newGrid

	return nil
}

func createGridCells(rows int, cols int, grid [][]any) (*Cell, error) {

	cells := make([]Cell, rows*cols)
	cellGrid := make([][]*Cell, rows)
	pointers := make([]*Cell, rows*cols)

	for r := range rows {
		rowStart := r * cols
		cellGrid[r] = pointers[rowStart : rowStart+cols]

		for c := range cols {
			index := r*cols + c

			cell := &cells[index]

			cellGrid[r][c] = cell

			element := grid[r][c]
			if layout, ok := element.(Layout); ok {
				cell.SubGrid = layout
			} else {
				cell.Value = element
			}
		}
	}

	for r := range rows {
		for c := range cols {
			cell := cellGrid[r][c]

			if r < rows-1 {
				cell.Down = cellGrid[r+1][c]
			}
			if c < cols-1 {
				cell.Right = cellGrid[r][c+1]
			}
		}
	}

	return cellGrid[0][0], nil
}
