package events

import (
	"strconv"
	"strings"
)

// MouseEvent is a button going down or up. Col and Row are terminal cells from
// the top left, counted from zero like everything else.
type MouseEvent struct {
	Msg

	Col, Row int
	Button   int
	Pressed  bool
}

// mousePrefix starts an SGR mouse report: ESC [ < button ; col ; row then M for
// a press or m for a release.
const mousePrefix = "\x1b[<"

// parseMouse reads an SGR mouse report. It reports false for anything else, and
// for the motion and wheel reports we don't ask the terminal for.
func parseMouse(seq string) (MouseEvent, bool) {
	if !strings.HasPrefix(seq, mousePrefix) || len(seq) < len(mousePrefix)+1 {
		return MouseEvent{}, false
	}

	end := seq[len(seq)-1]
	if end != 'M' && end != 'm' {
		return MouseEvent{}, false
	}

	fields := strings.Split(seq[len(mousePrefix):len(seq)-1], ";")
	if len(fields) != 3 {
		return MouseEvent{}, false
	}

	numbers := make([]int, 3)
	for i, field := range fields {
		n, err := strconv.Atoi(field)
		if err != nil {
			return MouseEvent{}, false
		}

		numbers[i] = n
	}

	button := numbers[0]

	// Bit 5 is movement and bit 6 the wheel. Neither is a click.
	if button&32 != 0 || button&64 != 0 {
		return MouseEvent{}, false
	}

	// The terminal counts from one.
	return MouseEvent{
		Col:     numbers[1] - 1,
		Row:     numbers[2] - 1,
		Button:  button & 3,
		Pressed: end == 'M',
	}, true
}
