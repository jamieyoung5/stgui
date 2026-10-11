package keyboard

import (
	"bufio"
	"os"
	"time"
	"unicode"
	"unicode/utf8"
)

// Key is one press. Named keys use the constants below; anything else is the
// characters typed.
type Key string

// Keys with names. Anything else comes back as the characters typed.
const (
	KeyEnter     Key = "ENTER"
	KeyEscape    Key = "ESCAPE"
	KeyBackspace Key = "BACKSPACE"
	KeyTab       Key = "TAB"
	KeyCtrlC     Key = "CTRL+C"
	KeyUp        Key = "UP"
	KeyDown      Key = "DOWN"
	KeyLeft      Key = "LEFT"
	KeyRight     Key = "RIGHT"
)

const keySequenceTimeout = 200 * time.Millisecond

// Rune is the character typed, and whether there was one. Named keys and
// escape sequences report false.
func (k Key) Rune() (rune, bool) {
	r, size := utf8.DecodeRuneInString(string(k))
	if size != len(k) || r == utf8.RuneError {
		return 0, false
	}

	// A lone control character is a key press, not text.
	if unicode.IsControl(r) {
		return 0, false
	}

	return r, true
}

// Single runes that map straight to a named key.
var namedKeys = map[rune]Key{
	'\r':   KeyEnter,
	'\n':   KeyEnter,
	'\t':   KeyTab,
	'\x03': KeyCtrlC,
	'\x7f': KeyBackspace,
	'\b':   KeyBackspace,
}

// ReadInput blocks until a key is pressed, resolving escape sequences like the
// arrow keys to their names. file is the terminal behind reader, needed to time
// out a lone Escape press.
func ReadInput(reader *bufio.Reader, file *os.File) (Key, error) {
	rn, _, err := reader.ReadRune()
	if err != nil {
		return "", err
	}

	if key, ok := namedKeys[rn]; ok {
		return key, nil
	}

	seq := string(rn)

	if !sequencesTrie.StartsWith(seq) {
		return Key(seq), nil
	}

	for {

		if cmd, ok := Sequences[seq]; ok {
			return cmd, nil
		}

		if ok := sequencesTrie.StartsWith(seq); !ok {
			return Key(seq), nil
		}

		nextRn, err := readRuneWithDeadline(reader, file, keySequenceTimeout)
		if err != nil {
			if seq == "\x1b" {
				return KeyEscape, nil
			}
			return Key(seq), nil
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
