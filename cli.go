package cli

import (
	"fmt"
	"strings"
)

// Colour codes constants
const (
	Black = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// SGR-codes list constants
const (
	Reset = iota
	Bold
	Faint
	Italic
	Underline
)

const (
	// ClearScreen is an escape code to clear the screen
	ClearScreen = "\x1b[2J"
	// Move is an escape code for move the cursor to (x,y)
	Move = "\x1b[%d;%dH"
	// LineReset is an escape code to return cursor to the line beginning and line cleanup
	LineReset = "\r\x1b[K"
)

// getColor returns an escape code for colour with given colour code
func getColor(code int) string {
	return getParam(30 + code)
}

// getBgColor returns an escape code for background colour with given colour code
func getBgColor(code int) string {
	return getParam(40 + code)
}

// getParam returns an escape code for text parameters
func getParam(code int) string {
	return fmt.Sprintf("\x1b[%dm", code)
}

// Println like Printf but adds a new line to the end of string
func Println(str string, a ...interface{}) {
	Printf(str+"\n", a...)
}

// Printf the string str with specifier substitutions
func Printf(str string, a ...interface{}) {
	fmt.Printf(Colorize(str, a...))
}

// Sprintf returns a string str after specifier substitutions
func Sprintf(str string, a ...interface{}) string {
	return fmt.Sprintf(Colorize(str, a...))
}

// Colorize returns a coloured string, converting {-tags
// Ex.: cli.Colorize("{Rred string{0 and {Bblue part{0")
// Note.
//		Old format: {G
//		New format: {G|
//		Migration:
//		replace
//		  \{((?:(?:_)?(?:[wargybmcWARGYBMC]))|(?:[ius0]))
//		with
//		  {$1|
//		in string literals only
func Colorize(str string, a ...interface{}) string {
	const (
		prefix  = "{"
		postfix = "|"
	)

	str = fmt.Sprintf(str, a...)
	changeMap := map[string]string{
		prefix + "w" + postfix:  getColor(White),
		prefix + "a" + postfix:  getColor(Black),
		prefix + "r" + postfix:  getColor(Red),
		prefix + "g" + postfix:  getColor(Green),
		prefix + "y" + postfix:  getColor(Yellow),
		prefix + "b" + postfix:  getColor(Blue),
		prefix + "m" + postfix:  getColor(Magenta),
		prefix + "c" + postfix:  getColor(Cyan),
		prefix + "W" + postfix:  getParam(Bold) + getColor(White),
		prefix + "A" + postfix:  getParam(Bold) + getColor(Black),
		prefix + "R" + postfix:  getParam(Bold) + getColor(Red),
		prefix + "G" + postfix:  getParam(Bold) + getColor(Green),
		prefix + "Y" + postfix:  getParam(Bold) + getColor(Yellow),
		prefix + "B" + postfix:  getParam(Bold) + getColor(Blue),
		prefix + "M" + postfix:  getParam(Bold) + getColor(Magenta),
		prefix + "C" + postfix:  getParam(Bold) + getColor(Cyan),
		prefix + "_w" + postfix: getBgColor(White),
		prefix + "_a" + postfix: getBgColor(Black),
		prefix + "_r" + postfix: getBgColor(Red),
		prefix + "_g" + postfix: getBgColor(Green),
		prefix + "_y" + postfix: getBgColor(Yellow),
		prefix + "_b" + postfix: getBgColor(Blue),
		prefix + "_m" + postfix: getBgColor(Magenta),
		prefix + "_c" + postfix: getBgColor(Cyan),
		prefix + "_W" + postfix: getParam(Bold) + getBgColor(White),
		prefix + "_A" + postfix: getParam(Bold) + getBgColor(Black),
		prefix + "_R" + postfix: getParam(Bold) + getBgColor(Red),
		prefix + "_G" + postfix: getParam(Bold) + getBgColor(Green),
		prefix + "_Y" + postfix: getParam(Bold) + getBgColor(Yellow),
		prefix + "_B" + postfix: getParam(Bold) + getBgColor(Blue),
		prefix + "_M" + postfix: getParam(Bold) + getBgColor(Magenta),
		prefix + "_C" + postfix: getParam(Bold) + getBgColor(Cyan),
		prefix + "i" + postfix:  getParam(Italic),
		prefix + "u" + postfix:  getParam(Underline),
		prefix + "s" + postfix:  ClearScreen,
	}
	for key, value := range changeMap {
		str = strings.Replace(str, key, getParam(Reset)+value, -1)
	}
	str = strings.Replace(str, prefix+"0"+postfix, getParam(Reset), -1)
	return str
}
