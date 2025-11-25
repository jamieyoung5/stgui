package keyboard

import (
	"bufio"
	"os"
	"time"
)

const (
	EnterKey      = "ENTER"
	EscapeKey     = "ESCAPE"
	UpArrowKey    = "UP"
	DownArrowKey  = "DOWN"
	LeftArrowKey  = "LEFT"
	RightArrowKey = "RIGHT"

	keySequenceTimeout = 50 * time.Millisecond
)

func ReadInput(reader *bufio.Reader, file *os.File) (string, error) {
	rn, _, err := reader.ReadRune()
	if err != nil {
		return "", err
	}

	if rn == '\r' || rn == '\n' {
		return EnterKey, nil
	}

	seq := string(rn)

	if !sequencesTrie.StartsWith(seq) {
		return seq, nil
	}

	for {

		if cmd, ok := Sequences[seq]; ok {
			return cmd, nil
		}

		if ok := sequencesTrie.StartsWith(seq); !ok {
			return seq, nil
		}

		nextRn, err := readRuneWithDeadline(reader, file, keySequenceTimeout)
		if err != nil {
			if seq == "\x1b" {
				return EscapeKey, nil
			}
			return seq, nil
		}
		seq += string(nextRn)
	}
}

func readRuneWithDeadline(r *bufio.Reader, f *os.File, timeout time.Duration) (rune, error) {
	if err := f.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return 0, err
	}
	// ensure deadline is cleared so future reads dont fail instantly
	defer f.SetReadDeadline(time.Time{})

	rn, _, err := r.ReadRune()
	return rn, err
}
