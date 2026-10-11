//go:build windows

package keyboard

import (
	"maps"
	"slices"

	"github.com/jamieyoung5/gostrc"
)

var (
	Sequences = map[string]Key{
		// Legacy Windows Scancodes (Prefix 0 or 224 + KeyCode)
		"\x00H": KeyUp,
		"\xe0H": KeyUp,
		"\x00P": KeyDown,
		"\xe0P": KeyDown,
		"\x00K": KeyLeft,
		"\xe0K": KeyLeft,
		"\x00M": KeyRight,
		"\xe0M": KeyRight,

		// Modern Windows Terminal (ANSI Fallback)
		"\x1b[A": KeyUp,
		"\x1b[B": KeyDown,
		"\x1b[C": KeyRight,
		"\x1b[D": KeyLeft,

		// Application Cursor Keys
		"\x1bOA": KeyUp,
		"\x1bOB": KeyDown,
		"\x1bOC": KeyRight,
		"\x1bOD": KeyLeft,
	}

	sequencesTrie = gostrc.NewTrie(slices.Collect(maps.Keys(Sequences)))
)
