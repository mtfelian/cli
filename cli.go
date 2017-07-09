package cli

import (
	"fmt"
	"strings"
)

// список возможных кодов цветов
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

// список возможных SGR-кодов
const (
	Reset = iota
	Bold
	Faint
	Italic
	Underline
)

const (
	// ClearScreen это escape-код для очистки экрана
	ClearScreen = "\x1b[2J"
	// Move это escape-код для перемещения курсора на x, y
	Move = "\x1b[%d;%dH"
	// LineReset это escape-код для возврата курсора в начало строки и очистки её
	LineReset = "\r\x1b[K"
)

// getColor возвращает escape-код для цвета с кодом code
func getColor(code int) string {
	return getParam(30 + code)
}

// getBgColor возвращает escape-код для фона цвета с кодом code
func getBgColor(code int) string {
	return getParam(40 + code)
}

// getParam возвращает escape-код для параметров текста
func getParam(code int) string {
	return fmt.Sprintf("\x1b[%dm", code)
}

// Println печатает строку str и добавляет в конец перенос строки
func Println(str string, a ...interface{}) {
	Printf(str+"\n", a...)
}

// Printf печатает строку str и подставляет туда параметры
func Printf(str string, a ...interface{}) {
	fmt.Printf(Colorize(str, a...))
}

// Sprintf возвращает строку str и подставляет туда параметры
func Sprintf(str string, a ...interface{}) string {
	return fmt.Sprintf(Colorize(str, a...))
}

// Colorize возвращает цветную строку преобразуя {-теги
// Пример: cli.Colorize("{Rred string{0 and {Bblue part{0")
func Colorize(str string, a ...interface{}) string {
	const prefix = "{"
	str = fmt.Sprintf(str, a...)
	changeMap := map[string]string{
		prefix + "w":  getColor(White),
		prefix + "a":  getColor(Black),
		prefix + "r":  getColor(Red),
		prefix + "g":  getColor(Green),
		prefix + "y":  getColor(Yellow),
		prefix + "b":  getColor(Blue),
		prefix + "m":  getColor(Magenta),
		prefix + "c":  getColor(Cyan),
		prefix + "W":  getParam(Bold) + getColor(White),
		prefix + "A":  getParam(Bold) + getColor(Black),
		prefix + "R":  getParam(Bold) + getColor(Red),
		prefix + "G":  getParam(Bold) + getColor(Green),
		prefix + "Y":  getParam(Bold) + getColor(Yellow),
		prefix + "B":  getParam(Bold) + getColor(Blue),
		prefix + "M":  getParam(Bold) + getColor(Magenta),
		prefix + "C":  getParam(Bold) + getColor(Cyan),
		prefix + "_w": getBgColor(White),
		prefix + "_a": getBgColor(Black),
		prefix + "_r": getBgColor(Red),
		prefix + "_g": getBgColor(Green),
		prefix + "_y": getBgColor(Yellow),
		prefix + "_b": getBgColor(Blue),
		prefix + "_m": getBgColor(Magenta),
		prefix + "_c": getBgColor(Cyan),
		prefix + "_W": getParam(Bold) + getBgColor(White),
		prefix + "_A": getParam(Bold) + getBgColor(Black),
		prefix + "_R": getParam(Bold) + getBgColor(Red),
		prefix + "_G": getParam(Bold) + getBgColor(Green),
		prefix + "_Y": getParam(Bold) + getBgColor(Yellow),
		prefix + "_B": getParam(Bold) + getBgColor(Blue),
		prefix + "_M": getParam(Bold) + getBgColor(Magenta),
		prefix + "_C": getParam(Bold) + getBgColor(Cyan),
		prefix + "i":  getParam(Italic),
		prefix + "u":  getParam(Underline),
		prefix + "s":  ClearScreen,
	}
	for key, value := range changeMap {
		str = strings.Replace(str, key, getParam(Reset)+value, -1)
	}
	str = strings.Replace(str, prefix+"0", getParam(Reset), -1)
	return str
}
