package stgui

import "strconv"

// Sequences for Cell.Style. Concatenate to combine: StyleBold + StyleRed. The
// reset is added for you.
const (
	StyleBold      = "\033[1m"
	StyleDim       = "\033[2m"
	StyleUnderline = "\033[4m"
	StyleReverse   = "\033[7m"

	StyleBlack   = "\033[30m"
	StyleRed     = "\033[31m"
	StyleGreen   = "\033[32m"
	StyleYellow  = "\033[33m"
	StyleBlue    = "\033[34m"
	StyleMagenta = "\033[35m"
	StyleCyan    = "\033[36m"
	StyleWhite   = "\033[37m"

	StyleBGBlack   = "\033[40m"
	StyleBGRed     = "\033[41m"
	StyleBGGreen   = "\033[42m"
	StyleBGYellow  = "\033[43m"
	StyleBGBlue    = "\033[44m"
	StyleBGMagenta = "\033[45m"
	StyleBGCyan    = "\033[46m"
	StyleBGWhite   = "\033[47m"
	StyleBGGrey    = "\033[100m"
)

// Fg256 picks a text colour from the 256 colour palette, for the shades the
// constants above don't cover. Out of range values are clamped.
func Fg256(colour int) string {
	return "\033[38;5;" + strconv.Itoa(clampColour(colour)) + "m"
}

// Bg256 is Fg256 for backgrounds.
func Bg256(colour int) string {
	return "\033[48;5;" + strconv.Itoa(clampColour(colour)) + "m"
}

func clampColour(colour int) int {
	return min(max(colour, 0), 255)
}

// Align says where a cell's contents sit in their column.
type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)
