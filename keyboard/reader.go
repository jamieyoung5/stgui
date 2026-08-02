package keyboard

import (
	"bufio"
	"os"
	"time"
)

// Keys with names. Anything else comes back as the characters typed.
const (
	EnterKey      = "ENTER"
	EscapeKey     = "ESCAPE"
	BackspaceKey  = "BACKSPACE"
	TabKey        = "TAB"
	CtrlCKey      = "CTRL+C"
	UpArrowKey    = "UP"
	DownArrowKey  = "DOWN"
	LeftArrowKey  = "LEFT"
	RightArrowKey = "RIGHT"

	keySequenceTimeout = 200 * time.Millisecond
)

// Single runes that map straight to a named key.
var namedKeys = map[rune]string{
	'\r':   EnterKey,
	'\n':   EnterKey,
	'\t':   TabKey,
	'\x03': CtrlCKey,
	'\x7f': BackspaceKey,
	'\b':   BackspaceKey,
}

// ReadInput blocks until a key is pressed, resolving escape sequences like the
// arrow keys to their names. file is the terminal behind reader, needed to time
// out a lone Escape press.
func ReadInput(reader *bufio.Reader, file *os.File) (string, error) {
	rn, _, err := reader.ReadRune()
	if err != nil {
		return "", err
	}

	if key, ok := namedKeys[rn]; ok {
		return key, nil
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
	// Optimization: If data is already in the buffer, read immediately.
	// This avoids syscall overhead and potential errors with SetReadDeadline
	// when the full sequence is already waiting (which is typical for arrow keys).
	if r.Buffered() > 0 {
		rn, _, err := r.ReadRune()
		return rn, err
	}

	if err := f.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return 0, err
	}
	// ensure deadline is cleared so future reads dont fail instantly
	defer f.SetReadDeadline(time.Time{})

	rn, _, err := r.ReadRune()
	return rn, err
}
