//go:build windows

package keyboard

import (
	"maps"
	"slices"

	"github.com/jamieyoung5/gostrc"
)

var (
	Sequences = map[string]string{
		// Legacy Windows Scancodes (Prefix 0 or 224 + KeyCode)
		"\x00H": UpArrowKey,
		"\xe0H": UpArrowKey,
		"\x00P": DownArrowKey,
		"\xe0P": DownArrowKey,
		"\x00K": LeftArrowKey,
		"\xe0K": LeftArrowKey,
		"\x00M": RightArrowKey,
		"\xe0M": RightArrowKey,

		// Modern Windows Terminal (ANSI Fallback)
		"\x1b[A": UpArrowKey,
		"\x1b[B": DownArrowKey,
		"\x1b[C": RightArrowKey,
		"\x1b[D": LeftArrowKey,

		// Application Cursor Keys
		"\x1bOA": UpArrowKey,
		"\x1bOB": DownArrowKey,
		"\x1bOC": RightArrowKey,
		"\x1bOD": LeftArrowKey,
	}

	sequencesTrie = gostrc.NewTrie(slices.Collect(maps.Keys(Sequences)))
)
