//go:build !windows

package keyboard

import (
	"maps"
	"slices"

	"github.com/jamieyoung5/gostrc"
)

var (
	Sequences = map[string]string{
		"\x1b[A": UpArrowKey,
		"\x1b[B": DownArrowKey,
		"\x1b[C": RightArrowKey,
		"\x1b[D": LeftArrowKey,
	}

	sequencesTrie = gostrc.NewTrie(slices.Collect(maps.Keys(Sequences)))
)
