package term

import (
	"encoding/base64"
)

const ResetColourSequence = "\033[0m" // anything after this will have no custom colouring applied

func Base64Encode(sequence []byte) string {
	return base64.RawURLEncoding.EncodeToString(sequence)
}
