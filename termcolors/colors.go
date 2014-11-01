package termcolors

type Color []byte

const keyEscape = 27

var (
	Black   = Color{keyEscape, '[', '3', '0', 'm'}
	Red     = Color{keyEscape, '[', '3', '1', 'm'}
	Green   = Color{keyEscape, '[', '3', '2', 'm'}
	Yellow  = Color{keyEscape, '[', '3', '3', 'm'}
	Blue    = Color{keyEscape, '[', '3', '4', 'm'}
	Magenta = Color{keyEscape, '[', '3', '5', 'm'}
	Cyan    = Color{keyEscape, '[', '3', '6', 'm'}
	White   = Color{keyEscape, '[', '3', '7', 'm'}
	Reset   = Color{keyEscape, '[', '0', 'm'}
)

// not efficient.
func WrapString(c Color, s string) string {
	return string(c) + s + string(Reset)
}
