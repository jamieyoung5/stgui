package term

import (
	"bufio"
	"encoding/base64"
)

// TODO: improve this implementation to be a bit more robust
func ReadKeySequence(reader *bufio.Reader) ([]byte, error) {
	input, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	keySequence := []byte{input}
	if input == 27 { // ESC character, indicating the start of a control sequence
		// Read the full sequence (assuming fixed length for simplicity)
		for i := 0; i < 2; i++ {
			if nextByte, err := reader.ReadByte(); err == nil {
				keySequence = append(keySequence, nextByte)
			}
		}
	}

	if input > 'A' && input <= 'Z' {
		keySequence[0] += 32
	}

	return keySequence, nil
}

func Encode(sequence []byte) string {
	return base64.RawURLEncoding.EncodeToString(sequence)
}
