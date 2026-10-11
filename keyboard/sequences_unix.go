//go:build !windows

package keyboard

import (
	"maps"
	"slices"

	"github.com/jamieyoung5/gostrc"
)

var (
	Sequences = map[string]Key{
		"\x1b[A": KeyUp,
		"\x1b[B": KeyDown,
		"\x1b[C": KeyRight,
		"\x1b[D": KeyLeft,

		"\x1bOA": KeyUp,
		"\x1bOB": KeyDown,
		"\x1bOC": KeyRight,
		"\x1bOD": KeyLeft,
	}

	sequencesTrie = gostrc.NewTrie(slices.Collect(maps.Keys(Sequences)))
)
