package stgui

import (
	"runtime"
	"testing"
)

func makeLargeGrid(rows, cols int) *Grid {
	gridData := make([][]any, rows)
	for r := range gridData {
		gridData[r] = make([]any, cols)
		for c := range gridData[r] {
			gridData[r][c] = "Item\n" + string(rune('A'+r)) + "-" + string(rune('0'+c))
		}
	}
	g, _ := NewGrid(gridData, WithGridSymbols())
	return g
}

func BenchmarkGridRenderLines(b *testing.B) {
	g := makeLargeGrid(50, 50)

	b.ResetTimer()
	var result []string
	for i := 0; i < b.N; i++ {
		result = g.RenderLines()
	}
	runtime.KeepAlive(result)
}

func BenchmarkCreateGridCells(b *testing.B) {
	const rows, cols = 50, 50
	gridData := make([][]any, rows)
	for r := range gridData {
		gridData[r] = make([]any, cols)
		for c := range gridData[r] {
			gridData[r][c] = "Item " + string(rune('A'+r)) + string(rune('0'+c))
		}
	}

	var result [][]*Cell

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		res, err := createGridCells(rows, cols, gridData)
		if err != nil {
			b.Fatal(err)
		}
		result = res
	}

	runtime.KeepAlive(result)
}
