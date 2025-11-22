package keyboard

import (
	"bufio"
	"unicode"
)

// TODO: handle sequences
func ReadInput(reader *bufio.Reader) (string, error) {
	rn, _, err := reader.ReadRune()
	if err != nil {
		return "", err
	}

	if rn == 27 { // Ignore ESC sequences completely
		if reader.Buffered() > 0 {
			reader.Discard(reader.Buffered())
		}
		return "", nil
	}

	if rn == '\r' || rn == '\n' {
		return "ENTER", nil
	}

	return string(unicode.ToUpper(rn)), nil
}
